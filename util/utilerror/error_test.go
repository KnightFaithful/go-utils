package utilerror

import (
	"example.com/m/util/printhelper"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func a() *UtilError {
	err := b()
	if err != nil {
		return err.AddError("")
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
	return NewError("test")
}

func Test_err(t *testing.T) {
	Convey("err", t, func() {
		Convey("err test1", func() {
			err := a()
			So(err, ShouldNotEqual, nil)
			//err.DebugError()
			fmt.Println(err.DebugError())
		})
		Convey("err test2", func() {
			GetRes()
		})
	})
}

func GetRes() (res bool) {
	res = true
	defer func() {
		printhelper.Println(res)
	}()
	return false
}
