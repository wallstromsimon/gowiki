package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPageRoundTrip(t *testing.T) {
	t.Chdir(t.TempDir())

	p := new(page{Title: "FrontPage", Body: "hello wiki"})
	if err := p.save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := loadPage("FrontPage")
	if err != nil {
		t.Fatalf("loadPage: %v", err)
	}
	if got.Title != p.Title || got.Body != p.Body {
		t.Fatalf("got %+v, want %+v", got, p)
	}
}

func TestRejectInvalidTitle(t *testing.T) {
	t.Chdir(t.TempDir())

	p := new(page{Title: "../etc", Body: "nope"})
	if err := p.save(); err != os.ErrInvalid {
		t.Fatalf("save error = %v, want %v", err, os.ErrInvalid)
	}

	if _, err := loadPage("bad-title"); err != os.ErrInvalid {
		t.Fatalf("loadPage error = %v, want %v", err, os.ErrInvalid)
	}
}

func TestHandlersRejectInvalidTitle(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/view/bad-title", nil)
	req.SetPathValue("title", "bad-title")
	rr := httptest.NewRecorder()
	viewHandler(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("view status = %d, want %d", rr.Code, http.StatusNotFound)
	}

	req = httptest.NewRequest(http.MethodPost, "/save/bad-title", nil)
	req.SetPathValue("title", "bad-title")
	rr = httptest.NewRecorder()
	saveHandler(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("save status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestViewRedirectsMissingPage(t *testing.T) {
	t.Chdir(t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/view/MissingPage", nil)
	req.SetPathValue("title", "MissingPage")
	rr := httptest.NewRecorder()
	viewHandler(rr, req)
	if rr.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusFound)
	}
	if loc := rr.Header().Get("Location"); loc != "/edit/MissingPage" {
		t.Fatalf("Location = %q, want /edit/MissingPage", loc)
	}
}
