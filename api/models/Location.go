package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	gormbulk "github.com/bombsimon/gorm-bulk"
	"github.com/jinzhu/gorm"
)

// Location ...
type Location struct {
	ID   string `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	Type string `json:"type"`

	Slug         string `json:"slug"`
	NameWithType string `json:"name_with_type"`

	Path         string `json:"path,omitempty"`
	PathWithType string `json:"path_with_type,omitempty"`

	Code       string `json:"code" gorm:"primary_key"`
	ParentCode string `json:"parent_code"`
}

// LocationList ...
type LocationList []Location

// GernerateCoordinate ...
func (coordinate LocationList) GernerateCoordinate() LocationList {
	// Get Directory Have Location's Files
	rootDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path := rootDir + "/api/locationData/"

	// Bind Data From Json to Struct
	var mapCountry LocationList
	var data LocationList

	for number := 0; number < 14; number++ {
		fileName := fmt.Sprintf("%smap.division%s%s", path, strconv.Itoa(number), ".json")

		//Open json File
		jsonFile, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &mapCountry)
		data = append(data, mapCountry...)
	}

	coordinate = append(coordinate, data...)
	return coordinate
}

//ConvertInterface ...
func (coordinate LocationList) convertInterface(input []Location) []interface{} {
	var is = make([]interface{}, len(input))

	for i := range input {
		is[i] = input[i]
	}

	return is
}

// InstallLocaion ...
func (coordinate *LocationList) InstallLocaion(db *gorm.DB, dataInput []Location) error {
	dataAsInterface := coordinate.convertInterface(dataInput)

	coordinateSlice1 := dataAsInterface[0 : len(dataAsInterface)/2]
	coordinateSlice2 := dataAsInterface[len(dataAsInterface)/2:]

	if err := gormbulk.BulkInsert(dbv1, coordinateSlice1); err != nil {
		return err
	}

	if err := gormbulk.BulkInsert(dbv1, coordinateSlice2); err != nil {
		return err
	}
	return nil
}
