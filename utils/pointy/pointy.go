package pointy

import "time"

var zeroTime time.Time

// GetPointer returns pointer from comparable type
func GetPointer[T comparable](value T) *T {
	return &value
}

// GetSlicePointer returns pointer from comparable type slice
func GetSlicePointer[T comparable](value []T) (result []*T) {
	for _, v := range value {
		result = append(result, GetPointer(v))
	}
	return result
}

// GetStatusPointer returns the status pointer in rpc
func GetStatusPointer(value *uint32) (result *uint8) {
	if value == nil {
		return nil
	}

	result = new(uint8)

	*result = uint8(*value)

	return result
}

// GetTimePointer returns the time pointer from unix time int64
func GetTimePointer(value *int64, nsec int64) (result *time.Time) {
	if value == nil {
		return nil
	}

	result = new(time.Time)

	*result = time.Unix(*value, nsec)

	return result
}

// GetTimeMilliPointer returns the time pointer from unix milli time int64
func GetTimeMilliPointer(value *int64) (result *time.Time) {
	if value == nil {
		return nil
	}

	result = new(time.Time)

	*result = time.UnixMilli(*value)

	return result
}

// GetUnixMilliPointer returns nil when int64 is -621355968000, your time should not be from before 1970
//
// Example:
//
//	 var zeroTime time.Time
//		zeroTimeP := GetUnixMilliPointer(zeroTime.UnixMilli())
//		fmt.Println(zeroTimeP)
//
// Result:  <nil>
func GetUnixMilliPointer(value int64) *int64 {
	if value == zeroTime.UnixMilli() {
		return nil
	}
	return &value
}
