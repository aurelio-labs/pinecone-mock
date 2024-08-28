package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bruvduroiu/pinecone-mock/index"
	"github.com/bruvduroiu/pinecone-mock/vector"
)

func main() {
	router := http.NewServeMux()

	indexHandler := index.Handler{}
	vectorHandler := vector.Handler{
		Vectors: make(map[string]*vector.Vector),
	}

	router.HandleFunc("POST /indexes", indexHandler.Create)
	router.HandleFunc("GET /indexes", indexHandler.List)
	router.HandleFunc("GET /indexes/{index_name}", indexHandler.GetByName)
	router.HandleFunc("POST /describe_index_stats", indexHandler.DescribeIndexStats)
	router.HandleFunc("POST /vectors/upsert", vectorHandler.Upsert)
	router.HandleFunc("POST /query", vectorHandler.Query)
	router.HandleFunc("GET /vectors/fetch", vectorHandler.Query)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			// Read the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			// Verify that the body is valid JSON
			var jsonBody interface{}
			if err := json.Unmarshal(body, &jsonBody); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			log.Println(fmt.Sprintf("%+v", jsonBody))

			// Set the Content-Type header to application/json
			w.Header().Set("Content-Type", "application/json")

			// Write the JSON back to the response
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		}
		if r.Method == http.MethodGet {
			response := make(map[string]interface{})

			// Add query parameters to the response
			params := r.URL.Query()
			if len(params) > 0 {
				response["query_parameters"] = params
			}

			// Add headers to the response
			headers := make(map[string]string)
			for name, values := range r.Header {
				headers[name] = strings.Join(values, ", ")
			}
			response["headers"] = headers

			// Convert the response to JSON
			jsonResponse, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
				return
			}
			log.Println(fmt.Sprintf("%+v", jsonResponse))

			// Set the Content-Type header to application/json
			w.Header().Set("Content-Type", "application/json")

			// Write the JSON response
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
		}
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server listening on port :8080")
	server.ListenAndServe()
}
