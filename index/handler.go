package index

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	Index *Index
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var index Index
	err := decoder.Decode(&index)
	if err != nil {
		panic(err)
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

	h.Index = &index
	log.Println(fmt.Sprintf("%+v", index))
	json.NewEncoder(w).Encode(index)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var indexResponse map[string][]Index
	if h.Index == nil {
		indexResponse = map[string][]Index{"indexes": []Index{}}
	} else {
		indexResponse = map[string][]Index{
			"indexes": []Index{*h.Index},
		}
	}
	json.NewEncoder(w).Encode(indexResponse)
}

func (h *Handler) GetByName(w http.ResponseWriter, r *http.Request) {
	if h.Index == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(h.Index)
}

func (h *Handler) DescribeIndexStats(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]int64{"dimension": h.Index.Dimension})
}
