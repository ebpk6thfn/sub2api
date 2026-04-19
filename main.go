package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/your-org/sub2api/handler"
)

const defaultPort = 8080

func main() {
	var (
		port    int
		verbose bool
	)

	flag.IntVar(&port, "port", getEnvInt("PORT", defaultPort), "HTTP server port")
	flag.BoolVar(&verbose, "verbose", true, "Enable verbose logging") // personal: default verbose to true
	flag.Parse()

	if verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	// Subscription conversion endpoints
	mux.HandleFunc("/sub", handler.SubHandler)
	mux.HandleFunc("/clash", handler.ClashHandler)
	mux.HandleFunc("/singbox", handler.SingBoxHandler)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("sub2api listening on %s", addr)
	log.Printf("verbose logging: %v", verbose)

	// personal: use a custom server with timeouts to avoid hanging connections
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// getEnvInt reads an integer from an environment variable, returning fallback on error.
func getEnvInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}
