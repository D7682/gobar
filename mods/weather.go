package mods

import (
	barista "barista.run"
	"barista.run/bar"
	"barista.run/base/click"
	"barista.run/modules/weather"
	"barista.run/modules/weather/openweathermap"
	"barista.run/outputs"
	"barista.run/pango"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

func GetWeather() {
	var openWeatherInBrowser = click.RunLeft(os.Getenv("BROWSER"), fmt.Sprintf("https://openweathermap.org/city/%s", viper.GetString("openweather.cityid")))
	weatherModule := weather.New(openweathermap.New(viper.GetString("openweather.key")).CityID(viper.GetString("openweather.cityid"))).Output(func(w weather.Weather) bar.Output {
		var iconName string
		switch w.Condition {
		case weather.Thunderstorm,
			weather.TropicalStorm,
			weather.Hurricane:
			iconName = "stormy"
		case weather.Drizzle,
			weather.Hail:
			iconName = "shower"
		case weather.Rain:
			iconName = "downpour"
		case weather.Mist,
			weather.Smoke,
			weather.Whirls,
			weather.Haze,
			weather.Fog:
			iconName = "windy"
		case weather.Clear:
			if !w.Sunset.IsZero() && time.Now().After(w.Sunset) {
				iconName = "night"
			} else {
				iconName = "sunny"
			}
		case weather.PartlyCloudy:
			iconName = "partly-cloudy"
		case weather.Tornado,
			weather.Windy:
			iconName = "windy"
		}

		if iconName == "" {
			iconName = "warning-outline"
		} else {
			iconName = "weather-" + iconName
		}

		return outputs.Pango(
			pango.Icon("typecn-"+iconName), spacer,
			pango.Textf("%.1f°F", w.Temperature.Fahrenheit()),
			pango.Textf(" (provided by %s)", w.Attribution).XSmall(),
		).OnClick(openWeatherInBrowser)
	})

	barista.Add(weatherModule)
}
