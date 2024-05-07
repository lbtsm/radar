package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/filter/internal/api/handler"
	"github.com/mapprotocol/filter/internal/api/store/mysql"
	"github.com/pkg/errors"
	"time"
)

func initMiddleware(g *gin.Engine) {
	g.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           2 * time.Minute,
	}))
}

func initController(g *gin.Engine, dsn string) error {
	db, err := mysql.Init(dsn)
	if err != nil {
		return errors.Wrap(err, "init db")
	}
	v1 := g.Group("/v1")
	{
		pro := handler.NewProject(db)
		group := v1.Group("project")
		group.GET("", pro.Get)
		group.POST("", pro.Add)
	}
	{
		event := handler.NewEvent(db)
		group := v1.Group("event")
		group.GET("", event.Get)
		group.POST("", event.Add)
		group.DELETE("", event.Delete)
		group.GET("/list", event.List)
	}
	return nil
}
