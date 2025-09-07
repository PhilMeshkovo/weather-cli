package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const apiKey = "36bd95dcee3866f9d399083d5a9cf83c"

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func getWeather(city string) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ru", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}
	if weather.Name == "" {
		return nil, fmt.Errorf("город не найден")
	}
	return &weather, nil
}

func main() {
	var cmdWeather = &cobra.Command{
		Use:   "weather [город]",
		Short: "Показать текущую погоду в городе",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			city := args[0]
			weather, err := getWeather(city)
			if err != nil {
				fmt.Println("Ошибка:", err)
				os.Exit(1)
			}
			fmt.Printf("Погода в городе %s:\nТемпература: %.1f°C\nСостояние: %s\n",
				weather.Name, weather.Main.Temp, weather.Weather[0].Description)
		},
	}
	cmdWeather.Execute()
}
