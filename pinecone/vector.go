package pinecone

type SparseValues struct {
	Indices []int64   `json:"indices"`
	Values  []float64 `json:"values"`
}

type Vector struct {
	ID           string         `json:"id"`
	Values       []float64      `json:"values"`
	SparseValues *SparseValues  `json:"sparseValues,omitempty"`
	Metadata     map[string]any `json:"metadata"`
}

func (v *Vector) Update(updateQuery VectorUpdateQuery) {
	if updateQuery.Values != nil {
		v.Values = updateQuery.Values
	}
	if v.Metadata == nil {
		v.Metadata = updateQuery.SetMetadata
	} else {
		metadata := v.Metadata
		for key, value := range updateQuery.SetMetadata {
			metadata[key] = value
		}
		v.Metadata = metadata
	}
}

type VectorResult struct {
	ID       string         `json:"id"`
	Score    float64        `json:"score"`
	Values   []float64      `json:"values"`
	Metadata map[string]any `json:"metadata"`
}

type VectorUpsertQuery struct {
	Vectors   []Vector `json:"vectors"`
	Namespace string   `json:"namespace"`
}

type VectorQueryQuery struct {
	Namespace string `json:"namespace,omitempty"`
	TopK      int    `json:"topK"`
}

type VectorUpdateQuery struct {
	ID          string         `json:"id"`
	Values      []float64      `json:"values,omitempty"`
	SetMetadata map[string]any `json:"metadata,omitempty"`
	Namespace   string         `json:"namespace"`
}

type VectorDeleteQuery struct {
	IDs       []string       `json:"ids,omitempty"`
	DeleteAll bool           `json:"deleteAll"`
	Namespace string         `json:"namespace,omitempty"`
	Filter    map[string]any `json:"filter,omitempty"`
}
