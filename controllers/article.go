package controllers

import (
	"knowledge-base/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

//ArticleController ...
type ArticleController struct{}

var articleModel = new(models.ArticleModel)

//Create add new article
func (ctrl ArticleController) Create(c *gin.Context) {
	var article models.Article
	if err := c.BindJSON(&article); err != nil {
		c.JSON(400, c.Error(err).SetType(gin.ErrorTypeBind))
		c.Abort()
		return
	}

	err := article.Create()

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(201, article)
}

//Read get articles
//TODO: return all count and current count of entries
func (ctrl ArticleController) Read(c *gin.Context) {
	filter := c.Query("filter")
	offset := c.Query("skip")
	limit := c.Query("limit")
	_, offsetErr := strconv.Atoi(offset)
	_, limitErr := strconv.Atoi(limit)

	if limitErr != nil {
		offset = "NULL"
		limit = "NULL"
	}
	if offsetErr != nil {
		offset = "NULL"
	}
	articles, err, count, left := articleModel.Read(filter, limit, offset)
	if err != nil {
		c.JSON(404, c.Error(err))
		c.Abort()
		return
	}
	c.Writer.Header().Set("X-Total-Count", strconv.Itoa(count))
	c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(left))
	c.JSON(200, articles)
}

//ReadOne get article by ID
func (ctrl ArticleController) ReadOne(c *gin.Context) {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, c.Error(err))
		c.Abort()
		return
	}
	article, err := articleModel.ReadOne(ID)
	if err != nil {
		c.JSON(404, c.Error(err))
		c.Abort()
		return
	}
	c.JSON(200, article)
}

//Delete article by ID
func (ctrl ArticleController) Delete(c *gin.Context) {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, c.Error(err))
		c.Abort()
		return
	}

	err = articleModel.Delete(ID)
	if err != nil {
		c.JSON(404, c.Error(err))
		c.Abort()
		return
	}
	c.JSON(202, gin.H{"message": "Article with ID [" + strconv.FormatInt(ID, 10) + "] deleted"})
}

//Update article by ID
func (ctrl ArticleController) Update(c *gin.Context) {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	var article models.Article
	if err = c.BindJSON(&article); err != nil {
		c.JSON(400, c.Error(err).SetType(gin.ErrorTypeBind))
		c.Abort()
		return
	}

	err = article.Update(ID)
	if err != nil {
		c.JSON(404, c.Error(err))
		c.Abort()
		return
	}
	c.JSON(201, article)
}
