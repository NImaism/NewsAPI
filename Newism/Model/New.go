package Model

type New struct {
	Id          string      `bson:"_id,omitempty"`
	Creator     interface{} `bson:"Creator" validate:"required"`
	Title       string      `bson:"Title" json:"Title,omitempty" validate:"required"`
	Description string      `bson:"Description,omitempty" json:"Description,omitempty" validate:"required"`
	Image       string      `bson:"Image,omitempty" json:"Image,omitempty" validate:"required"`
	Tag         string      `bson:"Tag,omitempty" json:"Tag,omitempty" validate:"required"`
	CreateDate  string      `bson:"CreateDate"`
	Public      int         `bson:"Public"`
}

type GetNewByTag struct {
	Limit int    `json:"Limit" validate:"required"`
	Tag   string `json:"Tag" validate:"required"`
}

type LikePost struct {
	Id     string      `bson:"_id,omitempty"`
	PostId string      `bson:"PostId,omitempty" json:"PostId" validate:"required"`
	By     interface{} `bson:"By,omitempty" validate:"required"`
	Time   string      `bson:"Time,omitempty" validate:"required"`
}

type ReportPost struct {
	Id      string      `bson:"_id,omitempty"`
	PostId  string      `bson:"PostId,omitempty" json:"PostId" validate:"required"`
	Title   string      `bson:"Title,omitempty" json:"Title" validate:"required"`
	Content string      `bson:"Content,omitempty" json:"Content" validate:"required"`
	By      interface{} `bson:"By,omitempty" validate:"required"`
	Time    string      `bson:"Time,omitempty" validate:"required"`
}
