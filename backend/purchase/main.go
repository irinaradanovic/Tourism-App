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
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var tracer trace.Tracer

type grpcMetadataCarrier struct {
	md metadata.MD
}

func (c grpcMetadataCarrier) Get(key string) string {
	values := c.md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c grpcMetadataCarrier) Set(key, value string) {
	c.md.Set(key, value)
}

func (c grpcMetadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c.md))
	for key := range c.md {
		keys = append(keys, key)
	}
	return keys
}

func initTracer(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	if serviceName == "" {
		serviceName = "purchase-service"
	}
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://tempo:4318/v1/traces"
	}

	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx, resource.WithAttributes(attribute.String("service.name", serviceName)))
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer = otel.Tracer(serviceName)

	return provider.Shutdown, nil
}

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

type grpcCheckoutServer struct {
	pb.UnimplementedCheckoutServiceServer
	service *service.PurchaseService
}

func (s *grpcCheckoutServer) Checkout(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ctx = otel.GetTextMapPropagator().Extract(ctx, grpcMetadataCarrier{md: md})
	}
	ctx, span := tracer.Start(
		ctx,
		"PurchaseService.Checkout",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.service", "CheckoutService"),
			attribute.String("rpc.method", "Checkout"),
			attribute.Int64("tourist.id", req.TouristId),
		),
	)
	defer span.End()

	log.Printf("[gRPC Server] Checkout request for tourist: %d", req.TouristId)

	checkoutCtx, checkoutSpan := tracer.Start(
		ctx,
		"PurchaseService.CheckoutCartAsync",
		trace.WithAttributes(attribute.Int64("tourist.id", req.TouristId)),
	)
	defer checkoutSpan.End()
	err := s.service.CheckoutCartAsync(checkoutCtx, req.TouristId)
	if err != nil {
		checkoutSpan.RecordError(err)
		checkoutSpan.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Printf("[gRPC Server] Error in CheckoutCartAsync for tourist %d: %v", req.TouristId, err)
		return nil, err
	}
	checkoutSpan.SetStatus(codes.Ok, "")
	span.SetStatus(codes.Ok, "")
	return &pb.CheckoutResponse{Message: "Checkout initiated in background", Tokens: nil}, nil
}

func main() {
	tracerShutdown, err := initTracer(context.Background(), os.Getenv("OTEL_SERVICE_NAME"))
	if err != nil {
		log.Printf("OpenTelemetry tracing disabled: %v", err)
		tracer = otel.Tracer("purchase-service")
	} else {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := tracerShutdown(ctx); err != nil {
				log.Printf("Failed to shutdown tracer provider: %v", err)
			}
		}()
	}

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

	err = db.AutoMigrate(&model.ShoppingCart{}, &model.OrderItem{}, &model.TourPurchaseToken{})
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
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	rabbitConn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	ch, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel on RabbitMQ: %v", err)
	}
	defer ch.Close()

	repo := repository.NewPurchaseRepository(db)
	serv := service.NewPurchaseService(repo, toursClient, ch)
	serv.StartSagaConsumers()
	hand := handler.NewPurchaseHandler(serv, jwtSecret)

	// gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":9084")
		if err != nil {
			log.Fatalf("gRPC listen failed: %v", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterCartServiceServer(grpcServer, &grpcCartServer{service: serv})
		pb.RegisterCheckoutServiceServer(grpcServer, &grpcCheckoutServer{service: serv})
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
	r.HandleFunc("/api/purchase/check/{tourId}", hand.CheckPurchase).Methods("GET")
	r.HandleFunc("/api/purchase/checkout", hand.Checkout).Methods("POST")

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
