package hello_test

import (
	"openapi/internal/infra/env"
	oapicodegen "openapi/internal/infra/oapicodegen/hello"
	"testing"

	cmp "github.com/google/go-cmp/cmp"

	"encoding/json"
	"io"
	"net/http"
)

func TestGetSuccess(t *testing.T) {
	res, err := http.Get(env.GetServiceUrl() + "/hello")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, res.StatusCode)
	}

	resBodyByte, _ := io.ReadAll(res.Body)
	var actual = &oapicodegen.Hello{}
	json.Unmarshal(resBodyByte, &actual)

	var expect = &oapicodegen.Hello{
		Message: "Hello, World!",
	}

	if !cmp.Equal(actual, expect) {
		t.Errorf("expected %s, actual %s", expect, actual)
	}
}
