package PubTrans4Watch_Server

import "github.com/gin-gonic/gin"

func GetStationArrivalInfo(context *gin.Context) {
	context.JSON(200, gin.H{
		"test": "test",
	})
}
