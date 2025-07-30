package metrics

import (
	"sort"
	"time"
)

// Computes the middle value of a sorted integer array, handling even-length arrays
func calculateMedianInt(values []int) float64 {
	if len(values) == 0 {
		return 0
	}

	// Convert to float64 for consistent return type
	floatValues := make([]float64, len(values))
	for i, v := range values {
		floatValues[i] = float64(v)
	}

	return calculateMedianFloat(floatValues)
}

// Computes the middle value of a sorted float array, handling even-length arrays
func calculateMedianFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Sort the values
	sort.Float64s(values)

	// Find the median
	length := len(values)
	if length%2 == 0 {
		// Even number of values, average the middle two
		return (values[length/2-1] + values[length/2]) / 2
	}

	// Odd number of values, return the middle one
	return values[length/2]
}

// Determines the Monday of the ISO week containing the given date
func getStartOfISOWeek(date time.Time) time.Time {
	// Get the weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	weekday := int(date.Weekday())

	// Convert to ISO weekday (1 = Monday, ..., 7 = Sunday)
	isoWeekday := weekday
	if isoWeekday == 0 {
		isoWeekday = 7
	}

	// Calculate days to subtract to get to Monday
	daysToSubtract := isoWeekday - 1

	// Get the start of the day
	year, month, day := date.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, date.Location())

	// Subtract days to get to Monday
	return startOfDay.AddDate(0, 0, -daysToSubtract)
}
