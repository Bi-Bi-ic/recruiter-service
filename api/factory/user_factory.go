package factory

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/models"
	uuid "github.com/satori/go.uuid"
)

// Userable ...
type Userable struct {
	ID           uuid.UUID `json:"id"`
	UserName     string    `json:"username"`
	Fullname     string    `json:"fullname"`
	Email        string    `json:"email"`
	MailContact  string    `json:"mail_contact"`
	Token        string    `json:"token" sql:"-"`
	RefreshToken string    `json:"refresh_token" sql:"-"`
	Avatar       string    `json:"avatar"`
	Cover        string    `json:"cover"`
	GuestMode    bool      `json:"guestmode"`
	TimeLines    []TimeLineAble  `json:"time_line"`
}

// TimeLineAble ...
type TimeLineAble struct {
	ID       	uuid.UUID `json:"id"`
	Title    	string    `json:"title"`
	SubTitle 	string    `json:"sub_title"`
	FromTime 	int64  `json:"from_time"`
	ToTime   	int64  `json:"to_time"`
	Description string    `json:"description"`	
}

type TimeLineFactory struct{}

func (factory TimeLineFactory) Create(timeLineEntity models.TimeLine) TimeLineAble {
	return TimeLineAble{
		ID:       	timeLineEntity.ID,
		SubTitle: 	timeLineEntity.SubTitle,
		Title:    	timeLineEntity.Title,
		FromTime: 	timeLineEntity.FromTime,
		ToTime:   	timeLineEntity.ToTime,
		Description: timeLineEntity.Description,
	}
}

// UserInfoFactory ... this object create for anything what if you want about user
type UserInfoFactory struct{}

// URL_SERVER  ...
const URL_SERVER string = "https://api.huc.com.vn"

// Create is a list of 'Userable' fixed from User entity
func (factory UserInfoFactory) Create(user interface{}) Userable {
	customer := models.User{}

	err := mapstructure.Decode(user, &customer)
	if err != nil {
		panic(err)
	}

	timeLinesAble := []TimeLineAble{}

	timeLineFactory := TimeLineFactory{}

	for _, value := range customer.TimeLines {
		
		timeLinesAble = append(timeLinesAble, timeLineFactory.Create(value))
	}

	return Userable{
		ID:           customer.ID,
		UserName:     customer.UserName,
		Fullname:     customer.Fullname,
		Email:        customer.Email,
		Token:        customer.Token,
		RefreshToken: customer.RefreshToken,
		Avatar:       URL_SERVER + customer.Avatar,
		Cover:        URL_SERVER + customer.Cover,
		GuestMode:    customer.GuestMode,
		TimeLines: 	timeLinesAble,
	}
}

// CreateDetail ...
func (factory UserInfoFactory) CreateDetail(user interface{}) Userable {
	return factory.Create(user)
}
