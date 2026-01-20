package copier_test

import (
	"fmt"
	"testing"

	"github.com/Jsharkc/mygopkg/copier"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	type P struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		TreeIds []int  `json:"tree_ids"`
		HostIds []any  `json:"host_ids"`
	}

	p1 := &P{
		Name: "polairs",
		Age:  23,
	}
	p2 := &P{}

	err := copier.Copy(p2, p1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, p1.Name, p2.Name)
	assert.Equal(t, p1.Age, p2.Age)

	p := map[string]any{}

	m := map[string]any{
		"name":     "polaris",
		"age":      12,
		"tree_ids": []int{1, 2},
		"host_ids": []int{12, 1},
	}
	err = copier.Copy(&p, m)
	if err != nil {
		t.Error(err)
	}

	m["name"] = "xuxinhua"

	if p["name"] == "" {
		t.Error("failed")
	}
	fmt.Println(p)
}

func BenchmarkCopy(b *testing.B) {
	// p := &struct {
	// 	Name    string
	// 	Age     int
	// 	TreeIds []int `json:"tree_ids"`
	// 	HostIds []any `json:"host_ids"`
	// }{}
	p := map[string]any{}

	m := map[string]any{
		"name":     "polaris",
		"age":      12,
		"tree_ids": []int{1, 2},
		"host_ids": []int{12, 1},
	}

	for i := 0; i < b.N; i++ {
		_ = copier.Copy(&p, m)
	}
}
