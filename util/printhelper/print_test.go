package printhelper

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_PrintIntList(t *testing.T) {
	convey.Convey("PrintIntList", t, func() {
		convey.Convey("PrintIntList test1", func() {
			PrintIntList([]int{1, 3, 5, 7, 9, 11, 13, 15, 17, 99, 101, 202, 1001})
		})
		convey.Convey("Print test1", func() {
			Printf("test")
			Println("test")
		})
	})
}
