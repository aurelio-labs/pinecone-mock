package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bruvduroiu/pinecone-mock/pinecone"
)

func main() {
	router := http.NewServeMux()

	handler := pinecone.Handler{}

	router.HandleFunc("POST /indexes", handler.CreateIndex)
	router.HandleFunc("GET /indexes", handler.ListIndex)
	router.HandleFunc("GET /indexes/{index_name}", handler.GetIndexByName)
	router.HandleFunc("POST /describe_index_stats", handler.DescribeIndexStats)
	router.HandleFunc("POST /vectors/upsert", handler.UpsertVectors)
	router.HandleFunc("POST /query", handler.QueryVectors)
	router.HandleFunc("GET /vectors/fetch", handler.FetchVectors)
	router.HandleFunc("POST /vectors/update", handler.UpdateVector)
	router.HandleFunc("POST /vectors/delete", handler.DeleteVector)
	router.HandleFunc("GET /vectors/list", handler.ListVectorIDs)

	port := os.Getenv("PINECONE_MOCK_PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf("%s:%s", os.Getenv("PINECONE_MOCK_HOST"), port)
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	fmt.Println("Server listening on ", addr)
	server.ListenAndServe()
}
