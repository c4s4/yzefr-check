package main

import (
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	var err error
	if err = run([]string{"test/nominal.md"}); err != nil {
		t.Fatalf("error checking file: %v", err)
	}
	if err = run([]string{"test/unknown.md"}); err == nil {
		t.Fatal("should have failed")
	}
	if !strings.Contains(err.Error(), "field foo not found in type main.Header") {
		t.Fatalf("bad error message: %s", err.Error())
	}
	if err = run([]string{"test/missing.md"}); err == nil {
		t.Fatal("should have failed")
	}
	if !strings.Contains(err.Error(), "Field validation for 'Title' failed on the 'required' tag") {
		t.Fatalf("bad error message: %s", err.Error())
	}
}
