package models

type (
	Weather struct {
		CurrentCondition []CurrentCondition `json:"current_condition"`
		WeatherForecast  []WeatherForecast  `json:"weather"`
	}

	CurrentCondition struct {
		FeelsLikeC  string `json:"FeelsLikeC"`
		FeelsLikeF  string `json:"FeelsLikeF"`
		TempC       string `json:"temp_C"`
		TempF       string `json:"temp_F"`
		WeatherDesc []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
	}

	WeatherForecast struct {
		Date        string `json:"date"`
		MaxTempC    string `json:"maxtempC"`
		MaxTempF    string `json:"maxtempF"`
		MinTempC    string `json:"mintempC"`
		MinTempF    string `json:"mintempF"`
		TotalSnowCm string `json:"totalSnow_cm"`
		UvIndex     string `json:"uvIndex"`
	}

	WeatherResponse struct {
		CurrentCondition CurrentConditionResponse  `json:"current_condition"`
		WeatherForecast  []WeatherForecastResponse `json:"weatherForecast"`
	}

	CurrentConditionResponse struct {
		FeelsLike          string `json:"feelsLike"`
		Temperature        string `json:"temperature"`
		WeatherDescription string `json:"weatherDescription"`
	}

	WeatherForecastResponse struct {
		Date        string `json:"date"`
		MaxTemp     string `json:"maxTemperature"`
		MinTemp     string `json:"minTemperature"`
		TotalSnowCm string `json:"totalSnow_cm"`
		UvIndex     string `json:"uvIndex"`
	}
)
