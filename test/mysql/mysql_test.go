package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"example.com/m/test/testutil"
	"example.com/m/util/stringutil"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const SocManpowerShiftForecastTabName = "soc_manpower_shift_forecast_tab"

type SocManpowerShiftForecast struct {
	ID           int64        `gorm:"column:id;primaryKey" json:"id"`
	StationID    int64        `gorm:"column:station_id" json:"station_id"`
	Operator     string       `gorm:"column:operator" json:"operator"`
	ForecastDate int64        `gorm:"column:forecast_date" json:"forecast_date"` // 使用更友好的时间类型
	SolutionType int64        `gorm:"column:solution_type" json:"solution_type"`
	Version      *VersionInfo `gorm:"column:version" json:"version"` // JSON 文本
	Ctime        int64        `gorm:"column:ctime" json:"ctime"`
	Mtime        int64        `gorm:"column:mtime" json:"mtime"`
}

type VersionInfo struct {
	ShiftList []float64 `json:"shift_list"`
}

func (info *VersionInfo) Scan(value interface{}) error {
	bs, ok := value.([]byte)
	if !ok {
		return errors.New("value is not []byte")
	}
	return json.Unmarshal(bs, info)
}

func (info *VersionInfo) Value() (driver.Value, error) {
	return json.Marshal(info)
}

func Test_ScanValue(t *testing.T) {
	Convey("ScanValue", t, func() {
		Convey("ScanValue test1", func() {
			ctx := testutil.NewContext(testutil.NewContextRequest{})
			temp := &SocManpowerShiftForecast{
				ID:           0,
				StationID:    102,
				Operator:     "kf",
				ForecastDate: 1751299200,
				SolutionType: 1,
				Version: &VersionInfo{
					ShiftList: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0},
				},
				Ctime: 1751299200,
				Mtime: 1751299200,
			}
			db := testutil.GetDBCommon(ctx)
			err := db.Table(SocManpowerShiftForecastTabName).Create(temp).Error
			So(err, ShouldEqual, nil)
			var res []*SocManpowerShiftForecast
			err = db.Table(SocManpowerShiftForecastTabName).Find(&res).Error
			So(err, ShouldEqual, nil)
			fmt.Println(stringutil.Object2String(res))
		})
	})
}
