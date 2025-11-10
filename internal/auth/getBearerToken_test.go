package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken_LowerCase(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "bearer abc.def.ghi")
	_, err := GetBearerToken(h)
	if err == nil {
		t.Fatalf("error header not to spec: %v", err)
	}
}

func TestGetBearerToken_TrimSpaces(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "  Bearer    abc.def.ghi")
	got, err := GetBearerToken(h)
	if err != nil {
		t.Fatalf("unexptected error: %v", err)
	}
	if got != "abc.def.ghi" {
		t.Errorf("want abc.def.ghi, got %s", got)
	}
}
func TestGetBearerToken_ExtraSpaces(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer   abc.def.ghi")
	want := "abc.def.ghi"

	got, err := GetBearerToken(h)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestGetBearerToken_WrongScheme(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Basic abc.def.ghi")

	_, err := GetBearerToken(h)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestGetBearerToken_Success(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer abc.def.ghi")
	want := "abc.def.ghi"

	got, err := GetBearerToken(h)
	if err != nil {
		t.Errorf("unexptected error: %v", err)
	}
	if got != want {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestGetBearer_EmptyHeader(t *testing.T) {
	h := http.Header{}
	if _, err := GetBearerToken(h); err == nil {
		t.Fatalf("exptected error, got nil")
	}
}
