package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Location struct {
	City    string  `json:"city"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type Weather struct {
	TempC float64 `json:"temp_c"`
	Text  string  `json:"condition:text"`
}

type TimezoneResponse struct {
	Timezone   string `json:"timezone"`
	RawOffset  int    `json:"raw_offset"`
	DstOffset  int    `json:"dst_offset"`
	TimeOffset int    `json:"time_offset"`
}

func main() {
	location, err := getCurrentLocation()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Your current location: %s, %s\n", location.City, location.Country)

	weather, err := getCurrentWeather(location.Lat, location.Lon)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current weather: %.1f°C, %s\n", weather.TempC, weather.Text)

	localTime, err := getCurrentLocalTime(location.Lat, location.Lon)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Local time: %s\n", localTime.Format(time.RFC3339))
}

func getCurrentLocation() (Location, error) {
	res, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return Location{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Location{}, err
	}

	var location Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		return Location{}, err
	}

	return location, nil
}

func getCurrentWeather(lat, lon float64) (Weather, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=ec180872243c4f57a4f153631230105&q=Jaipur&days=1&aqi=no&alerts=no")
	res, err := http.Get(url)
	if err != nil {
		return Weather{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Weather{}, err
	}

	var weatherData struct {
		Current struct {
			TempC     float64 `json:"temp_c"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
		} `json:"current"`
	}

	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return Weather{}, err
	}

	weather := Weather{
		TempC: weatherData.Current.TempC,
		Text:  weatherData.Current.Condition.Text,
	}

	return weather, nil
}

func getCurrentLocalTime(lat, lon float64) (time.Time, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/timezone/json?location=%.6f,%.6f&timestamp=%d&key=%s", lat, lon, time.Now().Unix())
	res, err := http.Get(url)
	if err != nil {
		return time.Time{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return time.Time{}, err
	}

	var timezoneResponse TimezoneResponse
	err = json.Unmarshal(body, &timezoneResponse)
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now()
	secsOffset := int64(timezoneResponse.RawOffset + timezoneResponse.TimeOffset)
	localTime := now.Add(time.Duration(secsOffset) * time.Second)

	return localTime, nil
}
