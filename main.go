package main

import (
	"fmt"
	"math"
	"time"
)

func isLeapYear(date time.Time) bool {
	return date.Year()%4 == 0 && (date.Year()%100 != 0 || date.Year()%400 == 0)
}

func fractionalYear(date time.Time) float64 {
	dateUTC := date.UTC()
	if isLeapYear(date) {
		return (2 * math.Pi * (float64(dateUTC.YearDay()) - 1 + float64(dateUTC.Hour())/24 + float64(dateUTC.Minute())/(24*60) + float64(dateUTC.Second())/(24*60*60)) / 366)
	} else {
		return (2 * math.Pi * (float64(dateUTC.YearDay()) - 1 + float64(dateUTC.Hour())/24 + float64(dateUTC.Minute())/(24*60) + float64(dateUTC.Second())/(24*60*60)) / 365)
	}
}

func eqTime(date time.Time) float64 {
	return 0.0116 - 7.3453*math.Sin(fractionalYear(date)+6.229) - 9.9212*math.Sin(2*fractionalYear(date)+0.3877) - 0.3363*math.Sin(3*fractionalYear(date)+0.342) - 0.2316*math.Sin(4*fractionalYear(date)+0.7531)
}

func decl(date time.Time) float64 {
	return 0.006918 - 0.399912*math.Cos(fractionalYear(date)) + 0.070257*math.Sin(fractionalYear(date)) - 0.006758*math.Cos(2*fractionalYear(date)) + 0.000907*math.Sin(2*fractionalYear(date)) - 0.002697*math.Cos(3*fractionalYear(date)) + 0.00148*math.Sin(3*fractionalYear(date))
}

func timeOffset(date time.Time, longitude float64) float64 {
	var _, offset = date.Zone()
	return eqTime(date) + 4*longitude - math.Round(float64(offset)/60)
}

func tST(date time.Time, longitude float64) float64 {
	return float64(date.Hour())*60 + float64(date.Minute()) + float64(date.Second())/60 + timeOffset(date, longitude)
}

func hA(date time.Time, longitude float64) float64 {
	return (tST(date, longitude) / 4) - 180
}

func toDegrees(rad float64) float64 {
	return float64(rad) * (180.0 / math.Pi)
}

func toRadians(deg float64) float64 {
	return float64(deg) * (math.Pi / 180.0)
}

func elevation(date time.Time, latitude float64, longitude float64) float64 {
	return math.Asin(math.Sin(toRadians(latitude))*math.Sin(decl(date)) + math.Cos(toRadians(latitude))*math.Cos(decl(date))*math.Cos(toRadians(hA(date, longitude))))
}

func zenith(date time.Time, latitude float64, longitude float64) float64 {
	return math.Acos(math.Sin(toRadians(latitude))*math.Sin(decl(date)) + math.Cos(toRadians(latitude))*math.Cos(decl(date))*math.Cos(toRadians(hA(date, longitude))))
}

func hASunrise(date time.Time, latitude float64) float64 {
	return toDegrees(-math.Acos(math.Cos(toRadians(90.833))/(math.Cos(toRadians(latitude))*math.Cos(decl(date))) - (math.Tan(toRadians(latitude)) * math.Tan(decl(date)))))
}

func hASunset(date time.Time, latitude float64) float64 {
	return toDegrees(math.Acos(math.Cos(toRadians(90.833))/(math.Cos(toRadians(latitude))*math.Cos(decl(date))) - (math.Tan(toRadians(latitude)) * math.Tan(decl(date)))))
}

func sunrise(date time.Time, latitude float64, longitude float64) time.Time {
	var _, offset = date.Zone()
	sunrise := 720 - 4*(longitude-hASunrise(date, latitude)) - eqTime(date) + math.Round(float64(offset)/60)
	iHour, fHour := math.Modf(sunrise / 60)
	iMinute, fMinute := math.Modf(fHour * 60)
	iSecond, _ := math.Modf(fMinute * 60)
	return time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, date.Location())

}

func sunset(date time.Time, latitude float64, longitude float64) time.Time {
	var _, offset = date.Zone()
	sunset := 720 - 4*(longitude-hASunset(date, latitude)) - eqTime(date) + math.Round(float64(offset)/60)
	iHour, fHour := math.Modf(sunset / 60)
	iMinute, fMinute := math.Modf(fHour * 60)
	iSecond, _ := math.Modf(fMinute * 60)
	return time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, date.Location())
}

func solarNoon(date time.Time, longitude float64) time.Time {
	var _, offset = date.Zone()
	noon := 720 - 4*longitude - eqTime(date) + math.Round(float64(offset)/60)
	iHour, fHour := math.Modf(noon / 60)
	iMinute, fMinute := math.Modf(fHour * 60)
	iSecond, _ := math.Modf(fMinute * 60)
	return time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, date.Location())
}

func solarMidnight(date time.Time, longitude float64) time.Time {
	var _, offset = date.Zone()
	midnight := -4*longitude - eqTime(date) + math.Round(float64(offset)/60)
	iHour, fHour := math.Modf(midnight / 60)
	iMinute, fMinute := math.Modf(fHour * 60)
	iSecond, _ := math.Modf(fMinute * 60)
	return time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, date.Location())

}

func solarNoonElevation(date time.Time, latitude float64, longitude float64) float64 {
	return elevation(solarNoon(date, longitude), latitude, longitude)
}

func solarMidnightElevation(date time.Time, latitude float64, longitude float64) float64 {
	return elevation(solarMidnight(date, longitude), latitude, longitude)
}

func azimuth(date time.Time, latitude float64, longitude float64) float64 {
	solarMidnight := solarMidnight(date, longitude)
	solarNoon := solarNoon(date, longitude)
	midnight := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	oneDay, _ := time.ParseDuration("24h")
	solarMidnightSameDay := solarMidnight.After(midnight)
	solarNoonSameDay := solarNoon.Before(midnight.Add(oneDay))
	beforeMidnight := date.Before(time.Date(date.Year(), date.Month(), date.Day(), solarMidnight.Hour(), solarMidnight.Minute(), solarMidnight.Second(), 0, date.Location()))
	beforeNoon := date.Before(time.Date(date.Year(), date.Month(), date.Day(), solarNoon.Hour(), solarNoon.Minute(), solarNoon.Second(), 0, date.Location()))
	var azimuth float64
	if solarMidnightSameDay && solarNoonSameDay {
		if beforeMidnight && beforeNoon {
			azimuth = -math.Acos((math.Sin(decl(date))-math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude)))/(math.Cos(toRadians(latitude))*math.Sin(zenith(date, latitude, longitude)))) + 2*math.Pi
		} else if !beforeMidnight && beforeNoon {
			azimuth = math.Acos((math.Sin(decl(date)) - math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude))) / (math.Cos(toRadians(latitude)) * math.Sin(zenith(date, latitude, longitude))))
		} else if !beforeMidnight && !beforeNoon {
			azimuth = -math.Acos((math.Sin(decl(date))-math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude)))/(math.Cos(toRadians(latitude))*math.Sin(zenith(date, latitude, longitude)))) + 2*math.Pi
		}
	} else if !solarMidnightSameDay && solarNoonSameDay {
		if beforeMidnight && beforeNoon {
			azimuth = math.Acos((math.Sin(decl(date)) - math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude))) / (math.Cos(toRadians(latitude)) * math.Sin(zenith(date, latitude, longitude))))
		} else if beforeMidnight && !beforeNoon {
			azimuth = -math.Acos((math.Sin(decl(date))-math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude)))/(math.Cos(toRadians(latitude))*math.Sin(zenith(date, latitude, longitude)))) + 2*math.Pi
		} else if !beforeMidnight && !beforeNoon {
			azimuth = math.Acos((math.Sin(decl(date)) - math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude))) / (math.Cos(toRadians(latitude)) * math.Sin(zenith(date, latitude, longitude))))
		}
	} else if solarMidnightSameDay && !solarNoonSameDay {
		if beforeMidnight && beforeNoon {
			azimuth = math.Acos((math.Sin(decl(date)) - math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude))) / (math.Cos(toRadians(latitude)) * math.Sin(zenith(date, latitude, longitude))))
		} else if beforeMidnight && !beforeNoon {
			azimuth = -math.Acos((math.Sin(decl(date))-math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude)))/(math.Cos(toRadians(latitude))*math.Sin(zenith(date, latitude, longitude)))) + 2*math.Pi
		} else if !beforeMidnight && !beforeNoon {
			azimuth = math.Acos((math.Sin(decl(date)) - math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude))) / (math.Cos(toRadians(latitude)) * math.Sin(zenith(date, latitude, longitude))))
		}
	}
	return azimuth
}

func percentageElevationDay(date time.Time, latitude float64, longitude float64) float64 {
	maxElevation := toDegrees(solarNoonElevation(date, latitude, longitude))
	minElevation := -0.833
	actualElevation := toDegrees(elevation(date, latitude, longitude))
	return (actualElevation - minElevation) / (maxElevation - minElevation)
}

func percentageElevationCivilTwilight(date time.Time, latitude float64, longitude float64) float64 {
	maxElevation := -0.833
	minElevation := float64(-6)
	actualElevation := toDegrees(elevation(date, latitude, longitude))
	return (actualElevation - minElevation) / (maxElevation - minElevation)
}

func percentageElevationNauticalTwilight(date time.Time, latitude float64, longitude float64) float64 {
	maxElevation := float64(-6)
	minElevation := float64(-12)
	actualElevation := toDegrees(elevation(date, latitude, longitude))
	return (actualElevation - minElevation) / (maxElevation - minElevation)
}

func colorTemp(date time.Time, latitude float64, longitude float64) int64 {
	actualElevation := toDegrees(elevation(date, latitude, longitude))
	if actualElevation > -0.833 {
		return int64(math.Round(percentageElevationDay(date, latitude, longitude)*2500 + 3000))
	} else if actualElevation <= -0.833 && actualElevation > -6 {
		return int64(math.Round(percentageElevationCivilTwilight(date, latitude, longitude)*1000 + 2000))
	} else {
		return 2000
	}
}

func brightness(date time.Time, latitude float64, longitude float64) int64 {
	actualElevation := toDegrees(elevation(date, latitude, longitude))
	if actualElevation > -6 {
		return 100
	} else if actualElevation <= -6 && actualElevation > -12 {
		return int64(math.Round(percentageElevationNauticalTwilight(date, latitude, longitude)*50 + 50))
	} else {
		return 50
	}
}

func main() {

	date := time.Now()
	latitude := 48.87
	longitude := 2.67

	dateSunrise := sunrise(date, latitude, longitude)
	dateNoon := solarNoon(date, longitude)
	dateSunset := sunset(date, latitude, longitude)
	dateMidnight := solarMidnight(date, longitude)

	fmt.Printf("Solar midnight: %v\n", dateMidnight)
	fmt.Printf("Sunrise: %v\n", dateSunrise)
	fmt.Printf("Solar noon: %v\n", dateNoon)
	fmt.Printf("Sunset: %v\n", dateSunset)

	fmt.Printf("Sun position at midnight: (%f,%f)\n", toDegrees(azimuth(dateMidnight, latitude, longitude)), toDegrees(elevation(dateMidnight, latitude, longitude)))
	fmt.Printf("Sun position at sunrise: (%f,%f)\n", toDegrees(azimuth(dateSunrise, latitude, longitude)), toDegrees(elevation(dateSunrise, latitude, longitude)))
	fmt.Printf("Sun position at noon: (%f,%f)\n", toDegrees(azimuth(dateNoon, latitude, longitude)), toDegrees(elevation(dateNoon, latitude, longitude)))
	fmt.Printf("Sun position at sunset: (%f,%f)\n", toDegrees(azimuth(dateSunset, latitude, longitude)), toDegrees(elevation(dateSunset, latitude, longitude)))

	fmt.Printf("Circadian color temperature at midnight: %d kelvin and brightness %d%%\n", colorTemp(dateMidnight, latitude, longitude), brightness(dateMidnight, latitude, longitude))
	fmt.Printf("Circadian color temperature at sunrise: %d kelvin and brightness %d%%\n", colorTemp(dateSunrise, latitude, longitude), brightness(dateSunrise, latitude, longitude))
	fmt.Printf("Circadian color temperature at noon: %d kelvin and brightness %d%%\n", colorTemp(dateNoon, latitude, longitude), brightness(dateNoon, latitude, longitude))
	fmt.Printf("Circadian color temperature at sunset: %d kelvin and brightness %d%%\n", colorTemp(dateSunset, latitude, longitude), brightness(dateSunset, latitude, longitude))

	for h := 0; h < 24; h++ {
		for m := 0; m < 60; m = m + 10 {
			d := time.Date(2021, 1, 1, h, m, 0, 0, time.Now().UTC().Location())
			fmt.Printf("%v : %f, %f\n", d, toDegrees(azimuth(d, latitude, longitude)), toDegrees(elevation(d, latitude, longitude)))
		}
	}

}
