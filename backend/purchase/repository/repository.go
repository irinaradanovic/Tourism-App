package repository

import (
	"context"
	"purchase/model"

	"gorm.io/gorm"
)

type IPurchaseRepository interface {
	SaveCart(ctx context.Context, cart *model.ShoppingCart) error
	GetCartByTouristId(touristID int64) (*model.ShoppingCart, error)
	CreateCart(ctx context.Context, cart *model.ShoppingCart) error
	CreateItem(ctx context.Context, item *model.OrderItem) error
	GetItemById(ctx context.Context, itemID uint) (*model.OrderItem, error)
	DeleteItem(ctx context.Context, item *model.OrderItem) error
	HasToken(touristID int64, tourID string) (bool, error)
	CreateToken(ctx context.Context, token *model.TourPurchaseToken) error
	ClearCartItems(ctx context.Context, cartID uint) error
	DeleteToken(ctx context.Context, tokenID uint) error
}

type PurchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) IPurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) GetCartByTouristId(touristID int64) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart
	err := r.db.Preload("Items").Where("tourist_id = ?", touristID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *PurchaseRepository) CreateCart(ctx context.Context, cart *model.ShoppingCart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *PurchaseRepository) SaveCart(ctx context.Context, cart *model.ShoppingCart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

func (r *PurchaseRepository) CreateItem(ctx context.Context, item *model.OrderItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PurchaseRepository) GetItemById(ctx context.Context, itemID uint) (*model.OrderItem, error) {
	var item model.OrderItem
	err := r.db.WithContext(ctx).First(&item, itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PurchaseRepository) DeleteItem(ctx context.Context, item *model.OrderItem) error {
	return r.db.WithContext(ctx).Delete(item).Error
}

func (r *PurchaseRepository) HasToken(touristID int64, tourID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.TourPurchaseToken{}).
		Where("tourist_id = ? AND tour_id = ?", touristID, tourID).
		Count(&count).Error
	return count > 0, err
}

func (r *PurchaseRepository) CreateToken(ctx context.Context, token *model.TourPurchaseToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *PurchaseRepository) ClearCartItems(ctx context.Context, cartID uint) error {
	return r.db.WithContext(ctx).Where("shopping_cart_id= ?", cartID).Delete(&model.OrderItem{}).Error
}

func (r *PurchaseRepository) DeleteToken(ctx context.Context, tokenID uint) error {
	return r.db.WithContext(ctx).Delete(&model.TourPurchaseToken{}, tokenID).Error
}
