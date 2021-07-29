package factory

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/models"
)

// LocationAble ...
type LocationAble struct {
	Total  int                 `json:"total"`
	Result models.LocationList `json:"result"`
}

// LocationInfoFactory this object create for anything what if you want about location
type LocationInfoFactory struct{}

// Create is a list of 'LocationAble' fixed from Location entity
func (factory LocationInfoFactory) Create(location interface{}) LocationAble {
	var coordinates models.LocationList

	err := mapstructure.Decode(location, &coordinates)
	if err != nil {
		panic(err)
	}

	return LocationAble{
		Total:  len(coordinates),
		Result: coordinates,
	}
}
