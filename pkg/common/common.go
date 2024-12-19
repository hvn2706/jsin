package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"jsin/logger"
	"jsin/pkg/constants"
)

func LogStruct(v any) string {
	marshal, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func LoadTimeZone() *time.Location {
	location, err := time.LoadLocation(constants.TimeZone)
	if err != nil {
		logger.Errorf("Failed to load location: %v", err)
		return time.UTC
	}
	return location
}

func IsValidTimeFormat(input string) bool {
	_, err := time.Parse(constants.HourFormater, input)
	return err == nil
}

func ConvertToCronFormat(timeStr string) string {
	parts := strings.Split(timeStr, ":")

	hour := parts[0]
	minute := parts[1]

	return fmt.Sprintf("%s %s * * *", minute, hour)
}
