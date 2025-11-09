package generic

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Head[T any](list []T) (*T, bool) {
	if len(list) == 0 {
		return nil, false
	}
	return &list[0], true
}

func Head2[T int | string](list []T) (*T, bool) {
	if len(list) == 0 {
		return nil, false
	}
	return &list[0], true
}

func Test_generic(t *testing.T) {
	Convey("generic", t, func() {
		Convey("generic test1", func() {
			var list1 = []int{1, 2, 3}
			var list2 = []string{}
			fmt.Println(Head(list1))
			fmt.Println(Head(list2))
			fmt.Println(Head2(list1))
			fmt.Println(Head2(list2))
		})
	})
}
