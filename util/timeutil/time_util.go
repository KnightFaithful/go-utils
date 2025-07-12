package timeutil

import (
	"context"
	"example.com/m/util/utilerror"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	OneHourSecond = 3600
	OneDaySecond  = 86400
)

const (
	DateFormat               = "2006-01-02"
	DateFormat2              = "2006-01-02-15"
	DateFormatYYYYMMDDHHMMSS = "2006-01-02 15:04:05"
)

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func GetTimeDayWeek(ctx context.Context, timeStamp int64) string {
	tmpTime := TimeStampTotTime(ctx, timeStamp)
	dateStr := fmt.Sprintf("%s(%s)", TimeStampTotStringWithLayout(ctx, tmpTime.Unix(), DateFormat), tmpTime.Weekday().String())
	return dateStr
}

func TimeStampTotStringWithLayout(ctx context.Context, timestamp int64, layout string) string {
	return TimeStampTotTime(ctx, timestamp).Format(layout)
}

func TimeStampTotTime(ctx context.Context, timestamp int64) time.Time {
	return time.Unix(timestamp, 0).In(getTimeZone(ctx))
}

func getTimeZone(ctx context.Context) *time.Location {
	//return time.FixedZone(cid, TimeZoneMap[cid])
	//cid := testutil2.GetValueByCtxString(ctx, testutil.ContextKeyCid)
	timeZone, _ := time.LoadLocation(RegionID)
	return timeZone
}

const (
	RegionSG = "sg"
	RegionVN = "vn"
	RegionID = "id"
	RegionTH = "th"
	RegionTW = "tw"
	RegionMY = "my"
	RegionPH = "ph"
	RegionBR = "br"
	RegionMX = "mx"
)

var TimeZoneIDMap = map[string]string{
	RegionSG:                  "Asia/Singapore",
	strings.ToUpper(RegionSG): "Asia/Singapore",
	RegionVN:                  "Asia/Ho_Chi_Minh",
	strings.ToUpper(RegionVN): "Asia/Ho_Chi_Minh",
	RegionID:                  "Asia/Jakarta",
	strings.ToUpper(RegionID): "Asia/Jakarta",
	RegionTH:                  "Asia/Bangkok",
	strings.ToUpper(RegionTH): "Asia/Bangkok",
	RegionTW:                  "Asia/Taipei",
	strings.ToUpper(RegionTW): "Asia/Taipei",
	RegionMY:                  "Asia/Kuala_Lumpur",
	strings.ToUpper(RegionMY): "Asia/Kuala_Lumpur",
	RegionPH:                  "Asia/Manila",
	strings.ToUpper(RegionPH): "Asia/Manila",
	RegionBR:                  "America/Sao_Paulo",
	strings.ToUpper(RegionBR): "America/Sao_Paulo",
	RegionMX:                  "America/Mexico_City",
	strings.ToUpper(RegionMX): "America/Mexico_City",
}

func SlotTimeToString(timeInt int64) string {
	var ret string = ""
	//计算小时
	hour := timeInt / 3600
	minuteInt := timeInt % 3600

	//计算分钟
	minute := minuteInt / 60
	second := minuteInt % 60

	//组装
	addOneFlag := false
	hourStr := strconv.FormatInt(hour, 10)
	if hour >= 24 {
		hourStr = strconv.FormatInt(hour-24, 10)
		addOneFlag = true
	}
	if hour < 10 {
		ret = ret + "0" + hourStr
	} else {
		ret = ret + hourStr
	}
	ret += ":"
	minuteStr := strconv.FormatInt(minute, 10)
	if minute < 10 {
		ret = ret + "0" + minuteStr
	} else {
		ret = ret + minuteStr
	}
	ret += ":"
	secondStr := strconv.FormatInt(second, 10)
	if second < 10 {
		ret = ret + "0" + secondStr
	} else {
		ret = ret + secondStr
	}
	if addOneFlag {
		ret = ret + "(+1)"
	}
	return ret
}

func TimeStampToString(ctx context.Context, timestamp int64, layout string) string {
	return time.Unix(timestamp, 0).In(getTimeZone(ctx)).Format(layout)
}

func GetTimeRange(ctx context.Context, year, month, day int64) (int64, int64, *utilerror.UtilError) {
	if year == 0 {
		return 0, 0, utilerror.NewError("year can not be 0")
	}
	if month == 0 {
		left := time.Date(int(year), 1, 1, 0, 0, 0, 0, getTimeZone(ctx)).Unix()
		right := time.Date(int(year+1), 1, 1, 0, 0, 0, 0, getTimeZone(ctx)).Unix() - 1
		return left, right, nil
	}
	if day == 0 {
		left := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, getTimeZone(ctx)).Unix()
		right := time.Date(int(year), time.Month(month+1), 1, 0, 0, 0, 0, getTimeZone(ctx)).Unix() - 1
		return left, right, nil
	}
	if !IsDateValid(year, month, day) {
		return 0, 0, utilerror.NewError("Invalid date")
	}
	left := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, getTimeZone(ctx)).Unix()
	right := time.Date(int(year), time.Month(month), int(day+1), 0, 0, 0, 0, getTimeZone(ctx)).Unix() - 1
	return left, right, nil
}

func IsDateValid(year, month, day int64) bool {
	// 创建日期对象
	date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)

	// 检查日期是否有效
	return date.Year() == int(year) && int64(date.Month()) == month && date.Day() == int(day)
}

func DateFormatStrToTimeStamp(ctx context.Context, dateTimeStr string, format string) (int64, *utilerror.UtilError) {
	dt, err := DateFormatStrToTime(ctx, dateTimeStr, format)
	if err != nil {
		return 0, err.Mark()
	}
	return dt.Unix(), nil
}

func DateFormatStrToTime(ctx context.Context, datetimeStr string, format string) (time.Time, *utilerror.UtilError) {
	res, err := time.ParseInLocation(format, datetimeStr, getTimeZone(ctx))
	if err != nil {
		return res, utilerror.NewError(err.Error())
	}
	return res, nil
}

// GetLocalDayZeroTime 获取某天0点
func GetLocalDayZeroTime(ctx context.Context, d time.Time) int64 {
	//当地时区
	localTimeZone := getTimeZone(ctx)
	//转换为当地时间
	localDate := d.In(localTimeZone)
	//当地0点
	localZero := time.Date(localDate.Year(), localDate.Month(), localDate.Day(), 0, 0, 0, 0, localTimeZone)
	//转换为机器所在时区0点
	currentZoneZero := localZero.In(d.Location())
	return currentZoneZero.Unix()
}

func GetLocalDayZeroTimestamp(ctx context.Context, timestamp int64) int64 {
	return GetLocalDayZeroTime(ctx, time.Unix(timestamp, 0))
}
