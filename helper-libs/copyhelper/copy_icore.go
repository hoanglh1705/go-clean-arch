package copyhelper

import (
	"database/sql"
	"go-clean-arch/helper-libs/timehelper"
	"time"

	"github.com/jinzhu/copier"
)

type (
	iCoreEntityConverter struct {
		timeHelper timehelper.TimeHelper
	}
)

func NewICoreEntityConverter() ICoreEntityConverter {
	return &iCoreEntityConverter{
		timeHelper: timehelper.NewTimeHelper(&timehelper.TimeOptions{
			Timezone: timehelper.Timezone_Asia_Ho_Chi_Minh,
		}),
	}
}

func (h *iCoreEntityConverter) FromModel(to interface{}, from interface{}) {
	_ = copier.CopyWithOption(to, from, copier.Option{
		Converters: []copier.TypeConverter{
			{
				SrcType: copier.String,
				DstType: time.Time{},
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(string)
					if !ok {
						return time.Time{}, nil
					}

					parsedTime, err := h.timeHelper.ParseDateTimeFromStringYYYYMMDD(s)
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
					if !ok {
						return nil, nil
					}

					parsedTime, err := h.timeHelper.ParseDateTimeFromStringYYYYMMDD(s)
					if err != nil {
						return nil, nil
					}

					return parsedTime, nil
				},
			},
		},
	})
}

func (h *iCoreEntityConverter) ToModel(to interface{}, from interface{}) {
	_ = copier.CopyWithOption(to, from, copier.Option{
		Converters: []copier.TypeConverter{
			{
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src interface{}) (interface{}, error) {
					s, ok := src.(time.Time)
					if !ok {
						return nil, nil
					}

					return h.timeHelper.FormatDateYYYYMMDD(s), nil
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

					return h.timeHelper.FormatDateYYYYMMDD(s), nil
				},
			},
		},
	})
}
