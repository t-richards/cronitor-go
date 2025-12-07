package cronitor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/t-richards/cronitor-go"
)

const (
	urlPath        = "/p/your-api-key/nightly-job"
	defaultTimeout = 5 * time.Second
)

func TestCronitor_Run(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != urlPath {
			t.Fatalf("unexpected URL: %s", req.URL.Path)
		}

		if req.URL.RawQuery != "state=run" {
			t.Fatalf("unexpected query: %s", req.URL.RawQuery)
		}

		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	crn := cronitor.New(server.URL + urlPath)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := crn.Run(ctx)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
}

func TestCronitor_Complete(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != urlPath {
			t.Fatalf("unexpected URL: %s", req.URL.Path)
		}

		if req.URL.RawQuery != "state=complete" {
			t.Fatalf("unexpected query: %s", req.URL.RawQuery)
		}

		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	crn := cronitor.New(server.URL + urlPath)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := crn.Complete(ctx)
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}
}

func TestCronitor_Fail(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.URL.Path != urlPath {
			t.Fatalf("unexpected URL: %s", req.URL.Path)
		}

		if req.URL.RawQuery != "state=fail" {
			t.Fatalf("unexpected query: %s", req.URL.RawQuery)
		}

		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	crn := cronitor.New(server.URL + urlPath)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := crn.Fail(ctx)
	if err != nil {
		t.Fatalf("Fail failed: %v", err)
	}
}
