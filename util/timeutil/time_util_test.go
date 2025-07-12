package timeutil

import (
	"context"
	"example.com/m/test/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var ctx = testutil.NewContext(testutil.NewContextRequest{
	Host:   "",
	Cookie: "",
	Cid:    "ID",
})

func DateFormatStrToTimeStamp2(ctx context.Context, dateTimeStr string) int64 {
	dt, err := DateFormatStrToTime(ctx, dateTimeStr, DateFormatYYYYMMDDHHMMSS)
	if err != nil {
		return 0
	}
	return dt.Unix()
}

func Test_GetTimeRange(t *testing.T) {
	Convey("GetTimeRange", t, func() {
		Convey("GetTimeRange test1", func() {
			left, right, err := GetTimeRange(ctx, 2024, 6, 10)
			So(err, ShouldEqual, nil)
			So(left, ShouldEqual, DateFormatStrToTimeStamp2(ctx, "2024-06-10 00:00:00"))
			So(right, ShouldEqual, DateFormatStrToTimeStamp2(ctx, "2024-06-10 23:59:59"))

			left, right, err = GetTimeRange(ctx, 2024, 6, 13)
			So(err, ShouldEqual, nil)
			So(left, ShouldEqual, DateFormatStrToTimeStamp2(ctx, "2024-06-13 00:00:00"))
			So(right, ShouldEqual, DateFormatStrToTimeStamp2(ctx, "2024-06-13 23:59:59"))

			left, right, err = GetTimeRange(ctx, 2024, 6, 0)
			So(err, ShouldEqual, nil)
			So(left, ShouldEqual, DateFormatStrToTimeStamp2(ctx, "2024-06-01 00:00:00"))
			So(right, ShouldEqual, DateFormatStrToTimeStamp2(ctx, "2024-06-30 23:59:59"))

			left, right, err = GetTimeRange(ctx, 2024, 0, 0)
			So(err, ShouldEqual, nil)
			So(left, ShouldEqual, DateFormatStrToTimeStamp2(ctx, "2024-01-01 00:00:00"))
			So(right, ShouldEqual, DateFormatStrToTimeStamp2(ctx, "2024-12-31 23:59:59"))
		})
	})
}
