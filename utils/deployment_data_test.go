package utils

import "testing"

func TestDeploymentDefault(t *testing.T) {
	const expected = "mongodb://localhost/test"
	d := &DeploymentData{
		DatabaseName:  "test",
		ServerAddress: "localhost",
	}
	if d.MakeDialString() != expected {
		t.Fatalf("expected %s, got %s", expected, d.MakeDialString())
	}
}
func TestDeploymentNoAuth(t *testing.T) {
	d := &DeploymentData{
		DatabaseName:  "test",
		ServerAddress: "localhost",
	}
	if d.IsAuthenticated() {
		t.Fatalf("should not be authenticated without username and password")
	}
}

func TestDeploymentWithAuth(t *testing.T) {
	// we have username and password, but we're not authenticated yet => admin
	const expected = "mongodb://test:123456@localhost/admin"
	d := &DeploymentData{
		DatabaseName:  "test",
		ServerAddress: "localhost",
		Username:      "test",
		Password:      "123456",
	}
	if d.MakeDialString() != expected {
		t.Fatalf("expected %s, got %s", expected, d.MakeDialString())
	}
}
