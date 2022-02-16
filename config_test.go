package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestConfig(t *testing.T) {
	want := []byte(`{"targets":[{"ip":"10.0.0.10","port":"8080"}]}`)

	c := Parse("test_data/config.json")
	got, _ := json.Marshal(c)

	if !bytes.Equal(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
