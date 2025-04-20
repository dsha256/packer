package safeconv_test

import (
	"testing"

	"github.com/dsha256/packer/pkg/safeconv"
)

// TestParseInt64 tests the ParseInt64 function.
func TestParseInt64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
		want int64
	}{
		{
			name: "Positive Number",
			data: "100",
			want: 100,
		},
		{
			name: "Negative Number",
			data: "-200",
			want: -200,
		},
		{
			name: "Zero",
			data: "0",
			want: 0,
		},
		{
			name: "Non Numeric",
			data: "abcd",
			want: 0,
		},
		{
			name: "Empty String",
			data: "",
			want: 0,
		},
		{
			name: "Large Number",
			data: "9223372036854775807",
			want: 9223372036854775807,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := safeconv.ParseInt64(tt.data); got != tt.want {
				t.Errorf("ParseInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseInt tests the ParseInt function.
func TestParseInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
		want int
	}{
		{
			name: "Positive Number",
			data: "100",
			want: 100,
		},
		{
			name: "Negative Number",
			data: "-200",
			want: -200,
		},
		{
			name: "Zero",
			data: "0",
			want: 0,
		},
		{
			name: "Non Numeric",
			data: "abcd",
			want: 0,
		},
		{
			name: "Empty String",
			data: "",
			want: 0,
		},
		{
			name: "Large Number",
			data: "2147483647",
			want: 2147483647,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := safeconv.ParseInt(tt.data); got != tt.want {
				t.Errorf("ParseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
