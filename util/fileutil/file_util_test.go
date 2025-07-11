package fileutil

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_ReadFolder(t *testing.T) {
	Convey("ReadFolder", t, func() {
		Convey("ReadFolder test1", func() {
			fileNameList, err := ReadFolder("../")
			So(err, ShouldEqual, nil)
			//[../ ../collection ../collection/set.go ../collection/string_set.go ../convert ../convert/convert_helper.go ../fileutil ../fileutil/file_util.go ../fileutil/file_util_test.go ../httputil ../httputil/post.go ../stringutil ../stringutil/string_util.go ../stringutil/string_util_test.go ../utilerror ../utilerror/error_code.go ../utilerror/util_error.go]
			fmt.Println(fileNameList)
		})
	})
}

func Test_FileExists(t *testing.T) {
	Convey("FileExists", t, func() {
		Convey("FileExists test1", func() {
			exists := FileExists("test")
			So(exists, ShouldEqual, true)

		})
	})
}

func Test_WriteString(t *testing.T) {
	Convey("WriteString", t, func() {
		Convey("WriteString test1", func() {
			err := WriteString("hello world!", "test.txt")
			So(err, ShouldEqual, nil)
		})
		Convey("WriteString test2", func() {
			err := WriteString("hello world2!", "test/test.txt")
			So(err, ShouldEqual, nil)
		})
	})
}

func Test_createFolderIfNotExist(t *testing.T) {
	Convey("CreateFolderIfNotExist", t, func() {
		Convey("CreateFolderIfNotExist test1", func() {
			err := CreateFolderIfNotExist("test/test")
			So(err, ShouldEqual, nil)

		})
		Convey("CreateFolderIfNotExist test2", func() {
			err := CreateFolderIfNotExist("test1/test1/test.go")
			So(err, ShouldEqual, nil)

		})
	})
}

func Test_DeletePath(t *testing.T) {
	Convey("DeletePath", t, func() {
		Convey("DeletePath test1", func() {
			err := DeletePath("test")
			So(err, ShouldEqual, nil)

		})
	})
}
