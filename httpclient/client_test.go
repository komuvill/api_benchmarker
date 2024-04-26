package httpclient

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetRequestBody(t *testing.T) {
	t.Run("raw string", func(t *testing.T) {
		body := "raw body"
		reader, cleanup, err := GetRequestBody(body)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer cleanup()

		// Convert the reader back to a string for comparison
		buf := new(strings.Builder)
		_, err = io.Copy(buf, reader) // Use io.Copy instead of buf.ReadFrom
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := buf.String()
		if got != body {
			t.Errorf("expected body to be %q, got %q", body, got)
		}
	})

	t.Run("file path", func(t *testing.T) {
		// Create a temporary file with some content
		content := "file body"
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write([]byte(content)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		reader, cleanup, err := GetRequestBody("@" + tmpfile.Name())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer cleanup()

		buf := new(strings.Builder)
		_, err = io.Copy(buf, reader) // Use io.Copy here as well
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := buf.String()
		if got != content {
			t.Errorf("expected body to be %q, got %q", content, got)
		}
	})
}

func TestHttpRequest(t *testing.T) {
	// Create a test server that responds with a dummy response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ok"}`))
	}))
	defer ts.Close()

	t.Run("successful request", func(t *testing.T) {
		body, statusCode, err := HttpRequest("GET", ts.URL, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if statusCode != http.StatusOK {
			t.Errorf("expected status code to be %d, got %d", http.StatusOK, statusCode)
		}
		expectedBody := `{"message": "ok"}`
		if body != expectedBody {
			t.Errorf("expected body to be %q, got %q", expectedBody, body)
		}
	})
}
