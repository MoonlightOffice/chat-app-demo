package lib

import (
	"testing"
)

func TestTrim(t *testing.T) {
	cases := map[string]string{
		" asdf ":          "asdf",
		"as df ":          "as df",
		"\tas df\n":       "as df",
		"\n \t as df \n ": "as df",
	}

	for before, after := range cases {
		trimmed := Trim(before)
		if trimmed != after {
			t.Fatalf("Expected %s, got %s\n", after, trimmed)
		}
	}
}

func TestCompareStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p1 := Person{
		Name: "Ichigo",
		Age:  22,
	}
	p2 := Person{
		Name: "Ichigo",
		Age:  21,
	}

	if CompareStructs(p1, p2) {
		t.Fatal("Must be false because the age is different")
	}

	p2.Age = 22
	if !CompareStructs(p1, p2) {
		t.Fatal("Invalid")
	}
}
