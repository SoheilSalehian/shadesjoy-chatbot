package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSanityMessengerVerifyHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/webhook", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MessengerVerify)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v was expecting %v", status, http.StatusOK)
	}
}
