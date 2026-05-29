package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

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
		if strings.HasPrefix(path, "/api/auth") || strings.HasPrefix(path, "/api/users") {
			serveReverseProxy(stakeholdersURL, w, r)
		} else if strings.HasPrefix(path, "/blogs") {
			serveReverseProxy(blogURL, w, r)
		} else if strings.HasPrefix(path, "/api/followers") {
			serveReverseProxy(followersURL, w, r)
		} else if strings.HasPrefix(path, "/api/tours") || strings.HasPrefix(path, "/api/position") {
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
