package client

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	k8utils "github.com/pytimer/k8sutil/apply"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type NAMESPACE string

func ReadYAMLFile(path string, name string) []byte {
	p := filepath.Join(path, name)
	b, err := ioutil.ReadFile(p)

	if err != nil {
		panic(err)
	}

	return AddDivisionLine(b)
}

func GetTestFileListToBytes(path string, files ...string) []byte {
	var fileList []byte
	for _, f := range files {
		fileList = append(fileList, ReadYAMLFile(path, f)...)
	}
	return fileList
}

func DecodeUnstructured(data []byte) ([]unstructured.Unstructured, error) {
	return k8utils.Decode(data)
}

func AddDivisionLine(data []byte) []byte {
	var byte_buf bytes.Buffer

	byte_buf.Write(data)
	byte_buf.WriteString("\n---\n")

	return byte_buf.Bytes()
}

func unstructuredToString(item *unstructured.Unstructured) (*string, error) {
	bytes, err := item.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var s = string(bytes)
	return &s, nil
}

func listToStrings(list *unstructured.UnstructuredList) ([]*string, error) {
	var ret []*string
	for _, item := range list.Items {
		s, err := unstructuredToString(&item)
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}
