package expression

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSplitStringToSlice(t *testing.T) {
	//src := &dto.SKUItem{SkuID: "aaa"}
	//def := &dto.SKUItem{SkuID: "bbb"}
	//dst := Default(src, def).(*dto.SKUItem)
	//fmt.Println("dst:", dst)
	//Assert(t, dst == src)
	//src = nil
	//dst = Default(src, def).(*dto.SKUItem)
	//Assert(t, dst == def)
	//dst = Default(nil, def).(*dto.SKUItem)
	//Assert(t, dst == def)
}

func Test_IfNilUseDefault(t *testing.T) {
	var num int64 = 100
	p := &num
	fmt.Println(IfNilUseDefaultInt64(p, 10))
	var q *int64
	fmt.Println(IfNilUseDefaultInt64(q, 10))
}

func TestAllMatchFunc(t *testing.T) {
	nums := []int64{0, 3}
	assertFunc := func(i int64) bool {
		return i <= 0
	}
	fmt.Println(AllMatch(assertFunc, nums...))
}

func Test_IsTheTimePeriodCrossed(t *testing.T) {
	Convey("IsTheTimePeriodCrossed", t, func() {
		Convey("IsTheTimePeriodCrossed test1", func() {
			So(IsTheTimePeriodCrossed(1000, 2000, 3000, 4000), ShouldEqual, false)
			So(IsTheTimePeriodCrossed(1000, 2000, 1000, 2000), ShouldEqual, true)
			So(IsTheTimePeriodCrossed(1000, 2000, 0, 2000), ShouldEqual, true)
			So(IsTheTimePeriodCrossed(1000, 2000, 3000, 4000), ShouldEqual, false)
			So(IsTheTimePeriodCrossed(10000, 40000, 20000, 30000), ShouldEqual, true)
			So(IsTheTimePeriodCrossed(20000, 30000, 10000, 40000), ShouldEqual, true)
			So(IsTheTimePeriodCrossed(1717779600+14*3600-4*3600, 1717779600+23*3600+4*3600, 1717866000+5*3600-4*3600, 1717866000+14*3600+4*3600), ShouldEqual, true)
		})
	})
}
