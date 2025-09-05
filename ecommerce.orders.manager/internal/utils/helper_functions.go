package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ParseUUID parses string to UUID with validation
func ParseUUID(s string) (uuid.UUID, error) {
	s = strings.TrimSpace(s)
	return uuid.Parse(s)
}

// ParseIntWithDefault parses string to int with default value
func ParseIntWithDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	
	return val
}

// ValidatePageSize validates and normalizes page size
func ValidatePageSize(size int, max int) int {
	if size < 1 {
		return 10
	}
	if size > max {
		return max
	}
	return size
}

// ValidatePage validates and normalizes page number
func ValidatePage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

func StringPtr(s string) *string {
	return &s
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func BoolPtr(b bool) *bool {
	return &b
}

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}