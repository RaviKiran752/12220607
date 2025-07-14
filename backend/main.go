package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"12220607/logging"
)

// ShortURL holds the data for a shortened URL
type ClickDetail struct {
	Timestamp time.Time `json:"timestamp"`
	Referrer  string    `json:"referrer"`
	Location  string    `json:"location"`
}

type ShortURL struct {
	URL       string
	CreatedAt time.Time
	Expiry    time.Time
	Hits      int
	Clicks    []ClickDetail
}

var (
	urlStore = make(map[string]*ShortURL)
	mu       sync.RWMutex
	letters  = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func generateShortcode(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		URL       string `json:"url"`
		Validity  int    `json:"validity"`
		Shortcode string `json:"shortcode"`
	}
	type Resp struct {
		Shortcode string `json:"shortcode"`
		Expiry    string `json:"expiry"`
	}
	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	if !isValidURL(req.URL) {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	if req.Validity <= 0 {
		req.Validity = 30
	}
	var code string
	mu.Lock()
	defer mu.Unlock()
	if req.Shortcode != "" {
		if _, exists := urlStore[req.Shortcode]; exists {
			http.Error(w, "Shortcode already exists", http.StatusConflict)
			return
		}
		code = req.Shortcode
	} else {
		for {
			code = generateShortcode(6)
			if _, exists := urlStore[code]; !exists {
				break
			}
		}
	}
	expiry := time.Now().Add(time.Duration(req.Validity) * time.Second)
	urlStore[code] = &ShortURL{
		URL:       req.URL,
		CreatedAt: time.Now(),
		Expiry:    expiry,
		Hits:      0,
		Clicks:    []ClickDetail{},
	}
	resp := Resp{Shortcode: code, Expiry: expiry.Format(time.RFC3339)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getShortURLStats(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/shorturls/")
	mu.RLock()
	defer mu.RUnlock()
	s, ok := urlStore[code]
	if !ok {
		http.Error(w, "Shortcode not found", http.StatusNotFound)
		return
	}
	resp := map[string]interface{}{
		"url":        s.URL,
		"created_at": s.CreatedAt.Format(time.RFC3339),
		"expiry":     s.Expiry.Format(time.RFC3339),
		"hits":       s.Hits,
		"clicks":     s.Clicks,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func redirectShortURL(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	mu.Lock()
	defer mu.Unlock()
	s, ok := urlStore[code]
	if !ok {
		http.Error(w, "Shortcode not found", http.StatusNotFound)
		return
	}
	if time.Now().After(s.Expiry) {
		http.Error(w, "Shortcode expired", http.StatusGone)
		return
	}
	s.Hits++
	ref := r.Referer()
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	loc := getLocationFromIP(ip)
	s.Clicks = append(s.Clicks, ClickDetail{
		Timestamp: time.Now(),
		Referrer:  ref,
		Location:  loc,
	})
	http.Redirect(w, r, s.URL, http.StatusFound)
}

func getLocationFromIP(ip string) string {
	// Dummy implementation: just return "Unknown" for now
	// In real-world, use a geo-IP service
	if ip == "" {
		return "Unknown"
	}
	// For demonstration, return first two octets for IPv4
	parts := strings.Split(ip, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1] + ".x.x"
	}
	return "Unknown"
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())
	mux := http.NewServeMux()
	mux.HandleFunc("/shorturls", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createShortURL(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/shorturls/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getShortURLStats(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && len(r.URL.Path) > 1 && !strings.HasPrefix(r.URL.Path, "/shorturls") {
			redirectShortURL(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	loggedMux := logging.LoggingMiddleware(mux)
	corsLoggedMux := CORSMiddleware(loggedMux)
	fmt.Println("Backend running on :3001")
	log.Fatal(http.ListenAndServe(":3001", corsLoggedMux))
}
