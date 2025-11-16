package binarysearch

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_BinarySearchInt(t *testing.T) {
	Convey("BinarySearchInt", t, func() {
		Convey("BinarySearchInt test1", func() {
			list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			for i, n := range list {
				So(BinarySearchInt(list, n), ShouldEqual, i)
			}
			list = []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
			for i, n := range list {
				So(-BinarySearchInt(list, n-1)-1, ShouldEqual, i)
			}
			res := BinarySearchInt(list, 5)
			fmt.Println(res, -res-1)
		})
	})
}
