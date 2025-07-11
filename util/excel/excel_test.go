package excel

import (
	"example.com/m/util/fileutil"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type ExcelItem struct {
	EventId    string `json:"event_id" title:"Event Id"`
	BizStaffId string `json:"biz_staff_id" title:"BizStaffId"`
}

func Test_GenerateExcelBytes(t *testing.T) {
	Convey("GenerateExcelBytes", t, func() {
		Convey("GenerateExcelBytes test1", func() {
			fileName := "test.xlsx"
			byteList, _, err := GenerateExcelBytes(fileName, []*ExcelSheetTab{
				{
					SheetName: "Sheet1",
					Data: []*ExcelItem{
						{
							EventId:    "event1",
							BizStaffId: "Ops1",
						},
						{
							EventId:    "Event2",
							BizStaffId: "Ops2",
						},
					},
					ExcludeTitles: nil,
				},
			})
			So(err, ShouldEqual, nil)
			err = fileutil.WriteByteList(byteList, fileName)
			So(err, ShouldEqual, nil)
		})
	})
}
