package index

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
}
