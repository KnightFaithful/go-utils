package esclient

import (
	"context"
	"encoding/json"
	testutil2 "example.com/m/test/testutil"
	"example.com/m/util/convert"
	"example.com/m/util/copier"
	"example.com/m/util/printhelper"
	"example.com/m/util/utilerror"
	"fmt"
	"github.com/olivere/elastic/v7"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

var (
	ctx      context.Context
	esClient *EsClient
	index    = "wfm_attendance_records"
)

func init() {
	ctx = testutil2.NewContext(testutil2.NewContextRequest{})
	esClient = NewEsClient()
}

type AttendanceClockRecordStatisticTab struct {
	Id                    int64  `gorm:"column:id" json:"id"`
	StaffId               string `gorm:"column:staff_id" json:"staff_id"`
	BizStaffId            string `gorm:"column:biz_staff_id" json:"biz_staff_id"`
	StaffType             int64  `gorm:"column:staff_type" json:"staff_type"`
	FunctionType          int64  `gorm:"column:function_type" json:"function_type"` // 打卡选择function
	StaffName             string `gorm:"column:staff_name" json:"staff_name"`
	StaffEmail            string `gorm:"column:staff_email" json:"staff_email"`
	StationId             int64  `gorm:"column:station_id" json:"station_id"` // staff打卡当下的站点，attendance data展示为reporting station
	EventId               string `gorm:"column:event_id" json:"event_id"`
	EventStationId        int64  `gorm:"column:event_station_id" json:"event_station_id"`     // 关联EVENT的站点
	EventFunctionId       string `gorm:"column:event_function_id" json:"event_function_id"`   // 对应roster_planning_origin_event_tab event_function_id
	ProfileStationId      string `gorm:"column:profile_station_id" json:"profile_station_id"` //  Driver/OPS关联的站点, profile station 多个
	Agency                string `gorm:"column:agency" json:"agency"`
	StartTime             int64  `gorm:"column:start_time" json:"start_time"` // EVENT的开始时间
	EndTime               int64  `gorm:"column:end_time" json:"end_time"`     // EVENT的结束时间
	ClockInRecordId       int64  `gorm:"column:clock_in_record_id" json:"clock_in_record_id"`
	ClockOutRecordId      int64  `gorm:"column:clock_out_record_id" json:"clock_out_record_id"`
	ClockInTime           int64  `gorm:"column:clock_in_time" json:"clock_in_time"`   // 针对Event的打卡时间
	ClockOutTime          int64  `gorm:"column:clock_out_time" json:"clock_out_time"` // 针对Event的打卡时间
	ClockInStatus         int64  `gorm:"column:clock_in_status" json:"clock_in_status"`
	ClockOutStatus        int64  `gorm:"column:clock_out_status" json:"clock_out_status"`
	ClockOutStationId     int64  `gorm:"column:clock_out_station_id" json:"clock_out_station_id"`
	Department            string `gorm:"column:department" json:"department"`                   // SPXFM-118212 整理后，该字段冗余了
	EventDepartmentId     int64  `gorm:"column:event_department_id" json:"event_department_id"` // 对应roster_planning_origin_event_tab event_department_id
	ContractType          string `gorm:"column:contract_type"  json:"contract_type"`
	UnclassifiedRecords   string `gorm:"column:unclassified_records" json:"unclassified_records"`
	OutOfEvent            int64  `gorm:"column:out_of_event" json:"out_of_event"` // 如果是在Event内打卡，则为NO｜如果是在EVENT外打卡，则为YES
	IsAttendance          int64  `gorm:"column:is_attendance" json:"is_attendance"`
	ContractGroup         int64  `gorm:"column:contract_group" json:"contract_group"`
	ClockInChannel        int64  `gorm:"column:clock_in_channel" json:"clock_in_channel"`
	ClockOutChannel       int64  `gorm:"column:clock_out_channel" json:"clock_out_channel"`
	ClockInRangeType      int64  `gorm:"column:clock_in_range_type" json:"clock_in_range_type"`
	ClockOutRangeType     int64  `gorm:"column:clock_out_range_type" json:"clock_out_range_type"`
	ClockInRemark         string `gorm:"column:clock_in_remark" json:"clock_in_remark"`
	ClockOutRemark        string `gorm:"column:clock_out_remark" json:"clock_out_remark"`
	ShiftPriority         int64  `gorm:"column:shift_priority" json:"shift_priority"`
	CreateVersion         string `gorm:"column:create_version" json:"create_version"` //创建这个统计的版本号
	AgencyId              int64  `gorm:"column:agency_id" json:"agency_id"`
	IsAgencyStaff         int64  `gorm:"column:is_agency_staff" json:"is_agency_staff"`
	Ctime                 int64  `gorm:"column:ctime" json:"ctime"`
	Mtime                 int64  `gorm:"column:mtime" json:"mtime"`
	ProfileFunctionTypeId int64  `gorm:"profile_function_type_id" json:"profile_function_type_id"` // driver function 掩码，跟EventFunctionId不同，因为要搜索关系，改掩码方式
	ProfileDepartmentId   int64  `gorm:"profile_department_id" json:"profile_department_id"`       // ops department
	SlotCode              string `gorm:"slot_code" json:"slot_code"`
	SlotId                string `gorm:"slot_id" json:"slot_id"`
	ClockInMatchingType   int64  `gorm:"column:clock_in_matching_type" json:"clock_in_matching_type"`
	ClockOutMatchingType  int64  `gorm:"column:clock_out_matching_type" json:"clock_out_matching_type"`
	FulfillWorkingHours   int64  `gorm:"column:fulfill_working_hours" json:"fulfill_working_hours"`
}

func Test_EsAdd(t *testing.T) {
	Convey("EsAdd", t, func() {
		Convey("EsAdd test1", func() {
			var list []*AttendanceClockRecordStatisticTab
			db := testutil2.GetDBCommon(ctx)
			limit := 10
			err := db.Table("attendance_clock_record_statistic_tab").Order("id desc").Limit(limit).Find(&list).Error
			So(err, ShouldEqual, nil)
			So(len(list), ShouldEqual, limit)
			esItem := &AttendanceClockRecordStatisticTab{}
			for _, temp := range list {
				err = copier.Copy(&temp, esItem)
				So(err, ShouldEqual, nil)
				err = esClient.EsAdd(ctx, index, convert.Int64ToString(temp.Id), esItem)
				So(err, ShouldEqual, nil)
			}
		})
	})
}

func Test_EsBatchGet(t *testing.T) {
	Convey("EsBatchGet", t, func() {
		Convey("EsBatchGet test1", func() {
			var resp []*AttendanceClockRecordStatisticTab
			err := esClient.EsBatchGet(ctx, index, "_id", []string{"518474", "518469", "518471"}, &resp)
			So(err, ShouldEqual, nil)
			fmt.Println(resp)
		})
	})
}

func Test_compareAndAdd(t *testing.T) {
	Convey("compareAndAdd", t, func() {
		Convey("compareAndAdd test1", func() {
			err := compareAndAdd()
			So(err, ShouldEqual, nil)
		})
	})
}

func compareAndAdd() *utilerror.UtilError {
	idGt := int64(0)
	limit := 1000
	for {
		printhelper.Println("offset:", idGt)
		var list []*AttendanceClockRecordStatisticTab
		db := testutil2.GetDBCommon(ctx)
		err := db.Table("attendance_clock_record_statistic_tab").Where("id > ?", idGt).Order("id asc").Limit(limit).Find(&list).Error
		if err != nil {
			return utilerror.NewError(err.Error())
		}
		if len(list) == 0 {
			return nil
		}
		idGt = list[len(list)-1].Id
		cErr := esClient.EsBatchCreate(ctx, list, func(data interface{}) []*elastic.BulkIndexRequest {
			itemList := data.([]*AttendanceClockRecordStatisticTab)
			var docs []*elastic.BulkIndexRequest
			for _, entity := range itemList {
				esJson, _ := json.Marshal(entity)
				doc := elastic.NewBulkIndexRequest().Index(index).Id(convert.Int64ToString(entity.Id)).Doc(string(esJson))
				docs = append(docs, doc)
			}
			return docs
		})
		if cErr != nil {
			return utilerror.NewError(cErr.Error())
		}
	}
}
