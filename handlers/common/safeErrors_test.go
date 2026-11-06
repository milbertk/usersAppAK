package common

import (
	"errors"
	"testing"
)

func TestClientMessage(t *testing.T) {
	cases := map[string]string{
		"the company name exists, please choose another":                      "the company name exists, please choose another",
		"❌ you don't have permission to update this company":                  "you don't have permission to update this company",
		"pq: duplicate key value violates unique constraint \"companies_pk\"": "fallback",
		"❌ DB connection error: dial tcp 10.0.0.1:5432: i/o timeout":          "fallback",
		"error validating company id: sql: no rows in result set":             "fallback",
		"firebase init error: cannot read credentials":                        "fallback",
	}
	for in, want := range cases {
		if got := ClientMessage(errors.New(in), "fallback"); got != want {
			t.Errorf("ClientMessage(%q) = %q, want %q", in, got, want)
		}
	}
	if got := ClientMessage(nil, "fallback"); got != "fallback" {
		t.Errorf("nil error: got %q", got)
	}
}
