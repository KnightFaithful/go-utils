package repository

import (
	"example.com/m/util/fileutil"
	"example.com/m/util/helper"
	"example.com/m/util/paginator"
	"example.com/m/util/utilerror"
	"fmt"

	"reflect"
)

type Module string

const (
	ModuleEvent       Module = "event"
	ModuleEventStaff  Module = "event_staff"
	ModuleStatistic   Module = "statistic"
	ModuleClockRecord Module = "clock_record"
)

const PathPrefix = "data"

func SaveArray(array interface{}, cid string, module Module, year, month, day int64) *utilerror.UtilError {
	pageSize, err := helper.CalculateRepositoryEveryPageSize(array)
	if err != nil {
		return err.Mark()
	}
	path := GetPath(cid, module, year, month, day)
	err = fileutil.DeletePath(path)
	if err != nil {
		return err.Mark()
	}
	for pageNo := int64(1); true; pageNo++ {
		cur := paginator.PageListAndReturn(array, pageNo, pageSize)
		if reflect.ValueOf(cur).Len() == 0 {
			return nil
		}
		fileName := fmt.Sprintf("%v.json", pageNo)
		if err := fileutil.SaveJson(cur, path, fileName); err != nil {
			return err.Mark()
		}
	}
	return nil
}

func GetPath(cid string, module Module, year, month, day int64) string {
	path := fmt.Sprintf("%v/%v/%v/%v/", PathPrefix, cid, module, year)
	if month > 0 {
		path += fmt.Sprintf("%v/", month)
	}
	if day > 0 {
		path += fmt.Sprintf("%v/", day)
	}
	return path
}
