package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"purchase/handler"
	"purchase/model"
	"purchase/pb"
	"purchase/repository"
	"purchase/service"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// gRPC server struct
type grpcCartServer struct {
	pb.UnimplementedCartServiceServer
	service *service.PurchaseService
}

func (s *grpcCartServer) GetCart(ctx context.Context, req *pb.CartRequest) (*pb.CartResponse, error) {
	log.Printf("[gRPC Server] Getting cart for tourist: %d", req.TouristId)
	cart, err := s.service.GetOrCreateCart(ctx, req.TouristId)
	if err != nil {
		return nil, err
	}

	var items []*pb.OrderItemProto
	for _, item := range cart.Items {
		items = append(items, &pb.OrderItemProto{
			Id:             uint32(item.ID),
			ShoppingCartId: uint32(item.ShoppingCartID),
			TourId:         item.TourID,
			TourName:       item.TourName,
			Price:          item.Price,
		})
	}

	return &pb.CartResponse{
		Id:         uint32(cart.ID),
		TouristId:  cart.TouristID,
		TotalPrice: cart.TotalPrice,
		Items:      items,
		CreatedAt:  cart.CreatedAt.String(),
		UpdatedAt:  cart.UpdatedAt.String(),
	}, nil
}

func main() {

	dsn := os.Getenv("DB_URL")

	log.Println("Connecting with DSN:", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "z2S4p9X8v6w3y1z5A7b9C0d2E4f6G8h0" // default
	}

	err = db.AutoMigrate(&model.ShoppingCart{}, &model.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// connect to Tours gRPC server
	toursGrpcAddr := os.Getenv("TOURS_GRPC_URL")
	if toursGrpcAddr == "" {
		toursGrpcAddr = "localhost:9083" // default
	}

	conn, err := grpc.Dial(toursGrpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Tours gRPC server: %v", err)
	}
	defer conn.Close()

	toursClient := pb.NewTourCheckServiceClient(conn)

	repo := repository.NewPurchaseRepository(db)
	serv := service.NewPurchaseService(repo, toursClient)
	hand := handler.NewPurchaseHandler(serv, jwtSecret)

	// gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":9084")
		if err != nil {
			log.Fatalf("gRPC listen failed: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterCartServiceServer(grpcServer, &grpcCartServer{service: serv})
		log.Println("gRPC server listening on :9084")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC serve failed: %v", err)
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/api/purchase/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "Purchase service healthy"}`))
	}).Methods("GET")

	r.HandleFunc("/api/purchase/cart", hand.GetCart).Methods("GET")
	r.HandleFunc("/api/purchase/cart/items", hand.AddItem).Methods("POST")
	r.HandleFunc("/api/purchase/cart/items/{id}", hand.RemoveItem).Methods("DELETE")

	port := ":8084"
	fmt.Printf("Purchase service listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, setupCORS(r)))
}

func setupCORS(router *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		router.ServeHTTP(w, r)
	})
}
