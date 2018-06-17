package main

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type k8s struct {
	clientset kubernetes.Interface
}

func newK8s() (*k8s, error) {
	path := os.Getenv("HOME") + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, err
	}
	client := k8s{}
	client.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (o *k8s) getVersion() (string, error) {
	version, err := o.clientset.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", version), nil
}

func main() {
	k8s, err := newK8s()
	if err != nil {
		fmt.Println(err)
		return
	}
	v, err := k8s.getVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
}
