package backson

import (
	"fmt"
	"reflect"
	"testing"
)

type SpecialInt64 int64
type SpecialInt64Ptr *int64

func Test_typeAssert(t *testing.T) {
	typeAssert(1)
	typeAssert(SpecialInt64(2))
	typeAssert(new(SpecialInt64Ptr))
}

func typeAssert[T any](e T) {
	t := reflect.TypeOf(e)
	fmt.Println("type is:", t, " value is: ", e, " kind is: ", reflect.ValueOf(e).Kind())
	reflect.TypeFor[T]()
}
