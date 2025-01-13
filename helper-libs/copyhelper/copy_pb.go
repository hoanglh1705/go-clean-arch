package copyhelper

import (
	"database/sql"
	"go-clean-arch/helper-libs/timehelper"
	"time"

	"github.com/jinzhu/copier"
)

type (
	pbConverter struct {
		timeHelper timehelper.TimeHelper
	}
)

func NewPbConverter(
	opts *CopyOptions,
) PbConverter {
	timezone := timehelper.Timezone_UTC
	if opts.Timezone != "" {
		timezone = opts.Timezone
	}
	return &pbConverter{
		timeHelper: timehelper.NewTimeHelper(&timehelper.TimeOptions{
			Timezone: timezone,
		}),
	}
}

func (h *pbConverter) FromPb(to interface{}, from interface{}) {
	_ = copier.CopyWithOption(to, from, copier.Option{
		Converters: []copier.TypeConverter{
			{
				SrcType: copier.String,
				DstType: time.Time{},
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(string)
					if !ok || s == "" {
						return time.Time{}, nil
					}

					parsedTime, err := h.timeHelper.ParseDateTimeFromStringYYYY_MM_DD(s)
					if err != nil {
						return time.Time{}, nil
					}

					return parsedTime, nil
				},
			},
			{
				SrcType: copier.String,
				DstType: sql.NullTime{},
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(string)
					if !ok || s == "" {
						return nil, nil
					}

					parsedTime, err := h.timeHelper.ParseDateTimeFromStringYYYY_MM_DD(s)
					if err != nil {
						return nil, nil
					}

					return parsedTime, nil
				},
			},
		},
	})
}

func (h *pbConverter) ToPb(to interface{}, from interface{}) {
	_ = copier.CopyWithOption(to, from, copier.Option{
		Converters: []copier.TypeConverter{
			{
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(time.Time)
					if !ok {
						return src, nil
					}

					if s.IsZero() {
						return "", nil
					}

					return h.timeHelper.FormatDateTimeJSISOString(s), nil
				},
			},
			{
				SrcType: sql.NullTime{},
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(time.Time)
					if !ok {
						return nil, nil
					}

					if s.IsZero() {
						return "", nil
					}

					return h.timeHelper.FormatDateTimeJSISOString(s), nil
				},
			},
		},
	})
}
