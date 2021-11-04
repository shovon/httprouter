package httprouter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNothing(t *testing.T) {
	req, err := http.NewRequest("GET", "/anywhere", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := Router{}
	router.ServeHTTP(rr, req)

	str := rr.Body.String()
	if rr.Code != 404 {
		t.Errorf("Expected response code to be %d but got %d", 404, rr.Code)
	}

	if str != "Not found" {
		t.Errorf("Expected response body to be %s, but got %s", "Not found", str)
	}
}

func TestMethods(t *testing.T) {
	req, err := http.NewRequest("POST", "/somewhere", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := Router{}
	router.RegisterPathMethod(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Found"))
	}), "/somewhere", "POST")
	router.ServeHTTP(rr, req)

	str := rr.Body.String()
	if rr.Code != 200 {
		t.Errorf("Expected response code to be %d but got %d", 200, rr.Code)
	}

	if str != "Found" {
		t.Errorf("Expected response body to be %s, but got %s", "Found", str)
	}
}

func TestPathsPOST(t *testing.T) {
	req, err := http.NewRequest("POST", "/somewhere", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := Router{}
	router.RegisterPath(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Found"))
	}), "/somewhere")
	router.ServeHTTP(rr, req)

	str := rr.Body.String()
	if rr.Code != 200 {
		t.Errorf("Expected response code to be %d but got %d", 200, rr.Code)
	}

	if str != "Found" {
		t.Errorf("Expected response body to be %s, but got %s", "Found", str)
	}
}

func TestPathsGET(t *testing.T) {
	req, err := http.NewRequest("GET", "/somewhere", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := Router{}
	router.RegisterPath(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Found"))
	}), "/somewhere")
	router.ServeHTTP(rr, req)

	str := rr.Body.String()
	if rr.Code != 200 {
		t.Errorf("Expected response code to be %d but got %d", 200, rr.Code)
	}

	if str != "Found" {
		t.Errorf("Expected response body to be %s, but got %s", "Found", str)
	}
}

func TestDelegates(t *testing.T) {
	req, err := http.NewRequest("GET", "/foo/bar/baz", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := Router{}
	subRouter := Router{}
	subRouter.RegisterPathMethod(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Found"))
	}), "/bar/baz", "GET")
	router.Delegate(subRouter, "/foo")
	router.ServeHTTP(rr, req)

	str := rr.Body.String()
	if rr.Code != 200 {
		t.Errorf("Expected response code to be %d but got %d", 200, rr.Code)
	}

	if str != "Found" {
		t.Errorf("Expected response body to be %s, but got %s", "Found", str)
	}
}
