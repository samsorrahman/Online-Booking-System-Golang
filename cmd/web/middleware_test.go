package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var mh myHandler

	// Create a new request using http.NewRequest(). Pass in http.MethodGet as the
	h := NoSurf(&mh)

	switch v := h.(type) {
	case http.Handler:
		// We want to test that the middleware is working correctly, so we need to

	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var mh myHandler

	// Create a new request using http.NewRequest(). Pass in http.MethodGet as the
	h := SessionLoad(&mh)

	switch v := h.(type) {
	case http.Handler:
		// We want to test that the middleware is working correctly, so we need to

	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}
