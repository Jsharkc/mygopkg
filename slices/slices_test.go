package slices_test

import (
	"fmt"
	"testing"

	"github.com/Jsharkc/mygopkg/slices"
)

func TestContainsInt(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	result := slices.Contains(intSlice, 7)
	fmt.Printf("results=%#v\n", result)
}

func TestContainsString(t *testing.T) {
	strSlice := []string{"1", "2", "3", "4", "5"}
	result := slices.Contains(strSlice, "3")
	fmt.Printf("results=%#v\n", result)
}

func TestJoinInts(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	result := slices.JoinInts(intSlice, "")
	fmt.Printf("results=%#v\n", result)
}

func TestMap(t *testing.T) {
	intSlice := []string{"1", "2", "3", "4", "5"}
	result := slices.Map(intSlice, func(x string) string { return x + "abc" })
	fmt.Printf("results=%#v\n", result)
}

func TestUnion(t *testing.T) {
	intSlice := []string{"1", "2", "3", "4", "5"}

	intSlice2 := []string{"1", "6", "8", "4", "10"}
	result := slices.Union(intSlice, intSlice2)
	fmt.Printf("results=%#v\n", result)
}

func TestDistinct(t *testing.T) {
	// stringSlice := []string{"1", "2", "3", "4", "5", "1", "6", "8", "4", "10"}

	intSlice := []int{1, 2, 3, 4, 5, 1, 6, 8, 4, 10}

	result := slices.Distinct(intSlice)
	fmt.Printf("results=%#v\n", result)
}

func TestDiff(t *testing.T) {
	intSlice := []string{"1", "2", "3", "4", "5"}

	intSlice2 := []string{"1", "6", "8", "4", "10"}
	result := slices.Diff(intSlice, intSlice2)
	fmt.Printf("results=%#v\n", result)
}

func TestInter(t *testing.T) {
	intSlice := []string{"1", "2", "3", "4", "5"}

	intSlice2 := []string{"1", "6", "8", "4", "10"}
	result := slices.Inter(intSlice, intSlice2)
	fmt.Printf("results=%#v\n", result)
}

func TestGetMapFromSlice(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	ps := []*Person{
		{"zhangsan", 30},
		{"lisi", 20},
		{"", 40},
		{"aa", 0},
	}

	result1 := slices.GetMapFromSlice[*Person, string](ps, "Name")
	result2 := slices.GetMapFromSlice[*Person, string](ps, "Name", false)
	result3 := slices.GetMapFromSlice[*Person, int](ps, "Age")
	result4 := slices.GetMapFromSlice[*Person, int](ps, "Age", false)
	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)
}

func TestGetSliceFromStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	ps := []*Person{
		{"zhangsan", 30},
		{"lisi", 20},
		{"", 40},
		{"aa", 0},
	}

	result1 := slices.GetSliceFromStructs[*Person, string](ps, "Name")
	result2 := slices.GetSliceFromStructs[*Person, string](ps, "Name", false)
	result3 := slices.GetSliceFromStructs[*Person, int](ps, "Age")
	result4 := slices.GetSliceFromStructs[*Person, int](ps, "Age", false)
	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)
}

func TestTransferIntSlice(t *testing.T) {
	type Nint8 int8
	res := slices.TransferIntSlice[int, Nint8]([]int{1, 2, 3})
	fmt.Println(res)
}

func TestGetMapSliceFromMapSlice(t *testing.T) {
	var source = []map[string]any{
		{
			"a": "2022-12-12",
			"b": 123,
			"c": map[any]any{
				"ca": "ccc",
				1:    123,
			},
		},
		{
			"a": "2022-12-13",
			"b": 133,
			"c": []int{1, 2, 3},
		},
		{
			"a": "2022-12-12",
			"b": 1234,
		},
	}
	res := slices.GetMapSliceFromMapSlice[string, any, string](source, "a")
	fmt.Println(res)
}

func TestSplitSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	res := slices.SplitSlice(s, 11)
	fmt.Println(res)
}


func TestGetIntsFromStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	ps := []*Person{
		{"zhangsan", 30},
		{"lisi", 20},
		{"", 40},
	}

	result1 := slices.GetIntsFromStructs(ps, "Age")
	result2 := slices.GetIntsFromStructs(ps, "Age", false)
	fmt.Println(result1)
	fmt.Println(result2)
}

func TestGetMapFromStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	ps := []*Person{
		{"zhangsan", 30},
		{"lisi", 20},
		{"", 40},
		{"aa", 0},
	}

	result1 := slices.GetMapFromStructs(ps, "Age")
	result2 := slices.GetMapFromStructs(ps, "Age", false)
	fmt.Println(result1)
	fmt.Println(result2)
}

func TestGetMapSliceFromStructs(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	ps := []*Person{
		{"zhangsan", 30},
		{"lisi", 20},
		{"", 40},
		{"aa", 0},
		{"bb", 0},
	}

	result1 := slices.GetMapSliceFromStructs(ps, "Age")
	result2 := slices.GetMapSliceFromStructs(ps, "Age", false)
	fmt.Println(result1)
	fmt.Println(result2)
}


func TestGetMapSliceFromSlice(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	ps := []*Person{
		{"zhangsan", 30},
		{"lisi", 20},
		{"", 40},
		{"aa", 0},
		{"", 0},
	}

	result1 := slices.GetMapSliceFromSlice[*Person, string](ps, "Name")
	result2 := slices.GetMapSliceFromSlice[*Person, string](ps, "Name", false)
	fmt.Println(result1)
	fmt.Println(result2)
}
