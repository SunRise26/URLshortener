package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMainHandler(t *testing.T) {
	handler := http.HandlerFunc(mainHandler)
	serv.db = OpenDatabase()

	// test index page
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `<title>Bluefield test</title>`
	if strings.Contains(rr.Body.String(), expected) == false {
		t.Errorf("handler returned unexpected body: Doesn't contain %s",
			expected)
	}

	// test StatusNotFound
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/testNotFound", nil)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// test StatusFound
	rr = httptest.NewRecorder()
	serv.db.addLink("https://google.com", "testFound")
	defer serv.db.deleteLink("testFound")
	req = httptest.NewRequest("GET", "/testFound", nil)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}
}

func TestCreateLink(t *testing.T) {
	handler := http.HandlerFunc(createLink)
	serv.db = OpenDatabase()

	// test empty request
	rr := httptest.NewRecorder()
	reader := strings.NewReader("origURL=&shortURL=")
	req := httptest.NewRequest(http.MethodPost, "/create", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected := "Error: Empty request\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body:\nexpected:\n%s\ngot:\n%s",
			expected, rr.Body.String())
	}

	// test wrong custom link
	rr = httptest.NewRecorder()
	reader = strings.NewReader("origURL=https://www.google.com/&shortURL=testFound")
	req = httptest.NewRequest(http.MethodPost, "/create", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	expected = "Error: Shortened URL 'testFound' already exist\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body:\nexpected:\n%s\ngot:\n%s",
			expected, rr.Body.String())
	}
}

func TestShowLinks(t *testing.T) {
	handler := http.HandlerFunc(showLinks)
	serv.db = OpenDatabase()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/showLinks", nil)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := serv.db.getLinks() + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}

func TestShowRequests(t *testing.T) {
	handler := http.HandlerFunc(showRequests)
	serv.db = OpenDatabase()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/showRequests", nil)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := serv.db.getRequests() + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body")
	}
}
