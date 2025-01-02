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

func IsValidTimeFormat(input string) error {
	_, err := time.Parse(constants.HourFormater, input)
	if err != nil {
		logger.Errorf("Invalid time format: %v", err)
	}
	return err
}

func ConvertToCronFormat(timeStr string) string {
	parts := strings.Split(timeStr, ":")

	hour := parts[0]
	minute := parts[1]

	return fmt.Sprintf("%s %s * * *", minute, hour)
}

// FindKeyInMap to recursively search for a key and return its value
func FindKeyInMap(data map[string]interface{}, key string) (interface{}, bool) {
	for k, v := range data {
		if k == key {
			return v, true
		}

		if nestedMap, ok := v.(map[string]interface{}); ok {
			if result, found := FindKeyInMap(nestedMap, key); found {
				return result, true
			}
		}
	}
	return nil, false
}
