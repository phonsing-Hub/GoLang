package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/phonsing-Hub/GoLang/utils/response"

	"time"
)

var VLD = validator.New()

func StringToDate(dateString string, layout string) (time.Time, error) {
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, response.Fail(nil, "INVALID_DATE_FORMAT", "Invalid date format", 400)
	}
	return t, nil
}
