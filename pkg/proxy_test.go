package pkg

import (
	"testing"
)

func TestNewReverseProxy(t *testing.T) {
	proxy := NewReverseProxy("http", "localhost", "example.com", 8080)
	if proxy == nil {
		t.Fatal("Expected non-nil proxy")
	}
	if proxy.BaseDomain != "example.com" {
		t.Errorf("Expected BaseDomain to be 'example.com', got '%s'", proxy.BaseDomain)
	}
}
