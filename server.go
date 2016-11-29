package main

import (
	"knowledge-base/controllers"
	"knowledge-base/db"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db.Init()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Message": "Hello and welcome to the Hatch Knowledge Base"})
	})
	knowledge := r.Group("/knowledge")
	{
		knowledge.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"Message": "pong",
			})
		})

		article := new(controllers.ArticleController)

		knowledge.POST("/articles", article.Create)
		knowledge.GET("/articles", article.Read)
		knowledge.GET("/articles/:id", article.ReadOne)
		knowledge.DELETE("articles/:id", article.Delete)
		knowledge.PUT("articles/:id", article.Update)

		tags := new(controllers.TagsController)
		knowledge.GET("/tags", tags.All)
	}
	r.Run()
}
