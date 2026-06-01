package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"purchase/model"
	"purchase/pb"
	"purchase/repository"
	"time"

	"gorm.io/gorm"
)

type PurchaseService struct {
	repo        repository.IPurchaseRepository
	toursClient pb.TourCheckServiceClient
}

func NewPurchaseService(repo repository.IPurchaseRepository, toursClient pb.TourCheckServiceClient) *PurchaseService {
	return &PurchaseService{
		repo:        repo,
		toursClient: toursClient,
	}
}

func (s *PurchaseService) ValidateTourViaGrpc(ctx context.Context, tourID string) (string, string, float64, error) {
	resp, err := s.toursClient.CheckTour(ctx, &pb.TourCheckRequest{TourId: tourID})
	if err != nil {
		return "", "", 0, fmt.Errorf("error communicating with Tours gRPC service: %v", err)
	}
	return resp.Status, resp.TourName, resp.Price, nil
}

type TourServiceClientResponse struct {
	Status string `json:"status"`
}

func (s *PurchaseService) CheckTourStatusFromToursService(tourID string) (string, error) {
	toursURL := os.Getenv("TOURS_SERVICE_URL") + "/" + tourID
	resp, err := http.Get(toursURL)
	if err != nil {
		return "", fmt.Errorf("Tour service is unavailable: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("tour does not exist")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error on tours service, status code: %d", resp.StatusCode)
	}
	var tourData TourServiceClientResponse
	if err := json.NewDecoder(resp.Body).Decode(&tourData); err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	return tourData.Status, nil
}

func (s *PurchaseService) GetOrCreateCart(ctx context.Context, touristID int64) (*model.ShoppingCart, error) {
	cart, err := s.repo.GetCartByTouristId(touristID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cart = &model.ShoppingCart{
			TouristID:  touristID,
			TotalPrice: 0.0,
			Items:      []model.OrderItem{},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if createErr := s.repo.CreateCart(ctx, cart); createErr != nil {
			return nil, createErr
		}
		return cart, nil
	}
	return cart, err
}

func (s *PurchaseService) AddItemToCart(ctx context.Context, touristID int64, tourID string, tourName string, price float64) (*model.OrderItem, error) {
	status, realName, realPrice, err := s.ValidateTourViaGrpc(ctx, tourID)
	if err != nil {
		return nil, err
	}
	if status == "DRAFT" || status == "ARCHIVED" {
		return nil, fmt.Errorf("cannot add tour to cart because its status is: %s", status)
	}
	cart, err := s.GetOrCreateCart(ctx, touristID)
	if err != nil {
		return nil, err
	}
	item := &model.OrderItem{
		ShoppingCartID: cart.ID,
		TourID:         tourID,
		TourName:       realName,
		Price:          realPrice,
	}
	if err := s.repo.CreateItem(ctx, item); err != nil {
		return nil, err
	}
	cart.TotalPrice += realPrice
	cart.UpdatedAt = time.Now()
	if err := s.repo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *PurchaseService) RemoveItemFromCart(ctx context.Context, touristID int64, itemID uint) error {
	item, err := s.repo.GetItemById(ctx, itemID)
	if err != nil {
		return errors.New("Item does not exist in the cart")
	}
	if err := s.repo.DeleteItem(ctx, item); err != nil {
		return err
	}
	cart, err := s.repo.GetCartByTouristId(touristID)
	if err != nil {
		return err
	}
	cart.TotalPrice -= item.Price
	cart.UpdatedAt = time.Now()
	if cart.TotalPrice < 0 {
		cart.TotalPrice = 0
	}
	if err := s.repo.SaveCart(ctx, cart); err != nil {
		return err
	}
	return nil
}

func (s *PurchaseService) HasPurchasedTour(touristID int64, tourID string) (bool, error) {
	return s.repo.HasToken(touristID, tourID)
}

func (s *PurchaseService) CheckoutCart(ctx context.Context, touristID int64) ([]model.TourPurchaseToken, error) {
	cart, err := s.repo.GetCartByTouristId(touristID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return []model.TourPurchaseToken{}, nil
	}

	var createdTokens []model.TourPurchaseToken

	// Compensating action in case of failure during the checkout process
	rollback := func() {
		for _, token := range createdTokens {
			log.Printf("[SAGA] Rolling back token for tour %s", token.TourID)
			s.repo.DeleteToken(ctx, token.ID)
		}
	}

	// SAGA
	for _, item := range cart.Items {
		alreadyPurchased, err := s.repo.HasToken(touristID, item.TourID)
		if err != nil {
			rollback()
			return nil, fmt.Errorf("[SAGA] Error checking token: %v", err)
		}
		if alreadyPurchased {
			continue
		}

		log.Printf("[SAGA] Step 1: Validating tour %s via gRPC", item.TourID)
		status, _, _, err := s.ValidateTourViaGrpc(ctx, item.TourID)
		if err != nil {
			log.Printf("[SAGA] Step 1 FAILED: %v", err)
			rollback()
			return nil, fmt.Errorf("tour validation failed: %v", err)
		}
		if status == "ARCHIVED" || status == "DRAFT" {
			log.Printf("[SAGA] Step 1 FAILED: tour status is %s", status)
			rollback()
			return nil, fmt.Errorf("tour %s cannot be purchased, status: %s", item.TourID, status)
		}
		log.Printf("[SAGA] Step 1 OK: tour %s is %s", item.TourID, status)

		log.Printf("[SAGA] Step 2: Creating token for tour %s", item.TourID)
		token := model.TourPurchaseToken{
			TouristID: touristID,
			TourID:    item.TourID,
			TourName:  item.TourName,
			CreatedAt: time.Now(),
		}

		if err := s.repo.CreateToken(ctx, &token); err != nil {
			log.Printf("[SAGA] Step 2 FAILED: %v", err)
			rollback()
			return nil, fmt.Errorf("failed to create token: %v", err)
		}
		log.Printf("[SAGA] Step 2 OK: token created for tour %s", item.TourID)
		createdTokens = append(createdTokens, token)
	}

	log.Printf("[SAGA] Step 3: Clearing cart for tourist %d", touristID)
	if err := s.repo.ClearCartItems(ctx, cart.ID); err != nil {
		log.Printf("[SAGA] Step 3 FAILED, rolling back tokens")
		rollback()
		return nil, fmt.Errorf("failed to clear cart: %v", err)
	}
	cart.Items = []model.OrderItem{}
	cart.TotalPrice = 0
	cart.UpdatedAt = time.Now()

	if err := s.repo.SaveCart(ctx, cart); err != nil {
		return nil, err
	}
	log.Printf("[SAGA] Checkout SAGA completed successfully for tourist %d", touristID)
	return createdTokens, nil
}
