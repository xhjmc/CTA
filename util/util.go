package util

import (
	"errors"
	"strconv"
	"time"
)

// try once first, and if it fails, retry ${retryTimes} times
func Retry(retryTimes int, retryInterval time.Duration, f func() bool) bool {
	if f() {
		return true
	}
	for i := 0; i < retryTimes; i++ {
		b := f()
		if b {
			return true
		}
		time.Sleep(retryInterval)
	}
	return false
}

func Interface2Int(item interface{}) (int, error) {
	i64, err := Interface2Int64(item)
	return int(i64), err
}

func Interface2Int64(item interface{}) (int64, error) {
	switch val := item.(type) {
	case byte:
		return int64(val), nil
	case bool:
		if val {
			return 1, nil
		} else {
			return 0, nil
		}
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		return 0, errors.New("unsupported type")
	}
}

func Interface2Bool(item interface{}) (bool, error) {
	switch val := item.(type) {
	case byte:
		return val != 0, nil
	case bool:
		return val, nil
	case int:
		return val != 0, nil
	case int8:
		return val != 0, nil
	case int16:
		return val != 0, nil
	case int32:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case float32:
		return val != 0, nil
	case float64:
		return val != 0, nil
	case string:
		return strconv.ParseBool(val)
	default:
		return false, errors.New("unsupported type")
	}
}
