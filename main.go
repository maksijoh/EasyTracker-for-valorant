package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "\ufeff"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), "\"'")
		if key != "" && os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func getData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accessToken := os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		http.Error(w, "ACCESS_TOKEN is not configured", http.StatusInternalServerError)
		return
	}

	name := strings.TrimSpace(r.URL.Query().Get("name"))
	tag := strings.TrimSpace(r.URL.Query().Get("tag"))
	region := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("region")))
	if region == "" {
		region = "eu"
	}
	if name == "" || tag == "" {
		http.Error(w, "name and tag query parameters are required", http.StatusBadRequest)
		return
	}
	validRegions := map[string]bool{"eu": true, "na": true, "latam": true, "br": true, "ap": true, "kr": true}
	if !validRegions[region] {
		http.Error(w, "invalid region", http.StatusBadRequest)
		return
	}

	apiURL := "https://api.henrikdev.xyz/valorant/v4/matches/" +
		url.PathEscape(region) + "/pc/" + url.PathEscape(name) + "/" + url.PathEscape(tag) + "?size=10"

	req, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodGet,
		apiURL,
		nil,
	)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", accessToken)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "External API is unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("response write error: %v", err)
	}
}

func main() {
	if err := loadEnv(".env"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to load .env: %v", err)
	}

	http.HandleFunc("/api/data", getData)

	log.Println("http://localhost:8080/api/data")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
