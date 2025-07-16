package utils

import "time"

// String returns a pointer to the string value passed in
func String(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer or empty string if nil
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// Int returns a pointer to the int value passed in
func Int(v int) *int {
	return &v
}

// IntValue returns the value of the int pointer or 0 if nil
func IntValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

// Int64 returns a pointer to the int64 value passed in
func Int64(v int64) *int64 {
	return &v
}

// Int64Value returns the value of the int64 pointer or 0 if nil
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// Float32 returns a pointer to the float32 value passed in
func Float32(v float32) *float32 {
	return &v
}

// Float32Value returns the value of the float32 pointer or 0 if nil
func Float32Value(v *float32) float32 {
	if v != nil {
		return *v
	}
	return 0
}

// Float64 returns a pointer to the float64 value passed in
func Float64(v float64) *float64 {
	return &v
}

// Float64Value returns the value of the float64 pointer or 0 if nil
func Float64Value(v *float64) float64 {
	if v != nil {
		return *v
	}
	return 0
}

// Bool returns a pointer to the bool value passed in
func Bool(v bool) *bool {
	return &v
}

// BoolValue returns the value of the bool pointer or false if nil
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// Time returns a pointer to the time.Time value passed in
func Time(v time.Time) *time.Time {
	return &v
}

// TimeValue returns the value of the time.Time pointer or zero time if nil
func TimeValue(v *time.Time) time.Time {
	if v != nil {
		return *v
	}
	return time.Time{}
}