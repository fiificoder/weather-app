package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}

	var configFile apiConfigData

	err = json.Unmarshal(bytes, &configFile)
	if err != nil {
		return apiConfigData{}, err
	}

	return configFile, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Gophers"))
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig("json.apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()

	var forecast weatherData
	err = json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		return weatherData{}, nil
	}
	return forecast, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json utf=8")
			json.NewEncoder(w).Encode(data)
		})

	if http.ListenAndServe(":8000", nil); http.ListenAndServe(":8080", nil) != nil {
		fmt.Printf(" Cannot Start Server at port 8080\n")
	}
}
