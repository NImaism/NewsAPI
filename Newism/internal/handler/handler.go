package handler

import (
	"Newism/internal/middleware"
	"Newism/internal/model"
	"Newism/internal/store"
	"Newism/internal/utility"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	Logger *zap.Logger
	Store  *store.Mognodb
}

func New(Logger *zap.Logger, Store *store.Mognodb) *App {
	return &App{
		Logger: Logger,
		Store:  Store,
	}
}

func (a *App) CrNew(e *gin.Context) {
	TokenData, Err := utility.ExtractTokenData(e)
	if Err != nil {
		a.Logger.Error("utility.extracttokendata Error In Extract Token Data")
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	Time := time.Now()

	Post := new(model.New)
	if err := e.Bind(Post); err != nil {
		a.Logger.Info("handler.app.bind failed", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}
	Post.Public = 0
	Post.Creator = TokenData["UserName"]
	Post.CreateDate = Time.String()

	Validator := validator.New()
	if err := Validator.Struct(Post); err != nil {
		a.Logger.Info("validator data has missing", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := a.Store.CreateNew(*Post)
	if err != nil {
		a.Logger.Error("store.createnew failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(http.StatusOK, model.SuccessResponse(id))
}

func (a *App) DeletePost(e *gin.Context) {
	TokenData, Err := utility.ExtractTokenData(e)
	if Err != nil {
		a.Logger.Error("utility.extracttokendata Error In Extract Token Data")
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if TokenData["IsAdmin"] == false {
		a.Logger.Info("IsAdmin In TokenData Has False")
		e.AbortWithStatus(http.StatusForbidden)
		return
	}

	err := a.Store.DeleteNew(e.Query("id"))
	if err != nil {
		a.Logger.Error("store.deletenew failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	e.JSON(200, model.SuccessResponse(nil))
}

func (a *App) GetNews(e *gin.Context) {
	Limit, Err := strconv.Atoi(e.Query("limit"))
	if Err != nil {
		a.Logger.Info("Error In Convert Limit To Int", zap.Error(Err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Data, err := a.Store.GetNews(Limit)
	if err != nil {
		a.Logger.Error("store.getnews failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(200, model.SuccessResponse(Data))
}

func (a *App) LoginUser(e *gin.Context) {
	LoginModel := new(model.LoginUserViewModel)

	if err := e.Bind(LoginModel); err != nil {
		a.Logger.Info("handler.app.bind failed", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Validator := validator.New()
	if err := Validator.Struct(LoginModel); err != nil {
		a.Logger.Info("validator data has missing", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Ok, IsAdmin, err := a.Store.Login(LoginModel.UserName, LoginModel.Password)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			a.Logger.Info("store.login failed", zap.Error(err))
			e.AbortWithStatus(http.StatusNotFound)
			return
		}
		a.Logger.Error("store.login failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if Ok != true {
		e.JSON(http.StatusOK, model.ErrorResponse(nil, "Account Not Found"))
		return
	}
	Token, err := utility.GenerateToken(LoginModel.UserName, IsAdmin)
	if err != nil {
		a.Logger.Error("utility.generatetoken failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(http.StatusOK, model.SuccessResponse(Token))
}

func (a *App) CreateUser(e *gin.Context) {
	User := new(model.User)
	err := e.Bind(User)
	if err != nil {
		a.Logger.Info("handler.app.bind failed", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}
	User.IsAdmin = 0

	Validator := validator.New()
	if err := Validator.Struct(User); err != nil {
		a.Logger.Info("validator data has missing", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = a.Store.CreateUser(*User)
	if err != nil {
		a.Logger.Error("store.createuser failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(http.StatusOK, model.SuccessResponse("Done !"))
}

func (a *App) GetNewsByTag(e *gin.Context) {
	Data := new(model.GetNewByTag)
	if err := e.Bind(Data); err != nil {
		a.Logger.Info("handler.app.bind failed", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Validator := validator.New()
	if err := Validator.Struct(Data); err != nil {
		a.Logger.Info("validator data has missing", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Dts, err := a.Store.GetNewsByTag(*Data)
	if err != nil {
		a.Logger.Error("store.getnewsbytag failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(200, model.SuccessResponse(Dts))
}

func (a *App) VerifyNew(e *gin.Context) {
	TokenData, Err := utility.ExtractTokenData(e)
	if Err != nil {
		a.Logger.Error("utility.extracttokendata Error In Extract Token Data")
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if TokenData["IsAdmin"] == false {
		a.Logger.Info("IsAdmin In TokenData Has False")
		e.AbortWithStatus(http.StatusForbidden)
		return
	}

	err := a.Store.VerifyPost(e.Query("id"))
	if err != nil {
		a.Logger.Error("store.verifypost failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(200, model.SuccessResponse("Done !"))
}

func (a *App) LikePost(e *gin.Context) {
	Data := new(model.LikePost)
	if err := e.Bind(Data); err != nil {
		a.Logger.Info("handler.app.bind failed", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Time := time.Now()

	TokenData, Err := utility.ExtractTokenData(e)
	if Err != nil {
		a.Logger.Error("utility.extracttokendata Error In Extract Token Data")
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	Data.By = TokenData["UserName"]
	Data.Time = Time.String()

	Validator := validator.New()
	if err := Validator.Struct(Data); err != nil {
		a.Logger.Info("validator data has missing", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := a.Store.LikePost(*Data); err != nil {
		a.Logger.Error("store.likepost failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(200, model.SuccessResponse(nil))
}

func (a *App) ReportPost(e *gin.Context) {
	Data := new(model.ReportPost)
	if err := e.Bind(Data); err != nil {
		a.Logger.Info("handler.app.bind failed", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Time := time.Now()

	TokenData, Err := utility.ExtractTokenData(e)
	if Err != nil {
		a.Logger.Error("utility.extracttokendata Error In Extract Token Data")
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	Data.By = TokenData["UserName"]
	Data.Time = Time.String()

	Validator := validator.New()
	if err := Validator.Struct(Data); err != nil {
		a.Logger.Info("validator data has missing", zap.Error(err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := a.Store.ReportPost(*Data); err != nil {
		a.Logger.Error("store.reportpost failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(200, model.SuccessResponse(nil))
}

func (a *App) GetNotPublicPosts(e *gin.Context) {
	TokenData, Err := utility.ExtractTokenData(e)
	if Err != nil {
		a.Logger.Error("utility.extracttokendata Error In Extract Token Data")
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if TokenData["IsAdmin"] == false {
		a.Logger.Info("IsAdmin In TokenData Has False")
		e.AbortWithStatus(http.StatusForbidden)
		return
	}

	Limit, Err := strconv.Atoi(e.Query("limit"))
	if Err != nil {
		a.Logger.Info("Error In Convert Limit To Int", zap.Error(Err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Data, err := a.Store.GetNotPublicPosts(Limit)
	if err != nil {
		a.Logger.Error("store.getnotpublicposts failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(200, model.SuccessResponse(Data))
}

func (a *App) GetReports(e *gin.Context) {
	TokenData, Err := utility.ExtractTokenData(e)
	if Err != nil {
		a.Logger.Error("utility.extracttokendata Error In Extract Token Data")
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if TokenData["IsAdmin"] == false {
		a.Logger.Info("IsAdmin In TokenData Has False")
		e.AbortWithStatus(http.StatusForbidden)
		return
	}

	Limit, Err := strconv.Atoi(e.Query("limit"))
	if Err != nil {
		a.Logger.Info("Error In Convert Limit To Int", zap.Error(Err))
		e.AbortWithStatus(http.StatusBadRequest)
		return
	}

	Data, err := a.Store.GetReports(Limit)
	if err != nil {
		a.Logger.Error("store.getreports failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(200, model.SuccessResponse(Data))
}

func (a *App) GetLikes(e *gin.Context) {
	Data, err := a.Store.GetLikes(e.Query("id"))
	if err != nil {
		a.Logger.Error("store.getlikes failed", zap.Error(err))
		e.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	e.JSON(200, model.SuccessResponse(Data))
}

func (a *App) SetRouting(e *gin.Engine) {
	Account := e.Group("/account/v1")

	Account.POST("/Login/", a.LoginUser)
	Account.POST("/Register/", a.CreateUser)

	Api := e.Group("/api/v1")
	Api.Use(middleware.JwtAuthMiddleware())

	Api.POST("/CreatePost/", a.CrNew)
	Api.POST("/GetPostsByTag/", a.GetNewsByTag)
	Api.POST("/VerifyPost/", a.VerifyNew)
	Api.POST("/LikePost/", a.LikePost)
	Api.POST("/ReportPost/", a.ReportPost)
	Api.DELETE("/DeletePost/", a.DeletePost)
	Api.GET("/GetPostForAdmin/", a.GetNotPublicPosts)
	Api.GET("/GetReports/", a.GetReports)
	Api.GET("/GetAllPost/", a.GetNews)
	Api.GET("/GetLikes/", a.GetLikes)
}
