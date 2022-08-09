package Routing

import (
	"Newism/Controllers"
	"Newism/Middleware"
	"github.com/gin-gonic/gin"
)

func SetRouting(e *gin.Engine) {
	Account := e.Group("/account/v1")

	Api := e.Group("/api/v1")
	Api.Use(Middleware.JwtAuthMiddleware())

	Api.POST("/CreatePost/", Controllers.CrNew)
	Api.POST("/GetPostsByTag/", Controllers.GetNewsByTag)
	Api.POST("/VerifyPost/", Controllers.VerfyNew)
	Api.POST("/LikePost/", Controllers.LikePost)
	Api.POST("/ReportPost/", Controllers.ReportPost)
	Api.DELETE("/DeletePost/", Controllers.DeletePost)
	Api.GET("/GetPostForAdmin/", Controllers.GetNotPublicPosts)
	Api.GET("/GetReports/", Controllers.GetReports)
	Api.GET("/GetAllPost/", Controllers.GetNews)
	Api.GET("/GetLikes/", Controllers.GetLikes)

	Account.POST("/Login/", Controllers.LoginUser)
	Account.POST("/Register/", Controllers.CreateUser)
}
