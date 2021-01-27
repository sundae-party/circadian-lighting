package main

import (
	"math"
	"testing"
	"time"
)

func TestIsLeapYear(t *testing.T) {

	dates := make(map[time.Time]bool)
	dates[time.Date(1600, 1, 1, 0, 0, 0, 0, time.UTC)] = true
	dates[time.Date(1700, 1, 1, 0, 0, 0, 0, time.UTC)] = false
	dates[time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC)] = false
	dates[time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)] = false
	dates[time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)] = true
	dates[time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)] = false
	dates[time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC)] = false
	dates[time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC)] = false
	dates[time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)] = true

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
	dates[time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)] = 0
	dates[time.Date(2021, 4, 2, 6, 0, 0, 0, time.UTC)] = math.Pi / 2
	dates[time.Date(2021, 7, 2, 12, 0, 0, 0, time.UTC)] = math.Pi
	dates[time.Date(2021, 10, 1, 18, 0, 0, 0, time.UTC)] = 3 * math.Pi / 2

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
	dates[time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)] = 60 * -3.68
	dates[time.Date(2021, 2, 1, 12, 0, 0, 0, time.UTC)] = 60 * -13.62
	dates[time.Date(2021, 3, 1, 12, 0, 0, 0, time.UTC)] = 60 * -12.28
	dates[time.Date(2021, 4, 1, 12, 0, 0, 0, time.UTC)] = 60 * -3.77
	dates[time.Date(2021, 5, 1, 12, 0, 0, 0, time.UTC)] = 60 * 2.93
	dates[time.Date(2021, 6, 1, 12, 0, 0, 0, time.UTC)] = 60 * 2.11
	dates[time.Date(2021, 7, 1, 12, 0, 0, 0, time.UTC)] = 60 * -3.93
	dates[time.Date(2021, 8, 1, 12, 0, 0, 0, time.UTC)] = 60 * -6.34
	dates[time.Date(2021, 9, 1, 12, 0, 0, 0, time.UTC)] = 60 * 0.05
	dates[time.Date(2021, 10, 1, 12, 0, 0, 0, time.UTC)] = 60 * 10.4
	dates[time.Date(2021, 11, 1, 12, 0, 0, 0, time.UTC)] = 60 * 16.47
	dates[time.Date(2021, 12, 1, 12, 0, 0, 0, time.UTC)] = 60 * 10.88

	for k, v := range dates {
		got := eqTime(k) * 60
		if math.Abs(got-v) > 15 {
			t.Errorf("eqTime(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		} else {
			t.Logf("eqTime(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		}
	}

}

func TestDecl(t *testing.T) {

	dates := make(map[time.Time]float64)
	dates[time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)] = -22.96
	dates[time.Date(2021, 2, 1, 12, 0, 0, 0, time.UTC)] = -16.95
	dates[time.Date(2021, 3, 1, 12, 0, 0, 0, time.UTC)] = -7.4
	dates[time.Date(2021, 4, 1, 12, 0, 0, 0, time.UTC)] = 4.73
	dates[time.Date(2021, 5, 1, 12, 0, 0, 0, time.UTC)] = 15.23
	dates[time.Date(2021, 6, 1, 12, 0, 0, 0, time.UTC)] = 22.12
	dates[time.Date(2021, 7, 1, 12, 0, 0, 0, time.UTC)] = 23.06
	dates[time.Date(2021, 8, 1, 12, 0, 0, 0, time.UTC)] = 17.88
	dates[time.Date(2021, 9, 1, 12, 0, 0, 0, time.UTC)] = 8.09
	dates[time.Date(2021, 10, 1, 12, 0, 0, 0, time.UTC)] = -3.38
	dates[time.Date(2021, 11, 1, 12, 0, 0, 0, time.UTC)] = -14.59
	dates[time.Date(2021, 12, 1, 12, 0, 0, 0, time.UTC)] = -21.88

	for k, v := range dates {
		got := toDegrees(decl(k))
		if math.Abs(got-v) > 0.5 {
			t.Errorf("decl(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		} else {
			t.Logf("decl(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		}
	}

}

func TestElevation(t *testing.T) {

	// Paris UTC
	latitude := 48.87
	longitude := 2.67
	dates := make(map[time.Time]float64)
	dates[time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)] = 18.2
	dates[time.Date(2021, 2, 1, 12, 0, 0, 0, time.UTC)] = 24.21
	dates[time.Date(2021, 3, 1, 12, 0, 0, 0, time.UTC)] = 33.76
	dates[time.Date(2021, 4, 1, 12, 0, 0, 0, time.UTC)] = 45.86
	dates[time.Date(2021, 5, 1, 12, 0, 0, 0, time.UTC)] = 56.25
	dates[time.Date(2021, 6, 1, 12, 0, 0, 0, time.UTC)] = 63.14
	dates[time.Date(2021, 7, 1, 12, 0, 0, 0, time.UTC)] = 64.17
	dates[time.Date(2021, 8, 1, 12, 0, 0, 0, time.UTC)] = 59.01
	dates[time.Date(2021, 9, 1, 12, 0, 0, 0, time.UTC)] = 49.17
	dates[time.Date(2021, 10, 1, 12, 0, 0, 0, time.UTC)] = 37.57
	dates[time.Date(2021, 11, 1, 12, 0, 0, 0, time.UTC)] = 26.29
	dates[time.Date(2021, 12, 1, 12, 0, 0, 0, time.UTC)] = 19.14

	for k, v := range dates {
		got := toDegrees(elevation(k, latitude, longitude))
		if math.Abs(got-v) > 0.5 {
			t.Errorf("elevation(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		} else {
			t.Logf("elevation(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		}
	}

}

func TestAzimuth(t *testing.T) {

	// Paris UTC
	latitude := 48.87
	longitude := 2.67
	dates := make(map[time.Time]float64)
	dates[time.Date(2021, 1, 1, 6, 0, 0, 0, time.UTC)] = 106.83
	dates[time.Date(2021, 1, 1, 7, 0, 0, 0, time.UTC)] = 117.4
	dates[time.Date(2021, 1, 1, 8, 0, 0, 0, time.UTC)] = 128.46
	dates[time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC)] = 140.38
	dates[time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC)] = 153.34
	dates[time.Date(2021, 1, 1, 11, 0, 0, 0, time.UTC)] = 167.24
	dates[time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)] = 181.7
	dates[time.Date(2021, 1, 1, 13, 0, 0, 0, time.UTC)] = 196.08
	dates[time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC)] = 209.79
	dates[time.Date(2021, 1, 1, 15, 0, 0, 0, time.UTC)] = 222.5
	dates[time.Date(2021, 1, 1, 16, 0, 0, 0, time.UTC)] = 234.2
	dates[time.Date(2021, 1, 1, 17, 0, 0, 0, time.UTC)] = 245.12

	for k, v := range dates {
		got := toDegrees(azimuth(k, latitude, longitude))
		if math.Abs(got-v) > 0.5 {
			t.Errorf("azimuth(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		} else {
			t.Logf("azimuth(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		}
	}

	// Tongatapu UTC
	latitude = -21.133
	longitude = -175.217
	dates[time.Date(2021, 1, 1, 6, 0, 0, 0, time.UTC)] = 115.37
	dates[time.Date(2021, 1, 1, 7, 0, 0, 0, time.UTC)] = 110.42
	dates[time.Date(2021, 1, 1, 8, 0, 0, 0, time.UTC)] = 128.46
	dates[time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC)] = 140.38
	dates[time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC)] = 153.34
	dates[time.Date(2021, 1, 1, 11, 0, 0, 0, time.UTC)] = 167.24
	dates[time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)] = 181.7
	dates[time.Date(2021, 1, 1, 13, 0, 0, 0, time.UTC)] = 196.08
	dates[time.Date(2021, 1, 1, 14, 0, 0, 0, time.UTC)] = 209.79
	dates[time.Date(2021, 1, 1, 15, 0, 0, 0, time.UTC)] = 222.5
	dates[time.Date(2021, 1, 1, 16, 0, 0, 0, time.UTC)] = 234.2
	dates[time.Date(2021, 1, 1, 17, 0, 0, 0, time.UTC)] = 245.12

	for k, v := range dates {
		got := toDegrees(azimuth(k, latitude, longitude))
		if math.Abs(got-v) > 0.5 {
			t.Errorf("azimuth(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		} else {
			t.Logf("azimuth(%v) = %f, expected %f, diff %f", k, got, v, math.Abs(got-v))
		}
	}

}
