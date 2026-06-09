package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"purchase/service"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type PurchaseHandler struct {
	service   *service.PurchaseService
	jwtSecret string
}

func NewPurchaseHandler(service *service.PurchaseService, jwtSecret string) *PurchaseHandler {
	return &PurchaseHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

type AddItemPayload struct {
	TourID   string  `json:"tour_id"`
	TourName string  `json:"tour_name"`
	Price    float64 `json:"price"`
}

func (h *PurchaseHandler) GetUserIdFromToken(w http.ResponseWriter, r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return "", "", errors.New("Missing authorization header")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("Invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("Invalid claims")
	}
	userId := claims["sub"].(string)
	role := claims["role"].(string)
	return userId, role, nil
}

func (h *PurchaseHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userIdStr, role, err := h.GetUserIdFromToken(w, r)
	if err != nil {
		return
	}
	if role != "TOURIST" {
		http.Error(w, "Only tourists can view their cart", http.StatusForbidden)
		return
	}
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	cart, err := h.service.GetOrCreateCart(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cart)
}

func (h *PurchaseHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userIdStr, role, err := h.GetUserIdFromToken(w, r)
	if err != nil {
		return
	}
	if role != "TOURIST" {
		http.Error(w, "Only tourists can add items to the cart", http.StatusForbidden)
		return
	}
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	var payload AddItemPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.service.AddItemToCart(r.Context(), userId, payload.TourID, payload.TourName, payload.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *PurchaseHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userIdStr, role, err := h.GetUserIdFromToken(w, r)
	if err != nil {
		return
	}
	if role != "TOURIST" {
		http.Error(w, "Only tourists can remove items from the cart", http.StatusForbidden)
		return
	}
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	vars := mux.Vars(r)
	itemIDStr := vars["id"]
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	err = h.service.RemoveItemFromCart(r.Context(), userId, uint(itemID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Item removed"}`))
}

func (h *PurchaseHandler) CheckPurchase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	tourID := vars["tourId"]

	touristIdStr := r.URL.Query().Get("touristId")
	touristId, err := strconv.ParseInt(touristIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid touristId", http.StatusBadRequest)
		return
	}

	has, err := h.service.HasPurchasedTour(touristId, tourID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"purchased": has})
}

func (h *PurchaseHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIdStr, role, err := h.GetUserIdFromToken(w, r)
	if err != nil {
		return
	}
	if role != "TOURIST" {
		http.Error(w, "Only tourists can checkout", http.StatusForbidden)
		return
	}

	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	err = h.service.CheckoutCartAsync(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Checkout initiated. Processing in background.",
	})
}
