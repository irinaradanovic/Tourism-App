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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var cartClient pb.CartServiceClient
var proximityClient pb.ProximityServiceClient
var checkoutClient pb.CheckoutServiceClient
var tourCommandClient pb.TourCommandServiceClient
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
		serviceName = "gateway"
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
	tourCommandClient = pb.NewTourCommandServiceClient(connTours)
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
	userID, _, err := getUserClaimsFromToken(r, secret)
	return userID, err
}

func getUserClaimsFromToken(r *http.Request, secret string) (int64, string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, "", fmt.Errorf("missing auth header")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)
	var id int64
	fmt.Sscanf(sub, "%d", &id)
	role, _ := claims["role"].(string)
	return id, role, nil
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

	ctx, span := tracer.Start(
		r.Context(),
		"POST /api/purchase/checkout",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.route", "/api/purchase/checkout"),
		),
	)
	defer span.End()

	touristId, err := getTouristIdFromToken(r, jwtSecret)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "unauthorized")
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	span.SetAttributes(attribute.Int64("tourist.id", touristId))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	grpcCtx, grpcSpan := tracer.Start(
		ctx,
		"grpc.Checkout purchase-service",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.service", "CheckoutService"),
			attribute.String("rpc.method", "Checkout"),
		),
	)
	defer grpcSpan.End()

	md, ok := metadata.FromOutgoingContext(grpcCtx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.New(nil)
	}
	otel.GetTextMapPropagator().Inject(grpcCtx, grpcMetadataCarrier{md: md})
	grpcCtx = metadata.NewOutgoingContext(grpcCtx, md)

	resp, err := checkoutClient.Checkout(grpcCtx, &pb.CheckoutRequest{TouristId: touristId})
	if err != nil {
		grpcSpan.RecordError(err)
		grpcSpan.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "gRPC error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	grpcSpan.SetStatus(codes.Ok, "")
	span.SetStatus(codes.Ok, "")

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

func handleCreateTourGrpc(w http.ResponseWriter, r *http.Request, jwtSecret string) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	userID, role, err := getUserClaimsFromToken(r, jwtSecret)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	if role != "GUIDE" {
		http.Error(w, "Only guides can create tours", http.StatusForbidden)
		return
	}

	var body struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Difficulty  string   `json:"difficulty"`
		Tags        []string `json:"tags"`
		Durations   []struct {
			TransportType string `json:"transportType"`
			Minutes       int32  `json:"minutes"`
		} `json:"durations"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var durations []*pb.TourDurationProto
	for _, d := range body.Durations {
		durations = append(durations, &pb.TourDurationProto{
			TransportType: toGrpcTransport(d.TransportType),
			Minutes:       d.Minutes,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := tourCommandClient.CreateTour(ctx, &pb.CreateTourRequest{
		Title:       body.Title,
		Description: body.Description,
		Difficulty:  toGrpcDifficulty(body.Difficulty),
		Tags:        body.Tags,
		AuthorId:    userID,
		Durations:   durations,
	})
	if err != nil {
		http.Error(w, "gRPC create tour error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Tour)
}

func handlePublishTourGrpc(w http.ResponseWriter, r *http.Request, jwtSecret string) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	userID, role, err := getUserClaimsFromToken(r, jwtSecret)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	tourID := parts[2]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := tourCommandClient.PublishTour(ctx, &pb.PublishTourRequest{
		TourId:   tourID,
		AuthorId: userID,
		Role:     role,
	})
	if err != nil {
		http.Error(w, "gRPC publish tour error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func toGrpcDifficulty(value string) pb.Difficulty {
	switch strings.ToUpper(value) {
	case "EASY":
		return pb.Difficulty_EASY
	case "MEDIUM":
		return pb.Difficulty_MEDIUM
	case "HARD":
		return pb.Difficulty_HARD
	default:
		return pb.Difficulty_DIFFICULTY_UNSPECIFIED
	}
}

func toGrpcTransport(value string) pb.TransportType {
	switch strings.ToUpper(value) {
	case "WALKING":
		return pb.TransportType_WALKING
	case "BICYCLE":
		return pb.TransportType_BICYCLE
	case "CAR":
		return pb.TransportType_CAR
	default:
		return pb.TransportType_TRANSPORT_UNSPECIFIED
	}
}

func isProximityGrpcRoute(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return len(parts) == 5 &&
		parts[0] == "api" &&
		parts[1] == "executions" &&
		parts[3] == "proximity" &&
		parts[4] == "grpc"
}

func isPublishTourGrpcRoute(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return len(parts) == 4 &&
		parts[0] == "api" &&
		parts[1] == "tours" &&
		parts[3] == "publish"
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
	tracerShutdown, err := initTracer(context.Background(), os.Getenv("OTEL_SERVICE_NAME"))
	if err != nil {
		log.Printf("OpenTelemetry tracing disabled: %v", err)
		tracer = otel.Tracer("gateway")
	} else {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := tracerShutdown(ctx); err != nil {
				log.Printf("Failed to shutdown tracer provider: %v", err)
			}
		}()
	}

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
		} else if path == "/api/tours" && (r.Method == "POST" || r.Method == "OPTIONS") {
			handleCreateTourGrpc(w, r, jwtSecret)
		} else if isPublishTourGrpcRoute(path) && (r.Method == "POST" || r.Method == "OPTIONS") {
			handlePublishTourGrpc(w, r, jwtSecret)
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
