package pinecone

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	Index *Index
}

func (h *Handler) CreateIndex(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var index Index
	err := decoder.Decode(&index)
	if err != nil {
		http.Error(w, "Error creating index", http.StatusInternalServerError)
		return
	}

	if index.Metric == "" {
		index.Metric = "cosine"
	}
	if index.Status.State == "" {
		index.Status = IndexStatus{
			Ready: true,
			State: "Ready",
		}
	}
	index.Host = "http://localhost:8080"
	index.Namespaces = make(map[string]map[string]*Vector)

	h.Index = &index
	log.Println(fmt.Sprintf("Received index creation: %+v", index))

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(index)
}

func (h *Handler) ListIndex(w http.ResponseWriter, r *http.Request) {
	var indexResponse map[string][]Index
	if h.Index == nil {
		indexResponse = map[string][]Index{"indexes": []Index{}}
	} else {
		indexResponse = map[string][]Index{
			"indexes": []Index{*h.Index},
		}
	}
	log.Println("Received list index request")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(indexResponse)
}

func (h *Handler) GetIndexByName(w http.ResponseWriter, r *http.Request) {
	if h.Index == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println("Received get index request, returning index")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.Index)
}

func (h *Handler) DescribeIndexStats(w http.ResponseWriter, r *http.Request) {
	log.Println("Received describe index request")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	var total_vector_count int
	for _, vectors := range h.Index.Namespaces {
		total_vector_count += len(vectors)
	}
	json.NewEncoder(w).Encode(map[string]any{"dimension": h.Index.Dimension, "total_vector_count": total_vector_count})
}

func (h *Handler) UpsertVectors(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var vectorUpsertBody VectorUpsertQuery
	err := decoder.Decode(&vectorUpsertBody)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	for _, vector := range vectorUpsertBody.Vectors {
		if h.Index.Namespaces[vectorUpsertBody.Namespace] == nil {
			h.Index.Namespaces[vectorUpsertBody.Namespace] = make(map[string]*Vector)
		}
		h.Index.Namespaces[vectorUpsertBody.Namespace][vector.ID] = &vector
	}

	log.Println(fmt.Sprintf("%+v", vectorUpsertBody))
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"upsertedCount": len(vectorUpsertBody.Vectors)})
}

func (h *Handler) QueryVectors(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var vectorQueryQuery VectorQueryQuery
	err := decoder.Decode(&vectorQueryQuery)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var vectors []Vector
	remainingTopK := vectorQueryQuery.TopK
	if vectorQueryQuery.Namespace != "" {
		for _, vector := range h.Index.Namespaces[vectorQueryQuery.Namespace] {
			if remainingTopK > 0 {
				vectors = append(vectors, *vector)
				remainingTopK--
			} else {
				break
			}
		}
	} else {
		for _, ns := range h.Index.Namespaces {
			for _, vector := range ns {
				if remainingTopK > 0 {
					vectors = append(vectors, *vector)
					remainingTopK--
				} else {
					break
				}

			}
		}
	}

	var vectorsResult []VectorResult
	for _, value := range vectors {
		vectorsResult = append(vectorsResult, VectorResult{
			ID:     value.ID,
			Score:  0,
			Values: value.Values,
		})
	}
	response := make(map[string]any)
	response["matches"] = vectors
	response["namespace"] = vectorQueryQuery.Namespace
	response["usage"] = map[string]int{"read_units": len(vectors)}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) FetchVectors(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]any{
		"vectors":   []Vector{},
		"namespace": "not-implemented",
		"usage":     map[string]int{"readUnits": 1},
	})
}

func (h *Handler) UpdateVector(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var vectorUpdate VectorUpdateQuery
	err := decoder.Decode(&vectorUpdate)
	if err != nil {
		http.Error(w, "Error updating vector", http.StatusInternalServerError)
		return
	}

	vector, exists := h.Index.Namespaces[vectorUpdate.Namespace][vectorUpdate.ID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if vectorUpdate.Values != nil {
		vector.Values = vectorUpdate.Values
	}

	if vectorUpdate.SetMetadata != nil {
		for k, v := range vectorUpdate.SetMetadata {
			vector.Metadata[k] = v
		}
	}

	h.Index.Namespaces[vectorUpdate.Namespace][vectorUpdate.ID] = vector

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Pinecone returns empty map
	json.NewEncoder(w).Encode(make(map[string]string, 0))

}

func (h *Handler) DeleteVector(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var vectorDelete VectorDeleteQuery
	err := decoder.Decode(&vectorDelete)
	if err != nil {
		http.Error(w, "Error deleting vectors", http.StatusInternalServerError)
		return
	}

	if vectorDelete.Namespace != "" {
		if (vectorDelete.IDs == nil) && (!vectorDelete.DeleteAll) {
			http.Error(w, "Bad request", http.StatusBadRequest)
		} else if vectorDelete.DeleteAll {
			h.Index.Namespaces = make(map[string]map[string]*Vector)
		} else {
			for _, ns := range h.Index.Namespaces {
				for _, id := range vectorDelete.IDs {
					delete(ns, id)
				}
			}
		}
	} else {
		if (vectorDelete.IDs == nil) && (!vectorDelete.DeleteAll) {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		if vectorDelete.DeleteAll {
			h.Index.Namespaces[vectorDelete.Namespace] = make(map[string]*Vector)
		} else {
			for _, id := range vectorDelete.IDs {
				delete(h.Index.Namespaces[vectorDelete.Namespace], id)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Pinecone returns empty map
	json.NewEncoder(w).Encode(make(map[string]string, 0))
}

func (h *Handler) ListVectorIDs(w http.ResponseWriter, r *http.Request) {

	namespace := r.URL.Query().Get("namespace")

	vectors := make([]map[string]string, 0)
	for id, _ := range h.Index.Namespaces[namespace] {
		vectors = append(vectors, map[string]string{"id": id})
	}
	response := make(map[string]any)
	response["vectors"] = vectors
	response["namespace"] = namespace
	response["usage"] = map[string]int{"read_units": len(vectors)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
