package pinecone

import (
	"strings"
)

type IndexStatus struct {
	Ready bool   `json:"ready"`
	State string `json:"state"`
}

type Index struct {
	Name      string      `json:"name"`
	Dimension int         `json:"dimension"`
	Metric    string      `json:"metric"`
	Status    IndexStatus `json:"status"`
	Host      string      `json:"host"`
	Spec      struct {
	} `json:"spec"`
	Namespaces map[string]map[string]*Vector
}

func (i *Index) CreateNamespace(name string) {
	i.Namespaces[name] = make(map[string]*Vector)
}

func (i *Index) UpsertVector(namespace string, vector *Vector) {
	ns, exists := i.Namespaces[namespace]
	if !exists {
		ns = map[string]*Vector{vector.ID: vector}
	} else {
		ns[vector.ID] = vector
	}
	i.Namespaces[namespace] = ns
}

func (i *Index) Query(namespace string, topK int) ([]*Vector, error) {
	var vectors []*Vector
	remainingTopK := topK

	if namespace != "" {
		for _, vector := range i.Namespaces[namespace] {
			if remainingTopK > 0 {
				vectors = append(vectors, vector)
				remainingTopK--
			} else {
				break
			}
		}
	} else {
		for _, ns := range i.Namespaces {
			for _, vector := range ns {
				if remainingTopK > 0 {
					vectors = append(vectors, vector)
					remainingTopK--
				} else {
					break
				}

			}
		}
	}

	return vectors, nil
}

func (i *Index) Fetch(namespace string, ids []string) []*Vector {
	vectors := make([]*Vector, 0)

	if namespace == "" {
		for _, ns := range i.Namespaces {
			for _, id := range ids {
				vector, exists := ns[id]
				if exists {
					vectors = append(vectors, vector)
				}
			}
		}
	} else {
		ns, exists := i.Namespaces[namespace]
		if !exists {
			return vectors
		} else {
			for _, id := range ids {
				vector, exists := ns[id]
				if exists {
					vectors = append(vectors, vector)
				}
			}
		}
	}

	return vectors
}

func (i *Index) GetVector(namespace string, id string) *Vector {
	ns, exists := i.Namespaces[namespace]
	if !exists {
		return &Vector{}
	}
	vector, exists := ns[id]

	return vector
}

func (i *Index) DeleteVector(vectorDelete VectorDeleteQuery) {
	if vectorDelete.Namespace != "" {
		if vectorDelete.DeleteAll {
			i.Namespaces = make(map[string]map[string]*Vector)
		} else {
			for _, ns := range i.Namespaces {
				for _, id := range vectorDelete.IDs {
					delete(ns, id)
				}
			}
		}
	} else {
		if vectorDelete.DeleteAll {
			i.Namespaces[vectorDelete.Namespace] = make(map[string]*Vector)
		} else {
			for _, id := range vectorDelete.IDs {
				delete(i.Namespaces[vectorDelete.Namespace], id)
			}
		}
	}
}

func (i *Index) ListVectorIDs(namespace string, prefix string) []map[string]string {
	vectors := make([]map[string]string, 0)

	if namespace != "" {
		ns, exists := i.Namespaces[namespace]
		if !exists {
			return vectors
		}
		for id, _ := range ns {
			if strings.HasPrefix(id, prefix) {
				vectors = append(vectors, map[string]string{"id": id})
			}
		}
	} else {
		for _, ns := range i.Namespaces {
			for id, _ := range ns {
				if strings.HasPrefix(id, prefix) {
					vectors = append(vectors, map[string]string{"id": id})
				}
			}
		}

	}

	return vectors
}
