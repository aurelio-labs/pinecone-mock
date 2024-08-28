package pinecone

type IndexStatus struct {
	Ready bool   `json:"ready"`
	State string `json:"state"`
}

type Index struct {
	Name      string      `json:"name"`
	Dimension int64       `json:"dimension"`
	Metric    string      `json:"metric"`
	Status    IndexStatus `json:"status"`
	Host      string      `json:"host"`
	Spec      struct {
	} `json:"spec"`
	Namespaces map[string]map[string]*Vector
}

type Vector struct {
	ID           string    `json:"id"`
	Values       []float64 `json:"values"`
	SparseValues struct {
		Indices []int64   `json:"indices,omitempty"`
		Values  []float64 `json:"values,omitempty"`
	} `json:"sparseValues,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type VectorResult struct {
	ID     string    `json:"id"`
	Score  float64   `json:"score"`
	Values []float64 `json:"values"`
}

type VectorUpsertQuery struct {
	Vectors   []Vector `json:"vectors"`
	Namespace string   `json:"namespace"`
}

type VectorQueryQuery struct {
	Namespace       string         `json:"namespace"`
	TopK            int            `json:"topK"`
	Filter          map[string]any `json:"filter"`
	IncludeValues   bool           `json:"includeValues"`
	IncludeMetadata bool           `json:"includeMetadata"`
	Vector          []float64      `json:"vector"`
	ID              string         `json:"id"`
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
