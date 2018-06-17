package main

import (
	"testing"

	"k8s.io/apimachinery/pkg/version"
	discoveryfake "k8s.io/client-go/discovery/fake"
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

func TestIsVersionOK(t *testing.T) {
	expectedMajor := "1"
	expectedMinor := "9"
	k8s := newTestSimpleK8s()
	k8s.clientset.Discovery().(*discoveryfake.FakeDiscovery).FakedServerVersion = &version.Info{
		Major: expectedMajor,
		Minor: expectedMinor,
	}
	v, err := k8s.isVersion(expectedMajor, expectedMinor)
	if err != nil {
		t.Fatal("isVersion should not raise an error")
	}
	if v != true {
		t.Fatal("isVersion should return true")
	}
}

func TestIsVersionErrorMajor(t *testing.T) {
	expectedMajor := "1"
	expectedMinor := "9"
	k8s := newTestSimpleK8s()
	k8s.clientset.Discovery().(*discoveryfake.FakeDiscovery).FakedServerVersion = &version.Info{
		Major: "wrong",
		Minor: "wrong",
	}
	_, err := k8s.isVersion(expectedMajor, expectedMinor)
	if err == nil {
		t.Fatal("isVersion should raise an error")
	}
	expected := "Major version does not match"
	if err.Error() != expected {
		t.Fatal("Raised error should be: " + expected)
	}
}

func TestIsVersionErrorMinor(t *testing.T) {
	expectedMajor := "1"
	expectedMinor := "9"
	k8s := newTestSimpleK8s()
	k8s.clientset.Discovery().(*discoveryfake.FakeDiscovery).FakedServerVersion = &version.Info{
		Major: expectedMajor,
		Minor: "wrong",
	}
	_, err := k8s.isVersion(expectedMajor, expectedMinor)
	if err == nil {
		t.Fatal("isVersion should raise an error")
	}
	expected := "Minor version does not match"
	if err.Error() != expected {
		t.Fatal("Raised error should be: " + expected)
	}
}
