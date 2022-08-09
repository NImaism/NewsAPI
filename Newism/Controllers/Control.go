package Controllers

import (
	"Newism/Model"
	"Newism/Service"
	"Newism/Utility"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"time"
)

func CrNew(e *gin.Context) {
	TokenData, Err := Utility.ExtractTokenData(e)
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Extract Token Data"))
		return
	}

	Time := time.Now()

	Post := new(Model.New)
	if err := e.Bind(Post); err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Binding Data"))
		return
	}
	Post.Public = 0
	Post.Creator = TokenData["UserName"]
	Post.CreateDate = Time.String()

	Validator := validator.New()
	if err := Validator.Struct(Post); err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Data Is Not Valid"))
		return
	}

	Srv := Service.NewNService()
	id, err := Srv.CreateNew(*Post)
	if err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Insert Data"))
		return
	}
	e.JSON(http.StatusOK, Model.SuccessResponse(id))
}

func DeletePost(e *gin.Context) {
	TokenData, Err := Utility.ExtractTokenData(e)
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Extract Token Data"))
		return
	}
	if TokenData["IsAdmin"] == false {
		e.JSON(403, Model.ErrorResponse(403, "You Are Not Admin"))
		return
	}
	Srv := Service.NewNService()
	err := Srv.DeleteNew(e.Query("id"))
	if err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Delete Doc"))
		return
	}
	e.JSON(200, Model.SuccessResponse(nil))
}

func GetNews(e *gin.Context) {
	Limit, Err := strconv.Atoi(e.Query("limit"))
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Get Limit Value"))
		return
	}

	Srv := Service.NewNService()

	Data, err := Srv.GetNews(Limit)
	if err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Get Docs"))
		return
	}
	e.JSON(200, Model.SuccessResponse(Data))
}

func LoginUser(e *gin.Context) {
	LoginModel := new(Model.LoginUserViewModel)

	if err := e.Bind(LoginModel); err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Error In Binding Data"))
		return
	}

	Validator := validator.New()
	if err := Validator.Struct(LoginModel); err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Data Is Not Valid"))
		return
	}

	Srv := Service.NewUService()
	Ok, IsAdmin, err := Srv.Login(LoginModel.UserName, LoginModel.Password)
	if err != nil {
		e.JSON(http.StatusInternalServerError, Model.ErrorResponse(err, "Error In Login"))
		return
	}
	if Ok != true {
		e.JSON(http.StatusOK, Model.ErrorResponse(nil, "Account Not Found"))
		return
	}
	Token, err := Utility.GenerateToken(LoginModel.UserName, IsAdmin)
	if err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Error In Create Token"))
		return
	}
	e.JSON(http.StatusOK, Model.SuccessResponse(Token))
}

func CreateUser(e *gin.Context) {
	User := new(Model.User)
	err := e.Bind(User)
	if err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Error In Binding Data"))
		return
	}
	User.IsAdmin = 0

	Validator := validator.New()
	if err := Validator.Struct(User); err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Data Is Not Valid"))
		return
	}

	Srv := Service.NewUService()
	err = Srv.CreateUser(*User)
	if err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Error In Create User"))
		return
	}
	e.JSON(http.StatusOK, Model.SuccessResponse("Done !"))
}

func UpdateProf(e *gin.Context) {
	UpdateModel := new(Model.UpdateModel)
	if err := e.Bind(UpdateModel); err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Binding Data"))
		return
	}
	Data, Err := Utility.ExtractTokenData(e)
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Get Username From Token"))
		return
	}

	Srv := Service.NewUService()
	if err := Srv.UpdateUser(Data["UserName"], *UpdateModel); err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Updating Doc"))
		return
	}
	e.JSON(200, Model.SuccessResponse("Success"))
}

func GetNewsByTag(e *gin.Context) {
	Data := new(Model.GetNewByTag)
	if err := e.Bind(Data); err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Binding Data"))
		return
	}

	Validator := validator.New()
	if err := Validator.Struct(Data); err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Data Is Not Valid"))
		return
	}

	Srv := Service.NewNService()
	Dts, err := Srv.GetNewsByTag(*Data)
	if err != nil {
		e.JSON(500, Model.ErrorResponse(Data, "Error In Find News"))
	}
	e.JSON(200, Model.SuccessResponse(Dts))
}

func VerfyNew(e *gin.Context) {
	TokenData, Err := Utility.ExtractTokenData(e)
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Extract Token Data"))
	}
	Srv := Service.NewNService()
	if TokenData["Admin"] == false {
		e.JSON(403, Model.ErrorResponse(403, "You Are Not Admin"))
		return
	}
	err := Srv.VerfyPost(e.Query("id"))
	if err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Verify Post"))
	}
	e.JSON(200, Model.SuccessResponse("Done !"))
}

func LikePost(e *gin.Context) {
	Data := new(Model.LikePost)
	if err := e.Bind(Data); err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Binding Data"))
		return
	}

	Time := time.Now()

	TokenData, Err := Utility.ExtractTokenData(e)
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Get Token Data"))
	}

	Data.By = TokenData["UserName"]
	Data.Time = Time.String()

	Validator := validator.New()
	if err := Validator.Struct(Data); err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Data Is Not Valid"))
		return
	}

	Srv := Service.NewNService()
	if err := Srv.LikePost(*Data); err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Like Post"))
		return
	}
	e.JSON(200, Model.SuccessResponse(nil))
}

func ReportPost(e *gin.Context) {
	Data := new(Model.ReportPost)
	if err := e.Bind(Data); err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Binding Data"))
		return
	}

	Time := time.Now()

	TokenData, Err := Utility.ExtractTokenData(e)
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Get Token Data"))
	}
	Data.By = TokenData["UserName"]
	Data.Time = Time.String()

	Validator := validator.New()
	if err := Validator.Struct(Data); err != nil {
		e.JSON(http.StatusBadRequest, Model.ErrorResponse(err, "Data Not Valid"))
		return
	}

	Srv := Service.NewNService()
	if err := Srv.ReportPost(*Data); err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Report Post"))
		return
	}
	e.JSON(200, Model.SuccessResponse(nil))
}

func GetNotPublicPosts(e *gin.Context) {
	TokenData, Err := Utility.ExtractTokenData(e)
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Extract Token Data"))
		return
	}
	if TokenData["IsAdmin"] == false {
		e.JSON(403, Model.ErrorResponse(403, "You Are Not Admin"))
		return
	}

	Srv := Service.NewNService()

	Limit, Err := strconv.Atoi(e.Query("limit"))
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Get Limit"))
		return
	}

	Data, err := Srv.ShowNotPublicPost(Limit)
	if err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Report Post From Mongodb"))
		return
	}
	e.JSON(200, Model.SuccessResponse(Data))
}

func GetReports(e *gin.Context) {
	TokenData, Err := Utility.ExtractTokenData(e)
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Extract Token Data"))
		return
	}
	if TokenData["IsAdmin"] == false {
		e.JSON(403, Model.ErrorResponse(403, "You Are Not Admin"))
		return
	}

	Srv := Service.NewNService()

	Limit, Err := strconv.Atoi(e.Query("limit"))
	if Err != nil {
		e.JSON(500, Model.ErrorResponse(Err, "Error In Get Limit"))
		return
	}

	Data, err := Srv.ShowReports(Limit)
	if err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Report Post From Mongodb"))
		return
	}
	e.JSON(200, Model.SuccessResponse(Data))
}

func GetLikes(e *gin.Context) {
	Srv := Service.NewNService()
	Data, err := Srv.GetLikes(e.Query("id"))
	if err != nil {
		e.JSON(500, Model.ErrorResponse(err, "Error In Get Likes"))
		return
	}
	e.JSON(200, Model.SuccessResponse(Data))
}
