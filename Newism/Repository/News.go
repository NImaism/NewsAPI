package Repository

import (
	database "Newism/Database"
	"Newism/Model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewRepository interface {
	GetNews(Limit int) ([]Model.New, error)
	GetNewsByTag(Data Model.GetNewByTag) ([]Model.New, error)
	CreateNew(User Model.New) (interface{}, error)
	DeleteNew(Title string) error
	Verfy(id string) error
	LikePost(Data Model.LikePost) error
	ShowPostForAdmins(Limit int) ([]Model.New, error)
	ReportPost(Data Model.ReportPost) error
	GetReports(Limit int) ([]Model.ReportPost, error)
	GetLikes(Id string) ([]Model.LikePost, error)
}

type newRepository struct{}

func NRepository() NewRepository { return newRepository{} }

func (newRepository) GetNews(Limit int) ([]Model.New, error) {
	newsCollection := database.GetCl(database.Data, "news")

	Ctx := context.TODO()

	opt := options.Find()
	opt.SetLimit(int64(Limit))

	cursor, err := newsCollection.Find(Ctx, bson.D{{"Public", 1}}, opt)
	if err != nil {
		return nil, err
	}

	var Result []Model.New
	err = cursor.All(Ctx, &Result)

	return Result, nil
}

func (newRepository) CreateNew(New Model.New) (interface{}, error) {
	newsCollection := database.GetCl(database.Data, "news")

	result, err := newsCollection.InsertOne(context.TODO(), New)
	if err != nil {
		return 0, err
	}

	return result.InsertedID, nil
}

func (newRepository) DeleteNew(Id string) error {
	newsCollection := database.GetCl(database.Data, "news")

	Object, Err := primitive.ObjectIDFromHex(Id)
	if Err != nil {
		return Err
	}

	_, err := newsCollection.DeleteOne(context.TODO(), bson.D{{"_id", Object}})
	if err != nil {
		return err
	}

	return nil
}

func (newRepository) GetNewsByTag(Data Model.GetNewByTag) ([]Model.New, error) {
	NewsCol := database.GetCl(database.Data, "news")

	Ctx := context.TODO()

	opt := options.Find()
	opt.SetLimit(int64(Data.Limit))

	Cr, Err := NewsCol.Find(Ctx, bson.D{{"Tag", Data.Tag}, {"Public", 1}}, opt)
	if Err != nil {
		return nil, Err
	}

	var Result []Model.New

	if err := Cr.All(Ctx, &Result); err != nil {
		return nil, err
	}

	return Result, nil
}

func (newRepository) Verfy(Id string) error {
	NewCol := database.GetCl(database.Data, "news")

	Object, Err := primitive.ObjectIDFromHex(Id)
	if Err != nil {
		return Err
	}

	_, Err = NewCol.UpdateOne(context.TODO(), bson.D{{"_id", Object}}, bson.D{{"$set", bson.D{{"Public", 1}}}})
	return Err
}

func (newRepository) LikePost(Data Model.LikePost) error {
	LikeCol := database.GetCl(database.Data, "likes")

	Ctx := context.TODO()

	var Result Model.LikePost

	_ = LikeCol.FindOne(Ctx, bson.D{{"PostId", Data.PostId}, {"By", Data.By}}).Decode(&Result)

	if Result.By == Data.By {
		_, err := LikeCol.DeleteOne(Ctx, bson.D{{"PostId", Data.PostId}, {"By", Data.By}})
		if err != nil {
			return err
		}
	} else {
		_, err := LikeCol.InsertOne(Ctx, Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (newRepository) ShowPostForAdmins(Limit int) ([]Model.New, error) {
	NewsCol := database.GetCl(database.Data, "news")

	var Results []Model.New

	Ctx := context.TODO()

	Opt := options.Find()
	Opt.SetLimit(int64(Limit))

	Cr, Err := NewsCol.Find(Ctx, bson.D{{"Public", 0}}, Opt)
	if Err != nil {
		return nil, Err
	}

	Err = Cr.All(Ctx, &Results)
	if Err != nil {
		return nil, Err
	}

	return Results, nil
}

func (newRepository) ReportPost(Data Model.ReportPost) error {
	ReportsCol := database.GetCl(database.Data, "reports")
	_, Err := ReportsCol.InsertOne(context.TODO(), Data)

	return Err
}

func (newRepository) GetReports(Limit int) ([]Model.ReportPost, error) {
	ReportsCol := database.GetCl(database.Data, "reports")

	Ctx := context.TODO()

	Opt := options.Find()
	Opt.SetLimit(int64(Limit))

	Cr, Err := ReportsCol.Find(Ctx, bson.D{}, Opt)
	if Err != nil {
		return nil, Err
	}

	var Results []Model.ReportPost

	Err = Cr.All(Ctx, &Results)
	if Err != nil {
		return nil, Err
	}

	return Results, nil
}

func (newRepository) GetLikes(Id string) ([]Model.LikePost, error) {
	LikeCol := database.GetCl(database.Data, "likes")

	Ctx := context.TODO()

	Cr, Err := LikeCol.Find(Ctx, bson.D{{"PostId", Id}})
	if Err != nil {
		return nil, Err
	}

	var Result []Model.LikePost

	if err := Cr.All(Ctx, &Result); err != nil {
		return nil, err
	}

	return Result, nil
}
