package PubTrans4Watch_Server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetStationArrivalInfo(context *gin.Context) {
	stationName := context.Param("station_name")
	url := fmt.Sprintf("http://swopenAPI.seoul.go.kr/api/subway/%s/json/realtimeStationArrival/0/5/%s", os.Getenv("SECRET"), stationName)

	if stationName == "" {
		err := errors.New("empty parameter: station_name")
		log.Fatal(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	//defer resp.Body.Close() // TODO 모든 처리 후 Body Close 처리

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	jsonData := make(map[string]interface{})

	// TODO: 예외처리
	json.Unmarshal(data, &jsonData)

	statusCode := 0
	if jsonData["code"] == "INFO-000" {
		statusCode = 200
	} else {
		statusCode = 500
		fmt.Println("TOPIS에서 오류 응답")
	}

	context.JSON(statusCode, gin.H{
		"data": jsonData,
	})
}
