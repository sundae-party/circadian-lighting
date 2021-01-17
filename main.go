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

//TODO add azimuth func

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
	date := time.Date(2021, 1, 17, 17, 44, 0, 0, &time.Location{})
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

	fmt.Printf("Sunrise: %v\n", dateSunrise)
	fmt.Printf("Solar noon: %v\n", dateNoon)
	fmt.Printf("Sunset: %v\n", dateSunset)
	fmt.Printf("Circadian color temperature: %d kelvin and brightness %d%%\n", colorTemp(date, latitude, longitude), brightness(date, latitude, longitude))

}
