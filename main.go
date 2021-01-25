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
	if isLeapYear(date) {
		return (2 * math.Pi * (float64(date.YearDay()) - 1 + float64(date.Hour())/24 + float64(date.Minute())/(24*60) + float64(date.Second())/(24*60*60)) / 366)
	} else {
		return (2 * math.Pi * (float64(date.YearDay()) - 1 + float64(date.Hour())/24 + float64(date.Minute())/(24*60) + float64(date.Second())/(24*60*60)) / 365)
	}
}

func eqTime(date time.Time) float64 {
	return 229.18 * (0.000075 + 0.001868*math.Cos(fractionalYear(date)) - 0.032077*math.Sin(fractionalYear(date)) - 0.014615*math.Cos(2*fractionalYear(date)) - 0.040849*math.Sin(2*fractionalYear(date)))
}

func decl(date time.Time) float64 {
	return 0.006918 - 0.399912*math.Cos(fractionalYear(date)) + 0.070257*math.Sin(fractionalYear(date)) - 0.006758*math.Cos(2*fractionalYear(date)) + 0.000907*math.Sin(2*fractionalYear(date)) - 0.002697*math.Cos(3*fractionalYear(date)) + 0.00148*math.Sin(3*fractionalYear(date))
}

func timeOffset(date time.Time, longitude float64) float64 {
	var _, offset = time.Now().Zone()
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

func sunrise(date time.Time, latitude float64, longitude float64) float64 {
	var _, offset = time.Now().Zone()
	return 720 - 4*(longitude-hASunrise(date, latitude)) - eqTime(date) + math.Round(float64(offset)/60)
}

func sunset(date time.Time, latitude float64, longitude float64) float64 {
	var _, offset = time.Now().Zone()
	return 720 - 4*(longitude-hASunset(date, latitude)) - eqTime(date) + math.Round(float64(offset)/60)
}

func solarNoon(date time.Time, longitude float64) float64 {
	var _, offset = time.Now().Zone()
	return 720 - 4*longitude - eqTime(date) + math.Round(float64(offset)/60)
}

func solarMidnight(date time.Time, longitude float64) float64 {
	var _, offset = time.Now().Zone()
	return -4*longitude - eqTime(date) + math.Round(float64(offset)/60)
}

func solarNoonElevation(date time.Time, latitude float64, longitude float64) float64 {
	noon := solarNoon(date, longitude)
	iHour, fHour := math.Modf(noon / 60)
	iMinute, fMinute := math.Modf(fHour * 60)
	iSecond, _ := math.Modf(fMinute * 60)
	dateNoon := time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, &time.Location{})
	return elevation(dateNoon, latitude, longitude)
}

func azimuth(date time.Time, latitude float64, longitude float64) float64 {
	dateSeconds := date.Hour()*60 + date.Minute() + date.Second()/60
	midnight := solarMidnight(date, longitude)
	noon := solarNoon(date, longitude)
	midnight0 := midnight > 0
	noon1440 := noon < 1440
	beforeMidnight := (dateSeconds - int(midnight)%1440) < 0
	beforeNoon := (dateSeconds - int(noon)%1440) < 0
	var azimuth float64
	if midnight0 && noon1440 {
		if beforeMidnight && beforeNoon {
			azimuth = -math.Acos((math.Sin(decl(date))-math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude)))/(math.Cos(toRadians(latitude))*math.Sin(zenith(date, latitude, longitude)))) + 2*math.Pi
		} else if !beforeMidnight && beforeNoon {
			azimuth = math.Acos((math.Sin(decl(date)) - math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude))) / (math.Cos(toRadians(latitude)) * math.Sin(zenith(date, latitude, longitude))))
		} else if !beforeMidnight && !beforeNoon {
			azimuth = -math.Acos((math.Sin(decl(date))-math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude)))/(math.Cos(toRadians(latitude))*math.Sin(zenith(date, latitude, longitude)))) + 2*math.Pi
		}
	} else if !midnight0 && noon1440 {
		if beforeMidnight && beforeNoon {
			azimuth = math.Acos((math.Sin(decl(date)) - math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude))) / (math.Cos(toRadians(latitude)) * math.Sin(zenith(date, latitude, longitude))))
		} else if beforeMidnight && !beforeNoon {
			azimuth = -math.Acos((math.Sin(decl(date))-math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude)))/(math.Cos(toRadians(latitude))*math.Sin(zenith(date, latitude, longitude)))) + 2*math.Pi
		} else if !beforeMidnight && !beforeNoon {
			azimuth = math.Acos((math.Sin(decl(date)) - math.Sin(toRadians(latitude))*math.Cos(zenith(date, latitude, longitude))) / (math.Cos(toRadians(latitude)) * math.Sin(zenith(date, latitude, longitude))))
		}
	} else if midnight0 && !noon1440 {
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

func solarMidnightElevation(date time.Time, latitude float64, longitude float64) float64 {
	noon := solarMidnight(date, longitude)
	iHour, fHour := math.Modf(noon / 60)
	iMinute, fMinute := math.Modf(fHour * 60)
	iSecond, _ := math.Modf(fMinute * 60)
	dateNoon := time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, &time.Location{})
	return elevation(dateNoon, latitude, longitude)
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
	latitude := 43.575997
	longitude := 1.480259

	d := sunrise(date, latitude, longitude)
	iHour, fHour := math.Modf(d / 60)
	iMinute, fMinute := math.Modf(fHour * 60)
	iSecond, _ := math.Modf(fMinute * 60)
	dateSunrise := time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, &time.Location{})

	d = solarNoon(date, longitude)
	iHour, fHour = math.Modf(d / 60)
	iMinute, fMinute = math.Modf(fHour * 60)
	iSecond, _ = math.Modf(fMinute * 60)
	dateNoon := time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, &time.Location{})

	d = sunset(date, latitude, longitude)
	iHour, fHour = math.Modf(d / 60)
	iMinute, fMinute = math.Modf(fHour * 60)
	iSecond, _ = math.Modf(fMinute * 60)
	dateSunset := time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, &time.Location{})

	d = solarMidnight(date, longitude)
	iHour, fHour = math.Modf(d / 60)
	iMinute, fMinute = math.Modf(fHour * 60)
	iSecond, _ = math.Modf(fMinute * 60)
	dateMidnight := time.Date(date.Year(), date.Month(), date.Day(), int(iHour), int(iMinute), int(iSecond), 0, &time.Location{})

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

}
