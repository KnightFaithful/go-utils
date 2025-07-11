package kafkaclient

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_singleConsumer(t *testing.T) {
	Convey("singleConsumer", t, func() {
		Convey("singleConsumer test1", func() {
			singleConsumer()
		})
	})
}
