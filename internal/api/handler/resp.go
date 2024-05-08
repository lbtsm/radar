package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/filter/internal/api/stream"
	"net/http"
)

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil { // todo custom error
		c.JSON(http.StatusOK, stream.CommonResp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, stream.CommonResp{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}
