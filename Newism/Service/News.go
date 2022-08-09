package Service

import (
	"Newism/Model"
	"Newism/Repository"
)

type NewsService interface {
	GetNews(Limit int) ([]Model.New, error)
	CreateNew(User Model.New) (interface{}, error)
	DeleteNew(Title string) error
	GetNewsByTag(Data Model.GetNewByTag) ([]Model.New, error)
	VerfyPost(Id string) error
	LikePost(Data Model.LikePost) error
	ReportPost(Data Model.ReportPost) error
	ShowNotPublicPost(Limit int) ([]Model.New, error)
	ShowReports(Limit int) ([]Model.ReportPost, error)
	GetLikes(Id string) ([]Model.LikePost, error)
}

type newsService struct {
}

func NewNService() NewsService {
	return newsService{}
}

func (newsService) GetNews(Limit int) ([]Model.New, error) {
	Srv := Repository.NRepository()
	return Srv.GetNews(Limit)
}

func (newsService) CreateNew(User Model.New) (interface{}, error) {
	Srv := Repository.NRepository()
	return Srv.CreateNew(User)
}

func (newsService) DeleteNew(Id string) error {
	Service := Repository.NRepository()
	return Service.DeleteNew(Id)
}

func (newsService) GetNewsByTag(Data Model.GetNewByTag) ([]Model.New, error) {
	Server := Repository.NRepository()
	return Server.GetNewsByTag(Data)
}

func (newsService) VerfyPost(Id string) error {
	Server := Repository.NRepository()
	return Server.Verfy(Id)
}

func (newsService) LikePost(Data Model.LikePost) error {
	Server := Repository.NRepository()
	return Server.LikePost(Data)
}

func (newsService) ReportPost(Data Model.ReportPost) error {
	Server := Repository.NRepository()
	return Server.ReportPost(Data)
}

func (newsService) ShowNotPublicPost(Limit int) ([]Model.New, error) {
	Server := Repository.NRepository()
	return Server.ShowPostForAdmins(Limit)
}

func (newsService) ShowReports(Limit int) ([]Model.ReportPost, error) {
	Server := Repository.NRepository()
	return Server.GetReports(Limit)
}

func (newsService) GetLikes(Id string) ([]Model.LikePost, error) {
	Server := Repository.NRepository()
	return Server.GetLikes(Id)
}
