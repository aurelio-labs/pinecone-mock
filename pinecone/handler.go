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
			h.Index.CreateNamespace(vectorUpsertBody.Namespace)
		}
		h.Index.UpsertVector(vectorUpsertBody.Namespace, &vector)
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
		http.Error(w, fmt.Sprintf("Error reading request body %+v", r.Body), http.StatusInternalServerError)
		return
	}

	var vectorsResult []VectorResult
	vectors, err := h.Index.Query(vectorQueryQuery.Namespace, vectorQueryQuery.TopK)
	for _, value := range vectors {
		vectorsResult = append(vectorsResult, VectorResult{
			ID:       value.ID,
			Score:    0.9,
			Values:   value.Values,
			Metadata: value.Metadata,
		})
	}
	response := make(map[string]any)
	response["matches"] = vectorsResult
	response["namespace"] = vectorQueryQuery.Namespace
	response["usage"] = map[string]int{"read_units": len(vectors)}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
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

	vector, err := h.Index.GetVector(vectorUpdate.Namespace, vectorUpdate.ID)

	if vector == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	vector.Update(vectorUpdate)

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

	if (vectorDelete.IDs == nil) && (!vectorDelete.DeleteAll) {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	h.Index.DeleteVector(vectorDelete)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Pinecone returns empty map
	json.NewEncoder(w).Encode(make(map[string]string, 0))
}

func (h *Handler) ListVectorIDs(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")

	vectors, err := h.Index.ListVectorIDs(namespace)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := make(map[string]any)
	response["vectors"] = vectors
	response["namespace"] = namespace
	response["usage"] = map[string]int{"read_units": len(vectors)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
