package utils_test

import (
	"testing"
	"time"

	"github.com/rizome-dev/go-moonshot/pkg/utils"
)

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"empty string", ""},
		{"non-empty string", "hello"},
		{"string with spaces", "hello world"},
		{"unicode string", "こんにちは"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := utils.String(tt.value)
			if ptr == nil {
				t.Fatal("String() returned nil")
			}
			if *ptr != tt.value {
				t.Errorf("String() = %v, want %v", *ptr, tt.value)
			}
		})
	}
}

func TestStringValue(t *testing.T) {
	tests := []struct {
		name  string
		ptr   *string
		want  string
	}{
		{"nil pointer", nil, ""},
		{"empty string pointer", utils.String(""), ""},
		{"non-empty string pointer", utils.String("hello"), "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.StringValue(tt.ptr)
			if got != tt.want {
				t.Errorf("StringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		name  string
		value int
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"max int", int(^uint(0) >> 1)},
		{"min int", -int(^uint(0)>>1) - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := utils.Int(tt.value)
			if ptr == nil {
				t.Fatal("Int() returned nil")
			}
			if *ptr != tt.value {
				t.Errorf("Int() = %v, want %v", *ptr, tt.value)
			}
		})
	}
}

func TestIntValue(t *testing.T) {
	tests := []struct {
		name  string
		ptr   *int
		want  int
	}{
		{"nil pointer", nil, 0},
		{"zero pointer", utils.Int(0), 0},
		{"positive pointer", utils.Int(42), 42},
		{"negative pointer", utils.Int(-42), -42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.IntValue(tt.ptr)
			if got != tt.want {
				t.Errorf("IntValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		name  string
		value int64
	}{
		{"zero", 0},
		{"positive", 42},
		{"negative", -42},
		{"large positive", 9223372036854775807},
		{"large negative", -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := utils.Int64(tt.value)
			if ptr == nil {
				t.Fatal("Int64() returned nil")
			}
			if *ptr != tt.value {
				t.Errorf("Int64() = %v, want %v", *ptr, tt.value)
			}
		})
	}
}

func TestInt64Value(t *testing.T) {
	tests := []struct {
		name  string
		ptr   *int64
		want  int64
	}{
		{"nil pointer", nil, 0},
		{"zero pointer", utils.Int64(0), 0},
		{"positive pointer", utils.Int64(42), 42},
		{"negative pointer", utils.Int64(-42), -42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.Int64Value(tt.ptr)
			if got != tt.want {
				t.Errorf("Int64Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32(t *testing.T) {
	tests := []struct {
		name  string
		value float32
	}{
		{"zero", 0},
		{"positive", 3.14},
		{"negative", -3.14},
		{"small", 0.000001},
		{"large", 1000000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := utils.Float32(tt.value)
			if ptr == nil {
				t.Fatal("Float32() returned nil")
			}
			if *ptr != tt.value {
				t.Errorf("Float32() = %v, want %v", *ptr, tt.value)
			}
		})
	}
}

func TestFloat32Value(t *testing.T) {
	tests := []struct {
		name  string
		ptr   *float32
		want  float32
	}{
		{"nil pointer", nil, 0},
		{"zero pointer", utils.Float32(0), 0},
		{"positive pointer", utils.Float32(3.14), 3.14},
		{"negative pointer", utils.Float32(-3.14), -3.14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.Float32Value(tt.ptr)
			if got != tt.want {
				t.Errorf("Float32Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"zero", 0},
		{"positive", 3.14159265359},
		{"negative", -3.14159265359},
		{"small", 0.000000000001},
		{"large", 1000000000000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := utils.Float64(tt.value)
			if ptr == nil {
				t.Fatal("Float64() returned nil")
			}
			if *ptr != tt.value {
				t.Errorf("Float64() = %v, want %v", *ptr, tt.value)
			}
		})
	}
}

func TestFloat64Value(t *testing.T) {
	tests := []struct {
		name  string
		ptr   *float64
		want  float64
	}{
		{"nil pointer", nil, 0},
		{"zero pointer", utils.Float64(0), 0},
		{"positive pointer", utils.Float64(3.14), 3.14},
		{"negative pointer", utils.Float64(-3.14), -3.14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.Float64Value(tt.ptr)
			if got != tt.want {
				t.Errorf("Float64Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name  string
		value bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := utils.Bool(tt.value)
			if ptr == nil {
				t.Fatal("Bool() returned nil")
			}
			if *ptr != tt.value {
				t.Errorf("Bool() = %v, want %v", *ptr, tt.value)
			}
		})
	}
}

func TestBoolValue(t *testing.T) {
	tests := []struct {
		name  string
		ptr   *bool
		want  bool
	}{
		{"nil pointer", nil, false},
		{"true pointer", utils.Bool(true), true},
		{"false pointer", utils.Bool(false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.BoolValue(tt.ptr)
			if got != tt.want {
				t.Errorf("BoolValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		value time.Time
	}{
		{"current time", now},
		{"zero time", time.Time{}},
		{"specific time", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := utils.Time(tt.value)
			if ptr == nil {
				t.Fatal("Time() returned nil")
			}
			if !ptr.Equal(tt.value) {
				t.Errorf("Time() = %v, want %v", *ptr, tt.value)
			}
		})
	}
}

func TestTimeValue(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		ptr   *time.Time
		want  time.Time
	}{
		{"nil pointer", nil, time.Time{}},
		{"current time pointer", utils.Time(now), now},
		{"zero time pointer", utils.Time(time.Time{}), time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.TimeValue(tt.ptr)
			if !got.Equal(tt.want) {
				t.Errorf("TimeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}