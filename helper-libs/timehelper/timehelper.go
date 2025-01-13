package timehelper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-clean-arch/helper-libs/loghelper"
)

type (
	TimeHelper interface {
		NewDate() time.Time
		NewDateTime() time.Time
		ParseDateTimeJSISOString(dateTimeStr string) (time.Time, error)
		ParseDateTimeFromICoreDateTime(iCoreDate, iCoreTime string) (time.Time, error)
		ParseDateTimeFromYMD(year, month, day int) time.Time
		ParseDateTimeFromStringYYYYMMDD(date string) (time.Time, error)
		ParseDateTimeFromStringYYYY_MM_DD(date string) (time.Time, error)
		FormatDateTimeJSISOString(datetime time.Time) string
		FormatDateYYYYMMDD(datetime time.Time) string
		FormatDateYYYY_MM_DD(datetime time.Time) string
	}

	TimeOptions struct {
		Timezone Timezone
	}
)

type timeHelper struct {
	timeLocation *time.Location
	Timezone     Timezone
}

func NewTimeHelper(
	opts *TimeOptions,
) TimeHelper {
	timezone := Timezone_UTC
	if opts.Timezone != "" {
		timezone = opts.Timezone
	}

	timeLocation, err := time.LoadLocation(string(timezone))
	if err != nil {
		loghelper.Logger.Errorf("failed to create time now in %v, use UTC instead", timeLocation)
		timeLocation = time.UTC
	}

	return &timeHelper{
		timeLocation: timeLocation,
		Timezone:     timezone,
	}
}

func (h *timeHelper) NewDate() time.Time {
	timeNow := time.Now().In(h.timeLocation)
	return time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, timeNow.Location())
}

func (h *timeHelper) NewDateTime() time.Time {
	return time.Now().In(h.timeLocation)
}
func ParseDateTimeJSISOString(dateTimeStr string) (time.Time, error) {
	return time.Parse(DATETIME_FORMAT_JS_ISO_STRING, dateTimeStr)
}

func ParseDateTimeFromICoreDateTime(iCoreDate, iCoreTime string, tz Timezone) (time.Time, error) {
	tz0 := Timezone_UTC
	if tz != "" {
		tz0 = tz
	}

	tzoffset0, ok := TimezoneOffsetMap[tz0]
	if !ok {
		loghelper.Logger.Errorf("failed to load timezone %v", tzoffset0)
		return time.Time{}, errors.New("timezone_invalid")
	}

	formattedTime := strings.Join(strings.Split(iCoreTime, " "), ".")
	return time.Parse(DATETIME_FORMAT_ICORE_TRXN, fmt.Sprintf("%sT%s%s", iCoreDate, formattedTime, tzoffset0))
}

func (h *timeHelper) ParseDateTimeJSISOString(dateTimeStr string) (time.Time, error) {
	return ParseDateTimeJSISOString(dateTimeStr) // dateTimeStr already includes timezone
}

func (h *timeHelper) ParseDateTimeFromICoreDateTime(iCoreDate, iCoreTime string) (
	time.Time, error) {
	return ParseDateTimeFromICoreDateTime(iCoreDate, iCoreTime, h.Timezone)
}

func (h *timeHelper) ParseDateTimeFromYMD(year, month, day int) time.Time {
	timeNow := time.Now().In(h.timeLocation)
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, timeNow.Location())
}

func (h *timeHelper) ParseDateTimeFromStringYYYYMMDD(date string) (time.Time, error) {
	_, err := time.Parse("20060102", date)
	if err != nil {
		loghelper.Logger.Infof("failed to parse date: %s, err: %v", date, err)
		return time.Time{}, errors.New("date_format_invalid")
	}

	year, _ := strconv.Atoi(date[:4])
	month, _ := strconv.Atoi(date[4:6])
	day, _ := strconv.Atoi(date[6:])
	timeNow := time.Now().In(h.timeLocation)
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, timeNow.Location()), nil
}

func (h *timeHelper) ParseDateTimeFromStringYYYY_MM_DD(date string) (time.Time, error) {
	// re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	// if !re.MatchString(date) {
	// 	return time.Time{}, errors.New("date_format_invalid")
	// }
	// // allSubMatchs := re.FindAllString(str1, -1)

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		loghelper.Logger.Infof("failed to parse date: %s, err: %v", date, err)
		return time.Time{}, errors.New("date_format_invalid")
	}

	dateParts := strings.Split(date, "-")
	year, _ := strconv.Atoi(dateParts[0])
	month, _ := strconv.Atoi(dateParts[1])
	day, _ := strconv.Atoi(dateParts[2])
	timeNow := time.Now().In(h.timeLocation)
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, timeNow.Location()), nil
}

func (h *timeHelper) FormatDateTimeJSISOString(datetime time.Time) string {
	return FormatDateTimeISOString(datetime.In(h.timeLocation))
}

func (h *timeHelper) FormatDateYYYYMMDD(datetime time.Time) string {
	return FormatDateYYYYMMDD(datetime)
}

func (h *timeHelper) FormatDateYYYY_MM_DD(datetime time.Time) string {
	return FormatDateYYYY_MM_DD(datetime)
}
