package draw

import (
	. "github.com/smartystreets/goconvey/convey"
	"gonum.org/v1/plot/plotter"
	"testing"
)

func Test_DrawPoint(t *testing.T) {
	Convey("DrawPoint", t, func() {
		Convey("DrawPoint test1", func() {
			err := DrawPoint(plotter.XYs{
				{
					X: 1,
					Y: 2,
				},
				{
					X: 5,
					Y: 6,
				},
				{
					X: 86,
					Y: 34,
				},
				{
					X: 56,
					Y: 23,
				},
				{
					X: 12,
					Y: 32,
				},
			}, 200, 200)
			So(err, ShouldEqual, nil)

		})
	})
}
