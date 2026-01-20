package copier

import (
	"encoding/json"
	"reflect"
)

// Copy 支持各种符合类型的复制
func Copy(dst, src any) error {
	return CopyWithOption(dst, src, Option{})
}

func CopyWithOption(dst, src any, opt Option) error {
	dstKind := getKind(dst)
	srcKind := getKind(src)
	if (isMap(dstKind) && isStruct(srcKind)) ||
		(isStruct(dstKind) && isMap(srcKind)) {
		return copyByJSON(dst, src)
	}

	return copyWithOption(dst, src, opt)
}

func copyByJSON(dst, src any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

func getKind(v any) reflect.Kind {
	reflectValue := reflect.ValueOf(v)
	return indirect(reflectValue).Kind()
}

func isMap(k reflect.Kind) bool {
	return k == reflect.Map
}

func isStruct(k reflect.Kind) bool {
	return k == reflect.Struct
}
