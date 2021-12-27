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
	"strconv"
)

type RemainSec struct {
	idx       int
	remainSec int
}

// Insert TODO: Generic 함수로 변경 후 공통으로 분리 (현재 Beta 단계)
func Insert(origin []RemainSec, idx int, target RemainSec) []RemainSec {
	fmt.Printf("Insert function: idx: %d   target: %+v\n", idx, target)
	fmt.Print("before: ")
	fmt.Println(origin)

	if len(origin) == 0 || len(origin) == idx {
		return append(origin, target)
	}

	result := make([]RemainSec, 1)
	copy(result, origin[:idx])

	//fmt.Println("result")
	//fmt.Println(result)

	result = append(result, target)

	for _, elem := range origin[idx:] {
		result = append(result, elem)
	}

	//fmt.Println(result)

	return result
}

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

	recvData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 응답 데이터 Json 변환
	originJsonData := make(map[string]interface{})
	orders := make([]RemainSec, 0)

	// TODO: 예외처리
	json.Unmarshal(recvData, &originJsonData)

	for idx, elem := range originJsonData["realtimeArrivalList"].([]interface{}) {
		fmt.Printf("%d 번째\n", idx)
		temp := elem.(map[string]interface{})
		remainSecond, _ := strconv.Atoi(temp["barvlDt"].(string))

		// 추가대상 데이터
		newOrder := new(RemainSec)
		newOrder.idx = idx
		newOrder.remainSec = remainSecond

		fmt.Printf("newOrder: %d  %d\n", newOrder.idx, newOrder.remainSec)
		//println(newOrder)

		if len(orders) == 0 {
			orders = Insert(orders, 0, *newOrder)
			fmt.Println("Init")
			fmt.Println(orders)
		} else {

			for orderIdx, o := range orders {
				if o.remainSec > remainSecond {
					fmt.Printf("insert: %d\n", orderIdx)
					orders = Insert(orders, orderIdx, *newOrder)
					fmt.Println("after result")
					fmt.Println(orders)
					break
				} else if orderIdx == len(orders)-1 {
					fmt.Printf("append\n")
					orders = Insert(orders, len(orders), *newOrder)
					fmt.Println("after result")
					fmt.Println(orders)
					break
				}
			}
		}

		//// TODO 임시코드
		//for _, inner_elem := range temp{
		//	if inner_elem != nil{
		//		fmt.Print(inner_elem)
		//		fmt.Print("    ")
		//	}
		//}
		//
		//fmt.Println()
	}

	fmt.Println("원본 데이터 건수")
	fmt.Println(len(originJsonData["realtimeArrivalList"].([]interface{})))
	fmt.Println("정렬 결과")
	fmt.Println(orders)

	// TODO: 정렬된 데이터순으로 map 데이터 추가
	orderedJsonData := make(map[string]interface{})
	orderedJsonData["errorMessage"] = originJsonData["errorMessage"]
	realtimeArrivalList := make([]interface{}, 0)

	for _, elem := range orders {
		realtimeArrivalList = append(realtimeArrivalList, originJsonData["realtimeArrivalList"].([]interface{})[elem.idx])
	}

	orderedJsonData["realtimeArrivalList"] = realtimeArrivalList
	fmt.Print(orderedJsonData)
	// 응답처리
	statusCode := 200
	// TODO TOPIS API 예외처리
	// 성공/실패시 응답형태 다름

	context.JSON(statusCode, gin.H{
		"data": orderedJsonData,
	})
}
