package client

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Workload struct {
	unstructList []unstructured.Unstructured
	namespace    NAMESPACE
}

func NewWorkload(namespace NAMESPACE, data []byte) *Workload {
	u, err := DecodeUnstructured(data)
	if err != nil {
		panic(err)
	}
	return &Workload{
		namespace:    namespace,
		unstructList: u,
	}
}

func decodeWorkload(workload *Workload) ([]byte, error) {
	var data []byte
	for _, u := range workload.unstructList {
		b, err := u.MarshalJSON()
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
		AddDivisionLine(data)
	}
	return data, nil
}
