package vector

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	Vectors map[string]*Vector
}

func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var vectorUpsertBody VectorUpsertQuery
	err := decoder.Decode(&vectorUpsertBody)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	for _, vector := range vectorUpsertBody.Vectors {
		h.Vectors[vector.ID] = &vector
	}

	log.Println(fmt.Sprintf("%+v", vectorUpsertBody))
	json.NewEncoder(w).Encode(map[string]int{"upsertedCount": len(vectorUpsertBody.Vectors)})
}

func (h *Handler) Query(w http.ResponseWriter, r *http.Request) {
	var vectors []VectorResult
	for _, value := range h.Vectors {
		vectors = append(vectors, VectorResult{
			ID:     value.ID,
			Score:  0,
			Values: value.Values,
		})
	}
	response := make(map[string]any)
	response["matches"] = vectors
	response["namespace"] = "example-namespace"
	response["usage"] = map[string]int{"read_units": len(vectors)}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	vectors := make(map[string]Vector, 0)

	for _, id := range params["ids"] {
		vector, exists := h.Vectors[id]
		if exists {
			vectors[id] = *vector
		}
	}
	json.NewEncoder(w).Encode(map[string]any{
		"vectors":   vectors,
		"namespace": "example-namespace",
		"usage":     map[string]int{"readUnits": 1},
	})
}
