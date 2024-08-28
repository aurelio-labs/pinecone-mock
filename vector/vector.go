package vector

type Vector struct {
	ID           string    `json:"id"`
	Values       []float64 `json:"values"`
	SparseValues struct {
		Indices []int64   `json:"indices"`
		Values  []float64 `json:"values"`
	} `json:"sparseValues,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type VectorResult struct {
	ID     string    `json:"id"`
	Score  float64   `json:"score"`
	Values []float64 `json:"values"`
}

type VectorFetchQuery struct {
	IDs       []string `json:"ids"`
	Namespace string   `json:"namespace"`
}

type VectorUpsertQuery struct {
	Vectors   []Vector `json:"vectors"`
	Namespace string   `json:"namespace"`
}
