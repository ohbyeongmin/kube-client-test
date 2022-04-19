package client

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	k8utils "github.com/pytimer/k8sutil/apply"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type NAMESPACE string

const (
	TestFilePath = "testdata"
	Namespace    = "namespace.yaml"
	Deployment   = "deployment.yaml"
	Pvc          = "persistent-volume-claim.yaml"
	Service      = "service.yaml"
)

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
