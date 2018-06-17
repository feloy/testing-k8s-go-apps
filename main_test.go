package main

import (
	"errors"
	"testing"

	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/version"
	discoveryfake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func newTestSimpleK8s() *k8s {
	client := k8s{}
	client.clientset = fake.NewSimpleClientset()
	return &client
}

func newTestK8s() *k8s {
	client := k8s{
		clientset: &fake.Clientset{},
	}
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

func TestCanICreateDeploymentsFalse(t *testing.T) {
	k8s := newTestSimpleK8s()
	c, err := k8s.canICreateDeployments()
	if err != nil {
		t.Fatal("canICreateDeployments should not raise an error")
	}
	if c != false {
		t.Fatal("canICreateDeployments should return false")
	}
}

func TestCanICreateDeplaoymentsTrue(t *testing.T) {
	k8s := newTestK8s()
	k8s.clientset.(*fake.Clientset).Fake.AddReactor("create", "selfsubjectaccessreviews", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		mysar := &authorizationv1.SelfSubjectAccessReview{
			Status: authorizationv1.SubjectAccessReviewStatus{
				Allowed: true,
				Reason:  "I want to test it",
			},
		}
		return true, mysar, nil
	})
	c, err := k8s.canICreateDeployments()
	if err != nil {
		t.Fatal("canICreateDeployments should not raise an error")
	}
	if c != true {
		t.Fatal("canICreateDeployments should return true")
	}
}

func TestCanICreateDeplaoymentsError(t *testing.T) {
	k8s := newTestK8s()
	k8s.clientset.(*fake.Clientset).Fake.AddReactor("create", "selfsubjectaccessreviews", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &authorizationv1.SelfSubjectAccessReview{}, errors.New("Error creating ssar")
	})
	_, err := k8s.canICreateDeployments()
	if err == nil {
		t.Fatal("canICreateDeployments should raise an error")
	}
}
