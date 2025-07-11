package fileutil

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"testing"
)

type Staff struct {
	Name   string `json:"name"`
	Age    string `json:"age"`
	Size   string `json:"size"`
	Height string `json:"height"`
	Money  string `json:"money"`
	Temp   string `json:"temp"`
}

func Test_Save(t *testing.T) {
	Convey("Save", t, func() {
		Convey("Save test1", func() {

			var res []*Staff
			for i := 0; i < 100; i++ {
				res = append(res, &Staff{
					Name:   fmt.Sprintf("Name1-%v", i),
					Age:    fmt.Sprintf("Age1-%v", i),
					Size:   fmt.Sprintf("Size1-%v", i),
					Height: fmt.Sprintf("Height1-%v", i),
					Money:  fmt.Sprintf("Money1-%v", i),
					Temp:   fmt.Sprintf("Temp1-%v", i),
				})
			}
			err := SaveJson(&res, "", "staff1.json")
			So(err, ShouldEqual, nil)
			var temp []*Staff
			err = ReadJson(&temp, "staff1.json")
			So(err, ShouldEqual, nil)
			So(len(temp), ShouldEqual, len(res))
			sort.Slice(temp, func(i, j int) bool {
				return temp[i].Name < temp[j].Name
			})
			sort.Slice(res, func(i, j int) bool {
				return res[i].Name < res[j].Name
			})
			for i := 0; i < len(res); i++ {
				So(temp[i].Name, ShouldEqual, res[i].Name)
				So(temp[i].Age, ShouldEqual, res[i].Age)
				So(temp[i].Size, ShouldEqual, res[i].Size)
				So(temp[i].Height, ShouldEqual, res[i].Height)
				So(temp[i].Money, ShouldEqual, res[i].Money)
				So(temp[i].Temp, ShouldEqual, res[i].Temp)
			}

		})
	})
}
