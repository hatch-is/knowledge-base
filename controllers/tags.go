package controllers

import (
	"knowledge-base/models"

	"github.com/gin-gonic/gin"
)

//TagsController ...
type TagsController struct{}

var tagsModel = new(models.TagsModel)

//All get all tags as []string
func (ctrl TagsController) All(c *gin.Context) {
	tags := tagsModel.All()
	c.JSON(200, tags)
}
