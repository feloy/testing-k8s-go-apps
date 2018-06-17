package main

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

func newTestSimpleK8s() *k8s {
	client := k8s{}
	client.clientset = fake.NewSimpleClientset()
	return &client
}

func TestGetVersionDefault(t *testing.T) {
	k8s := newTestSimpleK8s()
	v, err := k8s.getVersion()
	if err != nil {
		t.Fatal("getVersion should not raise an error")
	}
	expected := "v0.0.0-master+$Format:%h$"
	if v != expected {
		t.Fatal("getVersion should return " + expected)
	}
}
