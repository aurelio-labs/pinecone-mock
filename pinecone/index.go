package pinecone

import (
	"errors"
	"fmt"
)

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

func (i *Index) CreateNamespace(name string) {
	i.Namespaces[name] = make(map[string]*Vector)
}

func (i *Index) UpsertVector(namespace string, vector *Vector) error {
	ns, exists := i.Namespaces[namespace]
	var err error
	if !exists {
		err = errors.New(fmt.Sprintf("Namespace not found: %s", namespace))
	} else {
		ns[vector.ID] = vector
	}

	return err
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

func (i *Index) GetVector(namespace string, id string) (*Vector, error) {
	ns, exists := i.Namespaces[namespace]
	if !exists {
		err := errors.New(fmt.Sprintf("Namespace not found: %s", namespace))
		return nil, err
	}
	vector, exists := ns[id]
	if !exists {
		return nil, nil
	}

	return vector, nil

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

func (i *Index) ListVectorIDs(namespace string) ([]map[string]string, error) {
	vectors := make([]map[string]string, 0)
	if namespace != "" {
		ns, exists := i.Namespaces[namespace]
		if !exists {
			return nil, errors.New(fmt.Sprintf("Namespace %s does not exist"))
		}
		for id, _ := range ns {
			vectors = append(vectors, map[string]string{"id": id})
		}

	} else {
		for _, ns := range i.Namespaces {
			for id, _ := range ns {
				vectors = append(vectors, map[string]string{"id": id})
			}
		}
	}

	return vectors, nil
}
