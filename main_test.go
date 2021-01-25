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
	dates[time.Date(2021, 4, 2, 6, 0, 0, 0, &time.Location{})] = math.Pi / 2
	dates[time.Date(2021, 7, 2, 12, 0, 0, 0, &time.Location{})] = math.Pi
	dates[time.Date(2021, 10, 1, 18, 0, 0, 0, &time.Location{})] = 3 * math.Pi / 2

	for k, v := range dates {
		got := fractionalYear(k)
		if math.Abs(got-v) > float64(1)/float64(31536000) {
			t.Errorf("fractionalYear(%v) = %e, expected %e, diff %e", k, got, v, math.Abs(got-v))
		} else {
			t.Logf("fractionalYear(%v) = %e, expected %e, diff %e", k, got, v, math.Abs(got-v))
		}
	}

}

func TestEqTime(t *testing.T) {

	dates := make(map[time.Time]float64)
	dates[time.Date(2021, 1, 1, 12, 0, 0, 0, &time.Location{})] = -229
	dates[time.Date(2021, 2, 1, 12, 0, 0, 0, &time.Location{})] = -820
	dates[time.Date(2021, 3, 1, 12, 0, 0, 0, &time.Location{})] = -733
	dates[time.Date(2021, 4, 1, 12, 0, 0, 0, &time.Location{})] = -221
	dates[time.Date(2021, 5, 1, 12, 0, 0, 0, &time.Location{})] = 177
	dates[time.Date(2021, 6, 1, 12, 0, 0, 0, &time.Location{})] = 124
	dates[time.Date(2021, 7, 1, 12, 0, 0, 0, &time.Location{})] = -239
	dates[time.Date(2021, 8, 1, 12, 0, 0, 0, &time.Location{})] = -379
	dates[time.Date(2021, 9, 1, 12, 0, 0, 0, &time.Location{})] = 8
	dates[time.Date(2021, 10, 1, 12, 0, 0, 0, &time.Location{})] = 629
	dates[time.Date(2021, 11, 1, 12, 0, 0, 0, &time.Location{})] = 988
	dates[time.Date(2021, 12, 1, 12, 0, 0, 0, &time.Location{})] = 646

	for k, v := range dates {
		got := eqTime(k) * 60
		if math.Abs(got-v) > 8 {
			t.Errorf("eqTime(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		} else {
			t.Logf("eqTime(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		}
	}

}
