package store

import (
	"Newism/internal/database"
	"Newism/internal/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mognodb struct {
	Mongo *database.Mongo
}

func New(Database *database.Mongo) *Mognodb {
	return &Mognodb{Database}
}

func (s *Mognodb) GetNews(Limit int) ([]model.New, error) {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	newsCollection := s.Mongo.GetCl("news")

	opt := options.Find()
	opt.SetLimit(int64(Limit))

	cursor, err := newsCollection.Find(Ctx, bson.D{{"Public", 1}}, opt)
	if err != nil {
		return nil, err
	}

	var Result []model.New
	err = cursor.All(Ctx, &Result)

	return Result, nil
}

func (s *Mognodb) CreateNew(New model.New) (interface{}, error) {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	newsCollection := s.Mongo.GetCl("news")

	result, err := newsCollection.InsertOne(Ctx, New)
	if err != nil {
		return 0, err
	}

	return result.InsertedID, nil
}

func (s *Mognodb) DeleteNew(Id string) error {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	newsCollection := s.Mongo.GetCl("news")

	Object, Err := primitive.ObjectIDFromHex(Id)
	if Err != nil {
		return Err
	}

	var Result model.New

	_ = newsCollection.FindOne(Ctx, bson.D{{"_id", Object}}).Decode(&Result)

	if Result.Id != Id {
		return errors.New("Post NotFound")
	}

	_, err := newsCollection.DeleteOne(Ctx, bson.D{{"_id", Object}})
	if err != nil {
		return err
	}

	return nil
}

func (s *Mognodb) GetNewsByTag(Data model.GetNewByTag) ([]model.New, error) {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	NewsCol := s.Mongo.GetCl("news")

	opt := options.Find()
	opt.SetLimit(int64(Data.Limit))

	Cr, Err := NewsCol.Find(Ctx, bson.D{{"Tag", Data.Tag}, {"Public", 1}}, opt)
	if Err != nil {
		return nil, Err
	}

	var Result []model.New

	if err := Cr.All(Ctx, &Result); err != nil {
		return nil, err
	}

	return Result, nil
}

func (s *Mognodb) VerifyPost(Id string) error {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	NewCol := s.Mongo.GetCl("news")

	Object, Err := primitive.ObjectIDFromHex(Id)
	if Err != nil {
		return Err
	}

	var Result model.New

	_ = NewCol.FindOne(Ctx, bson.D{{"_id", Object}}).Decode(&Result)

	if Result.Id != Id {
		return errors.New("Post NotFound")
	}

	_, Err = NewCol.UpdateOne(Ctx, bson.D{{"_id", Object}}, bson.D{{"$set", bson.D{{"Public", 1}}}})
	return Err
}

func (s *Mognodb) LikePost(Data model.LikePost) error {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	LikeCol := s.Mongo.GetCl("likes")

	var Result model.LikePost

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

func (s *Mognodb) GetNotPublicPosts(Limit int) ([]model.New, error) {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	NewsCol := s.Mongo.GetCl("news")

	var Results []model.New

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

func (s *Mognodb) ReportPost(Data model.ReportPost) error {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	ReportsCol := s.Mongo.GetCl("reports")
	var Result model.ReportPost

	_ = ReportsCol.FindOne(Ctx, bson.D{{"PostId", Data.PostId}, {"By", Data.By}}).Decode(&Result)

	if Result.By == Data.By {
		return errors.New("You Have One Report For Post")
	}
	_, Err := ReportsCol.InsertOne(Ctx, Data)
	return Err
}

func (s *Mognodb) GetReports(Limit int) ([]model.ReportPost, error) {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	ReportsCol := s.Mongo.GetCl("reports")

	Opt := options.Find()
	Opt.SetLimit(int64(Limit))

	Cr, Err := ReportsCol.Find(Ctx, bson.D{}, Opt)
	if Err != nil {
		return nil, Err
	}

	var Results []model.ReportPost

	Err = Cr.All(Ctx, &Results)
	if Err != nil {
		return nil, Err
	}

	return Results, nil
}

func (s *Mognodb) GetLikes(Id string) ([]model.LikePost, error) {
	Ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	LikeCol := s.Mongo.GetCl("likes")

	Cr, Err := LikeCol.Find(Ctx, bson.D{{"PostId", Id}})
	if Err != nil {
		return nil, Err
	}

	var Result []model.LikePost

	if err := Cr.All(Ctx, &Result); err != nil {
		return nil, err
	}

	return Result, nil
}

func (s *Mognodb) Login(UserName string, Password string) (bool, bool, error) {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	usersCollection := s.Mongo.GetCl("users")

	var Result model.User

	err := usersCollection.FindOne(ctx, bson.D{{"UserName", UserName}, {"Password", Password}}).Decode(&Result)
	if err != nil {
		return false, false, err
	}

	if Result.IsAdmin == 1 {
		return true, true, nil
	}

	return true, false, nil
}

func (s *Mognodb) CreateUser(user model.User) error {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	var Result model.User

	UserCol := s.Mongo.GetCl("users")

	if Result.UserName == user.UserName {
		return errors.New("one account is active with this username")
	}

	_, err := UserCol.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
