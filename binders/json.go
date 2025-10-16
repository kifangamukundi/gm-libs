package binders

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SuccessMessage  = "Action successfully"
	ItemKey         = "item"
	ImagesKey       = "images"
	UserKey         = "user"
	AccessTokenKey  = "accessToken"
	RefreshTokenKey = "refreshToken"
	TotalCountKey   = "totalCount"
	CountKey        = "count"
	PageKey         = "page"
	LimitKey        = "limit"
	ItemsKey        = "items"
	OtherKey        = "other"
)

func ReturnJSONResponse(c *gin.Context, status int, success bool, data any) {
	c.JSON(status, gin.H{
		"success": success,
		"data":    data,
	})
}

func returnSuccessResponse(c *gin.Context, status int, data any) {
	ReturnJSONResponse(c, status, true, data)
}

type PaginatedResponse struct {
	TotalCount int `json:"totalCount"`
	Count      int `json:"count"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Items      any `json:"items"`
}

func ReturnJSONPaginateResponse(c *gin.Context, page, limit, totalCount, count int, data any) {
	paginatedResponse := PaginatedResponse{
		TotalCount: totalCount,
		Count:      count,
		Page:       page,
		Limit:      limit,
		Items:      data,
	}

	returnSuccessResponse(c, http.StatusOK, paginatedResponse)
}

func ReturnJSONGeneralResponse(c *gin.Context, data any) {
	returnSuccessResponse(c, http.StatusOK, gin.H{ItemKey: data})
}

func ReturnJSONCacheResponse(c *gin.Context, data any) {
	returnSuccessResponse(c, http.StatusOK, data)
}

func ReturnJSONBlogPostResponse(c *gin.Context, data, other any) {
	returnSuccessResponse(c, http.StatusOK, gin.H{ItemKey: data, OtherKey: other})
}

func ReturnJSONPermissionsResponse(c *gin.Context, data any) {
	returnSuccessResponse(c, http.StatusOK, data)
}

func ReturnJSONCreatedGenericResponse(c *gin.Context) {
	returnSuccessResponse(c, http.StatusCreated, SuccessMessage)
}

func ReturnJSONUpdatedGenericResponse(c *gin.Context) {
	returnSuccessResponse(c, http.StatusOK, SuccessMessage)
}

func ReturnJSONOkayGenericResponse(c *gin.Context) {
	returnSuccessResponse(c, http.StatusOK, SuccessMessage)
}

type TokenResponse struct {
	User struct {
		AccessToken  any `json:"accessToken"`
		RefreshToken any `json:"refreshToken"`
	} `json:"user"`
}

func ReturnJSONTokenResponse(c *gin.Context, access, refresh any) {
	tokenResponse := TokenResponse{}

	tokenResponse.User.AccessToken = access
	tokenResponse.User.RefreshToken = refresh

	returnSuccessResponse(c, http.StatusCreated, tokenResponse)
}

func ReturnJSONMediaUploadResponse(c *gin.Context, images any) {
	returnSuccessResponse(c, http.StatusCreated, gin.H{ImagesKey: images})
}
