package main

import (
	"github.com/gin-gonic/gin"
	topis "golony6449/PubTrans4Watch_Server/api/topis"
	"log"
)

func main() {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/ArrivalInfo/:station_name", topis.GetStationArrivalInfo)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
