package paginator

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
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

func Test_PageListAndReturn(t *testing.T) {
	Convey("PageListAndReturn", t, func() {
		Convey("PageListAndReturn test1", func() {
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
			temp := PageListAndReturn(res, 11, 10)
			fmt.Println(temp)
		})
	})
}
