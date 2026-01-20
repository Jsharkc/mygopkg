// Package slices provides common functions to operate or return slices.
// slices 包提供对切片的通用操作方法。
// 主要实现 Go 中缺少的对特定类型切片的操作。
package slices

import (
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/Jsharkc/mygopkg/constraints"
)

// JoinInts 使用给定分隔符，将整型切片拼接为单个字符串。
// 入参：1.用于拼接的整型切片 2.分隔符
// 返回：拼接结果(string)
func JoinInts(elems []int, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.Itoa(elems[0])
	}
	var b strings.Builder
	b.WriteString(strconv.Itoa(elems[0]))
	for _, elem := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(strconv.Itoa(elem))
	}
	return b.String()
}

// AppendAnySlice 将 []int、[]string 等展开后 append 到 []interface{} 中
// 返回最终的 []interface{}
func AppendAnySlice(interSlice []interface{}, inter interface{}) []interface{} {
	switch otherSlice := inter.(type) {
	case []int:
		for _, i := range otherSlice {
			interSlice = append(interSlice, i)
		}
	case []string:
		for _, s := range otherSlice {
			interSlice = append(interSlice, s)
		}
	}
	return interSlice
}

// Map 将Slice内元素逐个处理并将结果生成为新Slice
// 返回最终的 []R
func Map[T any, R any](t []T, predicate func(T) R) []R {
	r := make([]R, 0, len(t))
	for _, e := range t {
		r = append(r, predicate(e))
	}
	return r
}

// Filter 将Slice内元素逐个判断过滤并将结果生成为新Slice
// 返回最终的 []T
func Filter[T any](t []T, predicate func(T) bool) []T {
	r := make([]T, 0, len(t))
	for _, e := range t {
		if predicate(e) {
			r = append(r, e)
		}
	}
	return r
}

// Diff 获取列表t与列表r的差集
// 返回最终的 []T
func Diff[T comparable](t []T, r []T) []T {
	s := make(map[T]struct{})
	for _, e := range r {
		s[e] = struct{}{}
	}
	res := make([]T, 0, len(t))
	for _, e := range t {
		if _, ok := s[e]; !ok {
			res = append(res, e)
		}
	}
	return res
}

// Union 获取列表t与列表r的并集
// 返回最终的 []T
func Union[T comparable](t []T, r []T) []T {
	s := make(map[T]struct{})
	for _, e := range r {
		s[e] = struct{}{}
	}
	for _, e := range t {
		s[e] = struct{}{}
	}
	res := make([]T, 0, len(s))
	for e := range s {
		res = append(res, e)
	}
	return res
}

// Distinct 获取列表t去重后的新列表
// 返回最终的 []T
func Distinct[T comparable](t []T) []T {
	s := make(map[T]struct{})
	res := make([]T, 0, len(s))
	for _, e := range t {
		if _, ok := s[e]; !ok {
			s[e] = struct{}{}
			res = append(res, e)
		}
	}

	return res
}

// Inter 获取列表t与列表r的交集
// 返回最终的 []T
func Inter[T comparable](t []T, r []T) []T {
	s := make(map[T]struct{})
	for _, e := range r {
		s[e] = struct{}{}
	}
	res := make([]T, 0, len(t))
	for _, e := range t {
		if _, ok := s[e]; ok {
			res = append(res, e)
		}
	}
	return res
}

func Contains[T comparable](elems []T, i T) bool {
	for _, elem := range elems {
		if elem == i {
			return true
		}
	}
	return false
}

// GetIntsFromStructs 从 []struct 中的某个 struct 字段值提取出来成为 []int
// 如果 a 不是 struct slice，会 panic
//
// filterZeroVal ...bool 过滤零值，默认过滤。
func GetIntsFromStructs(a any, fieldname string, filterZeroVal ...bool) []int {
	value := reflect.ValueOf(a)

	length := value.Len()
	if length == 0 {
		return []int{}
	}

	filterZero := true
	if len(filterZeroVal) > 0 {
		filterZero = filterZeroVal[0]
	}

	intSlice := make([]int, 0, length)
	for i := 0; i < length; i++ {
		e := value.Index(i)
		if e.Kind() == reflect.Pointer {
			e = e.Elem()
		}
		val := e.FieldByName(fieldname)
		if filterZero && val.IsZero() {
			continue
		}
		intSlice = append(intSlice, int(val.Int()))
	}

	return intSlice
}

// GetSliceFromStructs 从 []struct 中的某个 struct 字段值提取出来成为 []S，其中 S 可以是任意可比较类型
// 如果 a 不是 struct slice，会 panic
//
// filterZeroVal ...bool 过滤零值，默认过滤。
func GetSliceFromStructs[T any, S comparable](slice []T, fieldname string, filterZeroVal ...bool) []S {
	length := len(slice)
	if length == 0 {
		return []S{}
	}

	filterZero := true
	if len(filterZeroVal) > 0 {
		filterZero = filterZeroVal[0]
	}

	newSlice := make([]S, 0, len(slice))
	for _, e := range slice {
		value := reflect.ValueOf(e)
		if value.Kind() == reflect.Pointer {
			value = value.Elem()
		}
		val := value.FieldByName(fieldname)
		if filterZero && val.IsZero() {
			continue
		}
		if k, ok := val.Interface().(S); ok {
			newSlice = append(newSlice, k)
		}
	}

	return newSlice
}

// GetMapFromSlice 从 []struct 中的某个 struct 字段值为 key，struct 为 value，组成 map
// 如果 slice 不是 struct slice，会 panic。支持任意 comparable 字段作为 key
//
// filterZeroVal ...bool 过滤零值，默认过滤。
func GetMapFromSlice[T any, S comparable](slice []T, fieldname string, filterZeroVal ...bool) map[S]T {
	length := len(slice)
	if length == 0 {
		return map[S]T{}
	}

	filterZero := true
	if len(filterZeroVal) > 0 {
		filterZero = filterZeroVal[0]
	}

	m := make(map[S]T)
	for _, e := range slice {
		value := reflect.ValueOf(e)
		if value.Kind() == reflect.Pointer {
			value = value.Elem()
		}
		val := value.FieldByName(fieldname)
		if filterZero && val.IsZero() {
			continue
		}
		if k, ok := val.Interface().(S); ok {
			m[k] = e
		}
	}

	return m
}

// GetMapFromStructs 从 []struct 中的某个 struct 字段值为 key，struct 为 value，组成 map
// 如果 slice 不是 struct slice，会 panic
// 只支持获取 int 为 key 的情况
//
// filterZeroVal ...bool 过滤零值，默认过滤。
func GetMapFromStructs[T any](slice []T, fieldname string, filterZeroVal ...bool) map[int]T {
	length := len(slice)
	if length == 0 {
		return map[int]T{}
	}

	filterZero := true
	if len(filterZeroVal) > 0 {
		filterZero = filterZeroVal[0]
	}

	m := make(map[int]T)
	for _, e := range slice {
		value := reflect.ValueOf(e)
		if value.Kind() == reflect.Pointer {
			value = value.Elem()
		}
		val := value.FieldByName(fieldname)
		if filterZero && val.IsZero() {
			continue
		}
		k := int(val.Int())
		m[k] = e
	}

	return m
}

// GetMapSliceFromStructs 从 []struct 中的某个 struct 字段值为 key，[]struct 为 value，组成 map
// 如果 slice 不是 struct slice，会 panic。只支持 int 作为 map 的 key
//
// filterZeroVal ...bool 过滤零值，默认过滤。
func GetMapSliceFromStructs[T any](slice []T, fieldname string, filterZeroVal ...bool) map[int][]T {
	length := len(slice)
	if length == 0 {
		return map[int][]T{}
	}

	filterZero := true
	if len(filterZeroVal) > 0 {
		filterZero = filterZeroVal[0]
	}

	m := make(map[int][]T)
	for _, e := range slice {
		value := reflect.ValueOf(e)
		if value.Kind() == reflect.Pointer {
			value = value.Elem()
		}
		val := value.FieldByName(fieldname)
		if filterZero && val.IsZero() {
			continue
		}
		k := int(val.Int())
		if tmpVal, ok := m[k]; ok {
			m[k] = append(tmpVal, e)
		} else {
			m[k] = []T{e}
		}
	}

	return m
}

// GetMapSliceFromSlice 从 []struct 中的某个 struct 字段值为 key，[]struct 为 value，组成 map
// 如果 slice 不是 struct slice，会 panic。支持任意可比较值作为 key。
//
// filterZeroVal ...bool 过滤零值，默认过滤。
func GetMapSliceFromSlice[T any, S comparable](slice []T, fieldname string, filterZeroVal ...bool) map[S][]T {
	length := len(slice)
	if length == 0 {
		return map[S][]T{}
	}
	filterZero := true
	if len(filterZeroVal) > 0 {
		filterZero = filterZeroVal[0]
	}

	m := make(map[S][]T)
	for _, e := range slice {
		value := reflect.ValueOf(e)
		if value.Kind() == reflect.Pointer {
			value = value.Elem()
		}
		val := value.FieldByName(fieldname)
		if filterZero && val.IsZero() {
			continue
		}
		if k, ok := val.Interface().(S); ok {
			if tmpVal, ok := m[k]; ok {
				m[k] = append(tmpVal, e)
			} else {
				m[k] = []T{e}
			}
		}
	}

	return m
}

// GetMapSliceFromMapSlice 从 map[comparable]any 中的某个key对应的val为 key，[]map[comparable]any 为 value，组成 map
// 支持任意可比较的值作为 key。
func GetMapSliceFromMapSlice[K comparable, V any, S comparable](slice []map[K]V, fieldname K) map[S][]map[K]V {
	length := len(slice)
	if length == 0 {
		return map[S][]map[K]V{}
	}

	m := make(map[S][]map[K]V)
	for _, e := range slice {
		if val, ok := e[fieldname]; ok {
			if realVal, ok := reflect.ValueOf(val).Interface().(S); ok {
				if tmpVal, ok := m[realVal]; ok {
					m[realVal] = append(tmpVal, e)
				} else {
					m[realVal] = []map[K]V{e}
				}
			}
		}
	}

	return m
}

// 整形切片格式互转，注意溢出问题
func TransferIntSlice[K, V constraints.Integer](source []K) (dest []V) {
	for _, v := range source {
		dest = append(dest, V(v))
	}
	return
}

// 将切片按照固定大小分割为子切片，返回结果为各子切片组成的切片
// 入参（1）list 要分割的源切片；（2）size 分割的固定大小
// 出参 resSlice分割后子切片组成的切片
// 示例 list : [1, 2, 3, 4, 5, 6, 7, 8, 9]  size : 4
// 结果 resSlice : [[1,2,3,4],[5,6,7,8],[9]]
func SplitSlice[T any](list []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	length := len(list)
	resLength := int(math.Ceil(float64(length) / float64(size)))
	spliltList := make([][]T, 0, resLength)
	for i := 0; i < resLength; i++ {
		tmpList := make([]T, 0, size)
		if i == int(resLength)-1 {
			tmpList = list[i*size:]
		} else {
			tmpList = list[i*size : i*size+size]
		}
		spliltList = append(spliltList, tmpList)
	}
	return spliltList
}
