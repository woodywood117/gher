package gher

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Person struct {
	Name string `json:"name"`
}

func Hello(p Person, _ *http.Request) (string, error) {
	return fmt.Sprintf("Hello, %s!", p.Name), nil
}

func HelloErr(p Person, _ *http.Request) (string, error) {
	return fmt.Sprintf("Hello, %s!", p.Name), fmt.Errorf("something went wrong")
}

func HelloReader(p Person, _ *http.Request) (*bytes.Buffer, error) {
	return bytes.NewBufferString(fmt.Sprintf("Hello, %s!", p.Name)), nil
}

func HelloStruct(p Person, _ *http.Request) (Person, error) {
	return p, nil
}

func HelloEmpty(_ struct{}, _ *http.Request) (string, error) {
	return "empty", nil
}

func HelloInString(name string, _ *http.Request) (string, error) {
	return fmt.Sprintf("Hello, %s!", name), nil
}

func HelloInPointer(p *Person, _ *http.Request) (string, error) {
	if p == nil {
		return "", fmt.Errorf("nil pointer")
	}
	return fmt.Sprintf("Hello, %s!", p.Name), nil
}

func TestGherSuccess(t *testing.T) {
	server := httptest.NewServer(Gher(Hello))

	body := bytes.NewBufferString(`{"name":"World"}`)
	resp, err := http.Post(server.URL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	output := new(bytes.Buffer)
	_, err = output.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if output.String() != "Hello, World!" {
		t.Fatalf("expected %q, got %q", "Hello, World!", output.String())
	}
}

func TestGherParseError(t *testing.T) {
	server := httptest.NewServer(Gher(Hello))

	body := bytes.NewBufferString(`invalid json`)
	resp, err := http.Post(server.URL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGherHandlerError(t *testing.T) {
	server := httptest.NewServer(Gher(HelloErr))

	body := bytes.NewBufferString(`{"name":"World"}`)
	resp, err := http.Post(server.URL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

	output := new(bytes.Buffer)
	_, err = output.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if output.String() != `{"error": "something went wrong"}` {
		t.Fatalf("expected %q, got %q", `{"error": "something went wrong"}`, output.String())
	}
}

func TestGherReader(t *testing.T) {
	server := httptest.NewServer(Gher(HelloReader))

	body := bytes.NewBufferString(`{"name":"World"}`)
	resp, err := http.Post(server.URL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	output := new(bytes.Buffer)
	_, err = output.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if output.String() != "Hello, World!" {
		t.Fatalf("expected %q, got %q", "Hello, World!", output.String())
	}
}

func TestGherStruct(t *testing.T) {
	server := httptest.NewServer(Gher(HelloStruct))

	body := bytes.NewBufferString(`{"name":"World"}`)
	resp, err := http.Post(server.URL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	output := new(bytes.Buffer)
	_, err = output.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if output.String() != "{\"name\":\"World\"}\n" {
		t.Fatalf("expected %q, got %q", "{\"name\":\"World\"}\n", output.String())
	}
}

func TestGherEmpty(t *testing.T) {
	server := httptest.NewServer(Gher(HelloEmpty))

	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	output := new(bytes.Buffer)
	_, err = output.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if output.String() != "empty" {
		t.Fatalf("expected %q, got %q", "empty", output.String())
	}
}

func TestGherInString(t *testing.T) {
	server := httptest.NewServer(Gher(HelloInString))

	body := bytes.NewBufferString("World")
	resp, err := http.Post(server.URL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	output := new(bytes.Buffer)
	_, err = output.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if output.String() != "Hello, World!" {
		t.Fatalf("expected %q, got %q", "Hello, World!", output.String())
	}
}

func TestGherInPointer(t *testing.T) {
	server := httptest.NewServer(Gher(HelloInPointer))

	body := bytes.NewBufferString(`{"name":"World"}`)
	resp, err := http.Post(server.URL, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	output := new(bytes.Buffer)
	_, err = output.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if output.String() != "Hello, World!" {
		t.Fatalf("expected %q, got %q", "Hello, World!", output.String())
	}
}

func TestGherInPointerNil(t *testing.T) {
	server := httptest.NewServer(Gher(HelloInPointer))

	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

	output := new(bytes.Buffer)
	_, err = output.ReadFrom(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if output.String() != `{"error": "nil pointer"}` {
		t.Fatalf("expected %q, got %q", `{"error": "nil pointer"}`, output.String())
	}
}
