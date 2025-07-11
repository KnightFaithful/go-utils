package task

import (
	"example.com/m/util/utilerror"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func Test_ConcurrentQueryTaskRunnerWithRoutineCount(t *testing.T) {
	Convey("ConcurrentQueryTaskRunnerWithRoutineCount", t, func() {
		Convey("ConcurrentQueryTaskRunnerWithRoutineCount test1", func() {
			res, err := ConcurrentQueryTaskRunnerWithRoutineCount([]int64{1, 2, 3, 4, 5, 6}, 1, func(i interface{}) (interface{}, *utilerror.UtilError) {
				fmt.Println(i)
				time.Sleep(2 * time.Second)
				return nil, nil
			}, 2, false)
			So(err, ShouldEqual, nil)
			fmt.Println(res)
		})
	})
}
