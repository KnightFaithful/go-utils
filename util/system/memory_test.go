package system

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_MemStat(t *testing.T) {
	Convey("MemStat", t, func() {
		Convey("MemStat test1", func() {
			MemStat()
			list := make([]int64, 1024*1024*1024)
			list[0] = 1
			MemStat()
			list[1] = 2
			list[2] = 3
		})
	})
}
