package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func parseResp(resp *http.Response) map[string]interface{} {
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	jsonObj := make(map[string]interface{})
	json.Unmarshal(data, &jsonObj)

	loc := jsonObj["name"].(string)
	temp := jsonObj["main"].(map[string]interface{})["temp"]
	weather := jsonObj["weather"].([]interface{})[0].(map[string]interface{})["main"]

	weatherMap := map[string]interface{}{"Location": loc, "Temp": temp, "Weather": weather}
	return weatherMap
}

func QueryWeatherApi(long float64, lat float64) map[string]interface{} {
	apiKey := os.Getenv("weatherAPI")
	reqUrl := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric",
		lat, long, apiKey,
	)
	fmt.Println(reqUrl)
	res, err := http.Get(reqUrl)

	if err != nil {
		log.Fatalln(err)
	}

	weatherMap := parseResp(res)
	return weatherMap
}
