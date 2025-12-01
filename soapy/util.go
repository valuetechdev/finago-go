package soapy

import (
	"time"

	"github.com/hooklift/gowsdl/soap"
)

// Helper function to ensure valid time for SOAP, considering DST transitions
// for Finago Office's SOAP API.
func EnsureValidTimeForSOAP(t time.Time) soap.XSDDateTime {
	// Ensure the time is in a specific timezone (e.g., the server's timezone or your local timezone)
	loc, _ := time.LoadLocation("UTC")

	// Convert the time to the target timezone
	timeInTargetTZ := t.In(loc)

	// Check if the time is valid in the timezone (not in a DST "gap")
	// This is a simplified check - if the time doesn't change when normalized, it's valid
	normalized := time.Date(
		timeInTargetTZ.Year(),
		timeInTargetTZ.Month(),
		timeInTargetTZ.Day(),
		timeInTargetTZ.Hour(),
		timeInTargetTZ.Minute(),
		timeInTargetTZ.Second(),
		timeInTargetTZ.Nanosecond(),
		loc,
	)

	// The time is valid
	if normalized.Equal(timeInTargetTZ) {
		return soap.CreateXsdDateTime(timeInTargetTZ, true)
	}

	// If invalid (in DST gap), adjust the time forward by 1 hour to skip the gap
	adjustedTime := timeInTargetTZ.Add(time.Hour)
	return soap.CreateXsdDateTime(adjustedTime, true)
}
