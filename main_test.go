package main

import (
	"math"
	"testing"
	"time"
)

func TestIsLeapYear(t *testing.T) {

	dates := make(map[time.Time]bool)
	dates[time.Date(1600, 1, 1, 0, 0, 0, 0, &time.Location{})] = true
	dates[time.Date(1700, 1, 1, 0, 0, 0, 0, &time.Location{})] = false
	dates[time.Date(1800, 1, 1, 0, 0, 0, 0, &time.Location{})] = false
	dates[time.Date(1900, 1, 1, 0, 0, 0, 0, &time.Location{})] = false
	dates[time.Date(2000, 1, 1, 0, 0, 0, 0, &time.Location{})] = true
	dates[time.Date(2100, 1, 1, 0, 0, 0, 0, &time.Location{})] = false
	dates[time.Date(2200, 1, 1, 0, 0, 0, 0, &time.Location{})] = false
	dates[time.Date(2300, 1, 1, 0, 0, 0, 0, &time.Location{})] = false
	dates[time.Date(2016, 1, 1, 0, 0, 0, 0, &time.Location{})] = true

	for k, v := range dates {
		got := isLeapYear(k)
		if got != v {
			t.Errorf("isLeapYear(%v) = %t, expected %t", k, got, v)
		} else {
			t.Logf("isLeapYear(%v) = %t, expected %t", k, got, v)
		}
	}

}

func TestFractionalYear(t *testing.T) {

	dates := make(map[time.Time]float64)
	dates[time.Date(2021, 1, 1, 0, 0, 0, 0, &time.Location{})] = 0
	dates[time.Date(2021, 12, 31, 23, 59, 59, 0, &time.Location{})] = 2 * math.Pi
	dates[time.Date(2021, 7, 2, 12, 0, 0, 0, &time.Location{})] = math.Pi

	for k, v := range dates {
		got := fractionalYear(k)
		if math.Abs(got-v) > 0.000001 {
			t.Errorf("fractionalYear(%v) = %f, expected %f", k, got, v)
		} else {
			t.Logf("fractionalYear(%v) = %f, expected %f", k, got, v)
		}
	}

}
