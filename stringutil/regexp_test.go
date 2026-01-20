package stringutil_test

import (
	"fmt"
	"testing"

	"github.com/Jsharkc/mygopkg/stringutil"
)

func TestPhoneMask(t *testing.T) {
	input := "我的手机号是18434858803,1111, +86 18434858802"
	res := stringutil.FindAllPhone(input, -1)
	fmt.Printf("%+v", res)
}

func TestPhone(t *testing.T) {
	phone := "13456780901"
	res := stringutil.CheckPhone(phone)
	fmt.Println(res)
}

func TestEmail(t *testing.T) {
	email := "123@test.com"
	res := stringutil.CheckEmail(email)
	fmt.Println(res)
}
