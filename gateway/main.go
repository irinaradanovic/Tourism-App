package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/pb"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cartClient pb.CartServiceClient
var proximityClient pb.ProximityServiceClient
var checkoutClient pb.CheckoutServiceClient

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func initGrpcClients() {
	connPurchase, err := grpc.NewClient("purchase:9084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to purchase gRPC: %v", err)
	}
	cartClient = pb.NewCartServiceClient(connPurchase)
	log.Println("gRPC client connected to purchase:9084")

	connTours, err := grpc.NewClient("tours:9083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to tours gRPC: %v", err)
	}
	proximityClient = pb.NewProximityServiceClient(connTours)
	log.Println("gRPC client connected to tours:9083")
}

func initCheckoutClient() {
	conn, err := grpc.Dial("purchase:9084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to purchase gRPC: %v", err)
	}
	checkoutClient = pb.NewCheckoutServiceClient(conn)
	log.Println("gRPC client connected to purchase:9084 (checkout)")
}

func getTouristIdFromToken(r *http.Request, secret string) (int64, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("missing auth header")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	var id int64
	fmt.Sscanf(sub, "%d", &id)
	return id, nil
}

func handleGetCartGrpc(w http.ResponseWriter, r *http.Request, jwtSecret string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	touristId, err := getTouristIdFromToken(r, jwtSecret)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cart, err := cartClient.GetCart(ctx, &pb.CartRequest{TouristId: touristId})
	if err != nil {
		http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}
func handleCheckoutGrpc(w http.ResponseWriter, r *http.Request, jwtSecret string) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	touristId, err := getTouristIdFromToken(r, jwtSecret)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := checkoutClient.Checkout(ctx, &pb.CheckoutRequest{TouristId: touristId})
	if err != nil {
		http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleCheckProximityGrpc(w http.ResponseWriter, r *http.Request, jwtSecret string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	touristId, err := getTouristIdFromToken(r, jwtSecret)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	executionId := parts[2]

	var body struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := proximityClient.CheckProximity(ctx, &pb.ProximityRequest{
		ExecutionId: executionId,
		TouristId:   touristId,
		Lat:         body.Lat,
		Lon:         body.Lon,
	})
	if err != nil {
		http.Error(w, "gRPC proximity error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[Gateway] Proximity response: executionId=%s status=%s", resp.ExecutionId, resp.Status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func isProximityGrpcRoute(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return len(parts) == 5 &&
		parts[0] == "api" &&
		parts[1] == "executions" &&
		parts[3] == "proximity" &&
		parts[4] == "grpc"
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, err := url.Parse(target)
	if err != nil {
		log.Printf("Error while parsing URL: %v", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		if auth := req.Header.Get("Authorization"); auth != "" {
			req.Header.Set("Authorization", auth)
		}
	}

	proxy.ServeHTTP(res, req)
}

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	initGrpcClients()
	initCheckoutClient()
	stakeholdersURL := "http://stakeholders:8082"
	blogURL := "http://blog:8081"
	followersURL := "http://followers:8000"
	toursURL := "http://tours:8083"
	purchaseURL := "http://purchase:8084"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("[Gateway] Received request: %s %s", r.Method, path)

		if path == "/api/purchase/cart" && (r.Method == "GET" || r.Method == "OPTIONS") {
			handleGetCartGrpc(w, r, jwtSecret)
		} else if isProximityGrpcRoute(path) && (r.Method == "POST" || r.Method == "OPTIONS") {
			handleCheckProximityGrpc(w, r, jwtSecret)
		} else if path == "/api/purchase/checkout" && (r.Method == "POST" || r.Method == "OPTIONS") {
			handleCheckoutGrpc(w, r, jwtSecret)
		} else if strings.HasPrefix(path, "/api/auth") || strings.HasPrefix(path, "/api/users") {
			serveReverseProxy(stakeholdersURL, w, r)
		} else if strings.HasPrefix(path, "/blogs") {
			serveReverseProxy(blogURL, w, r)
		} else if strings.HasPrefix(path, "/api/followers") {
			serveReverseProxy(followersURL, w, r)
		} else if strings.HasPrefix(path, "/api/tours") ||
			strings.HasPrefix(path, "/api/position") ||
			strings.HasPrefix(path, "/api/executions") {
			serveReverseProxy(toursURL, w, r)
		} else if strings.HasPrefix(path, "/api/purchase") {
			serveReverseProxy(purchaseURL, w, r)
		} else {
			http.Error(w, "Route not found on Gateway", http.StatusNotFound)
		}
	})

	log.Println("Custom API Gateway started on port :80...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("Gateway failed: %v", err)
	}
}
