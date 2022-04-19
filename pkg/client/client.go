package client

import (
	"context"
	"flag"
	"path/filepath"

	k8utils "github.com/pytimer/k8sutil/apply"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	client    *kubernetes.Clientset
	dynamic   dynamic.Interface
	discovery *discovery.DiscoveryClient

	workloads map[NAMESPACE]*Workload
}

func NewClient() *Client {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	dn, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &Client{
		client:    cs,
		dynamic:   dn,
		discovery: dc,

		workloads: make(map[NAMESPACE]*Workload),
	}
}

func (c *Client) createNamespace(namespace NAMESPACE) error {
	ns := apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: string(namespace),
		},
	}
	nc := c.client.CoreV1().Namespaces()
	_, err := nc.Create(context.TODO(), &ns, metav1.CreateOptions{})
	return err
}

func (c *Client) setNamespaceToWorkLoad(workload *Workload) {
	for _, u := range workload.unstructList {
		u.SetNamespace(string(workload.namespace))
	}
}

func (c *Client) AddWorkload(workload *Workload) {
	err := c.createNamespace(workload.namespace)
	if err != nil {
		panic(err)
	}

	c.setNamespaceToWorkLoad(workload)

	c.workloads[workload.namespace] = workload
}

func (c *Client) ApplyWorkload(namespace NAMESPACE) error {
	workload := c.workloads[namespace]
	applyOptions := k8utils.NewApplyOptions(c.dynamic, c.discovery)
	data, err := decodeWorkload(workload)
	if err != nil {
		return err
	}
	return applyOptions.Apply(context.TODO(), data)
}

func (c *Client) DeleteWorkload(namespace NAMESPACE) error {
	ctx := context.TODO()
	a := c.client.CoreV1().Namespaces()
	ns, err := a.Get(ctx, string(namespace), metav1.GetOptions{})
	if err != nil {
		return err
	}

	err = a.Delete(ctx, ns.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}
