package main

import (
	"errors"
	"fmt"
	"os"

	authorizationv1 "k8s.io/api/authorization/v1"
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

func (o *k8s) isVersion(major string, minor string) (bool, error) {
	version, err := o.clientset.Discovery().ServerVersion()
	if err != nil {
		return false, err
	}
	if version.Major != major {
		return false, errors.New("Major version does not match")
	}
	if version.Minor != minor {
		return false, errors.New("Minor version does not match")
	}
	return true, nil
}

func (o *k8s) canICreateDeployments() (bool, error) {
	ssar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Verb:     "create",
				Resource: "deployments",
			},
		},
	}
	ssar, err := o.clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(ssar)
	if err != nil {
		return false, err
	}
	return ssar.Status.Allowed, nil
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
	isV, err := k8s.isVersion("1", "9")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(isV)
}
