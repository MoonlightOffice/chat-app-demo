package lib

import (
	"encoding/json"
	"testing"
)

func TestJWT(t *testing.T) {
	type Person struct {
		Name string  `json:"name"`
		Age  float64 `json:"age"`
	}

	p1 := Person{Name: "John", Age: 24}
	p1json, _ := json.Marshal(p1)

	token, err := ToJWT(p1json)
	if err != nil {
		t.Fatal(err)
	}

	p2json, err := FromJWT(token)
	if err != nil {
		t.Fatal(err)
	}

	var p2 Person
	err = json.Unmarshal(p2json, &p2)
	if err != nil {
		t.Fatal(err)
	}

	if !CompareStructs(p1, p2) {
		t.Fatal("p1 and p2 don't match")
	}
}
