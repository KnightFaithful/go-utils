package asql

import (
	"example.com/m/util/stringutil"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type OpsTabEntity struct {
	Id    int64  `gorm:"column:id;primary_key" json:"id"`
	Photo string `gorm:"column:photo" json:"photo"`
}

func Test_GenerateBatchUpdateSQL(t *testing.T) {
	Convey("GenerateBatchUpdateSQL", t, func() {
		Convey("GenerateBatchUpdateSQL test1", func() {
			sqls, utilError := GenerateBatchUpdateSQL("ops_tab", []*OpsTabEntity{
				{
					Id:    1,
					Photo: "123",
				},
				{
					Id:    2,
					Photo: "456",
				},
				{
					Id:    3,
					Photo: "789",
				},
			})
			So(utilError, ShouldEqual, nil)
			stringutil.Println(sqls...)
			//UPDATE ops_tab SET photo = CASE id WHEN 1 THEN '123' WHEN 2 THEN '456' WHEN 3 THEN '789' ELSE photo END WHERE id IN (1,2,3);
		})
	})
}

type OpsTabEntity2 struct {
	Id    int64  `gorm:"column:id;primary_key" json:"id"`
	OpsId string `gorm:"column:ops_id" json:"ops_id"`
	Photo string `gorm:"column:photo" json:"photo"`
}

func Test_GenerateInsertSQL(t *testing.T) {
	Convey("GenerateInsertSQL", t, func() {
		Convey("GenerateInsertSQL test1", func() {
			sqls, utilError := GenerateInsertSQL("ops_tab", []*OpsTabEntity2{
				{
					Id:    1,
					OpsId: "ops1",
					Photo: "123",
				},
				{
					Id:    2,
					OpsId: "ops2",
					Photo: "456",
				},
				{
					Id:    3,
					OpsId: "ops3",
					Photo: "789",
				},
			})
			So(utilError, ShouldEqual, nil)
			stringutil.Println(sqls)
			//INSERT INTO ops_tab (ops_id,photo) VALUES ('ops1','123'),('ops2','456'),('ops3','789');
		})
	})
}
