package timehelper

import (
	"testing"
	"time"
)

func TestFormatDateJsISOString(t *testing.T) {
	expected := "2017-12-05T07:12:15.234Z"
	timeToFormat := time.Date(2017, 12, 5, 7, 12, 15, 234000000, time.UTC)

	acctual := FormatDateTimeJsISOString(&timeToFormat)
	if acctual != expected {
		t.Errorf("expected: %v, acctual: %v", expected, acctual)
	}
}
