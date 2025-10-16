package binders

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, obj any) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return err
	}

	return nil
}

func BindXML(c *gin.Context, obj any) error {
	if err := c.ShouldBindXML(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return err
	}

	return nil
}

func BindForm(c *gin.Context, obj any) error {
	if err := c.ShouldBind(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return err
	}

	return nil
}
