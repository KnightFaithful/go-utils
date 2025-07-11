package utilerror2

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func a() *UtilError {
	err := b()
	if err != nil {
		return err.AddError(-2, "")
	}
	return nil
}

func b() *UtilError {
	err := c()
	if err != nil {
		return err.Mark()
	}
	return nil
}

func c() *UtilError {
	return NewError(-1, "test")
}

func Test_err(t *testing.T) {
	Convey("err", t, func() {
		Convey("err test1", func() {
			err := a()
			So(err, ShouldNotEqual, nil)
			//err.DebugError()
			fmt.Println(err.DebugError())
		})
	})
}
