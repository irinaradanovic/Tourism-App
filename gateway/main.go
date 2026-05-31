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

func initGrpcClient() {
	conn, err := grpc.NewClient("purchase:9084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to purchase gRPC: %v", err)
	}
	cartClient = pb.NewCartServiceClient(conn)
	log.Println("gRPC client connected to purchase:9084")
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

// Help function to handle reverse proxying to target services
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, err := url.Parse(target)
	if err != nil {
		log.Printf("Error while parsing URL: %v", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update headers so services know where the request came from
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
	initGrpcClient()
	// Define routes to our services within the Docker network
	stakeholdersURL := "http://stakeholders:8082"
	blogURL := "http://blog:8081"
	followersURL := "http://followers:8000"
	toursURL := "http://tours:8083"
	purchaseURL := "http://purchase:8084"

	// Main handler that acts as an Nginx reverse proxy
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		log.Printf("[Gateway] Received request: %s %s", r.Method, path)

		// Route requests based on path prefixes
		if path == "/api/purchase/cart" && (r.Method == "GET" || r.Method == "OPTIONS") {
			handleGetCartGrpc(w, r, jwtSecret)
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
