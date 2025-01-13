package timehelper

import (
	"errors"
	"fmt"
	"time"

	"go-clean-arch/helper-libs/loghelper"
)

const (
	DATE_FORMAT_YYYYMMDD          string = "20060102"
	DATE_FORMAT_YYYY_MM_DD        string = "2006-01-02"
	DATETIME_FORMAT_ICORE_TRXN    string = "20060102T15:04:05.000Z07:00"
	DATETIME_FORMAT_JS_ISO_STRING string = "2006-01-02T15:04:05.000Z07:00"
	DATE_FORMAT_SIMPLE_STRING     string = "20060102"
	DATE_FORMAT_ISO_STRING        string = "2006-01-02"
	DATE_FORMAT_VIETNAMESE        string = "02/01/2006"
	TIME_FORMAT_SIMPLE_STRING     string = "150405"
	TIME_FORMAT_ISO_STRING        string = "15:04:05"
)

const (
	Location_UTC             string = "UTC"
	Location_AsiaHo_Chi_Minh string = "Asia/Ho_Chi_Minh"
)

type Timezone string

const (
	Timezone_UTC              Timezone = "UTC"
	Timezone_Local            Timezone = "Local"
	Timezone_Asia_Ho_Chi_Minh Timezone = "Asia/Ho_Chi_Minh"
)

var (
	TimezoneOffsetMap = map[Timezone]string{
		Timezone_UTC:              "Z",
		Timezone_Asia_Ho_Chi_Minh: "+07:00",
	}
)

func CheckTimezoneSupported(tz Timezone) error {
	_, err := time.LoadLocation(string(tz))
	return err
}

func NewLocalTime() time.Time {
	return time.Now()
}

func NewUTCTime() time.Time {
	return time.Now().UTC()
}

func NewTimeInLocation(name string) time.Time {
	location, err := time.LoadLocation(name)
	if err != nil {
		return time.Now().UTC()
	}
	return time.Now().In(location)
}

func NewUnixTimestamp() int64 {
	return time.Now().Unix()
}

func NewUnixTimestampMilli() int64 {
	return time.Now().UnixMilli()
}

func ParseDateTimeFromISOString(utcTimeStr string) (time.Time, error) {
	return time.Parse(DATETIME_FORMAT_JS_ISO_STRING, utcTimeStr)
}

func ParseDateTimeFromStringYYYYMMDD(date string, tz Timezone) (time.Time, error) {
	tz0 := Timezone_UTC
	if tz != "" {
		tz0 = tz
	}

	tzoffset0, ok := TimezoneOffsetMap[tz0]
	if !ok {
		loghelper.Logger.Errorf("failed to load timezone %v", tzoffset0)
		return time.Time{}, errors.New("timezone_invalid")
	}

	return time.Parse(DATETIME_FORMAT_ICORE_TRXN, fmt.Sprintf("%sT00:00:00.000%s", date, tzoffset0))
}

func ParseDateTimeFromStringYYYY_MM_DD(date string, tz Timezone) (time.Time, error) {
	tz0 := Timezone_UTC
	if tz != "" {
		tz0 = tz
	}

	tzoffset0, ok := TimezoneOffsetMap[tz0]
	if !ok {
		loghelper.Logger.Errorf("failed to load timezone %v", tzoffset0)
		return time.Time{}, errors.New("timezone_invalid")
	}

	return time.Parse(DATETIME_FORMAT_JS_ISO_STRING, fmt.Sprintf("%sT00:00:00.000%s", date, tzoffset0))
}

func ParseUnixTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func ParseUnixTimestampMilli(timestamp int64) time.Time {
	return time.UnixMilli(timestamp)
}

func FormatDateTimeJsISOString(datetime *time.Time) string {
	if datetime == nil {
		return time.Now().UTC().Format(DATETIME_FORMAT_JS_ISO_STRING)
	}
	return datetime.Format(DATETIME_FORMAT_JS_ISO_STRING)
}

func FormatDateSimpleString(utcTime *time.Time) string {
	if utcTime == nil {
		return time.Now().UTC().Format(DATE_FORMAT_SIMPLE_STRING)
	}
	return time.Now().UTC().Format(DATE_FORMAT_SIMPLE_STRING)
}

func FormatDateISOString(utcTime *time.Time) string {
	if utcTime == nil {
		return time.Now().UTC().Format(DATE_FORMAT_ISO_STRING)
	}
	return time.Now().UTC().Format(DATE_FORMAT_ISO_STRING)
}

func FormatDateVietnamese(utcTime *time.Time) string {
	if utcTime == nil {
		return time.Now().UTC().Format(DATE_FORMAT_VIETNAMESE)
	}

	return utcTime.Format(DATE_FORMAT_VIETNAMESE)
}

func FormatTimeSimpleString(utcTime *time.Time) string {
	if utcTime == nil {
		return time.Now().UTC().Format(TIME_FORMAT_SIMPLE_STRING)
	}
	return utcTime.Format(TIME_FORMAT_SIMPLE_STRING)
}

func FormatTimeISOString(utcTime *time.Time) string {
	if utcTime == nil {
		return time.Now().UTC().Format(TIME_FORMAT_ISO_STRING)
	}
	return utcTime.Format(TIME_FORMAT_ISO_STRING)
}

func FormatDateTimeISOString(datetime time.Time) string {
	return datetime.Format(DATETIME_FORMAT_JS_ISO_STRING)
}

func FormatDateYYYYMMDD(datetime time.Time) string {
	return datetime.Format(DATE_FORMAT_YYYYMMDD)
}

func FormatDateYYYY_MM_DD(datetime time.Time) string {
	return datetime.Format(DATE_FORMAT_YYYY_MM_DD)
}

func FormatDateNowYYYYMMDD() string {
	return FormatDateYYYYMMDD(time.Now())
}

func FormatDateNowYYYY_MM_DD() string {
	return FormatDateYYYY_MM_DD(time.Now())
}

func AddDate(utcTime *time.Time, y, m, d int) time.Time {
	return utcTime.AddDate(y, m, d)
}

func AddTime(utcTime *time.Time, duration time.Duration) time.Time {
	return utcTime.Add(duration)
}

func AddMonth(dateTime time.Time, num int) time.Time {
	date := dateTime.AddDate(0, num, 0)
	return date
}

func AddDay(dateTime time.Time, num int) time.Time {
	date := dateTime.AddDate(0, 0, num)
	return date
}
