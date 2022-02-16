package infoconf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	SkipValue = "~"
)

// defaults
var Defaults = map[string]interface{}{
	"DEnableInfos":         false,
	"DInfos":               []ConfiguredInfo{},
	"DEnableTickers":       false,
	"DTickers":             []ConfiguredTicker{},
	"DHideWatermarkTicker": false,
}

/*
var Defaults = map[string]interface{}{
	"DID":        60001070, // TODO: check for significance
	"DPriority":  1,
	"DStartTime": -2,
	"DEndTime":   -3,
	"DData": InfoData{
		"everyDay",
		"This is a default message. If you're seeing this, you probably didn't configure correctly!",
		"10000000",
		"text",
		SkipValue,
	},
}
*/

var DisplayTypes = map[string]string{
	"everyDay":     "0",
	"once":         "1",
	"fullTime":     "2",
	"onlyInfoPage": "3",
}

var InfoTypes = map[string]string{
	"text":          "0",  // text is specified in "message" field
	"image":         "1",  // show image with specified imageID
	"feed":          "2",  // share button
	"shopCancel":    "3",  // ?
	"feedUrl":       "4",  // same as 9, but with a share button
	"feedRoulette":  "5",  // same as 10, but with a share button
	"feedShop":      "6",  // same as 11, but with a share button (?)
	"feedEvent":     "7",  // same as 12, but with a share button
	"feedEventList": "8",  // same as 13, but with a share button
	"url":           "9",  // external link is provided in the "extra" field
	"roulette":      "10", // button for roulette
	"shop":          "11", // button for shop
	"event":         "12", // for events with a top page
	"eventList":     "13", // for events with a rewards list
	"rouletteInfo":  "14", // the roulette banner
	"quickInfo":     "15", // a banner across the timed mode button (and/or story mode button, depending on event type and data)
	"countryText":   "16", // based on region code in "extra" field
	"countryImage":  "17", // based on region code in "extra" field
}

type InfoData struct {
	DisplayType string `json:"displayType,omitempty"` // if left out of param, game uses
	Message     string `json:"message,omitempty"`     // if left out of param, game uses empty string
	ImageID     string `json:"imageID,omitempty"`     // if left out of param game uses "-1". if used with text, this is the image that appears in the info list (ad_<id>_<regionCode>.jpg)
	InfoType    string `json:"infoType,omitempty"`    // if left out of param, game uses Text
	Extra       string `json:"extra,omitempty"`       // desired region code if DisplayType == countryText or countryImage; web address otherwise
}

type ConfiguredTicker struct {
	Message   string `json:"message"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
}

type ConfiguredInfo struct {
	ID        int64    `json:"id"`
	Priority  int64    `json:"priority"`
	StartTime int64    `json:"startTime"`
	EndTime   int64    `json:"endTime"`
	Data      InfoData `json:"content"`
}

/*
* StartTime and EndTime can use certain negative values to represent time
functions that are not otherwise easily possible using a static
JSON file. -2 will be the beginning of the system's day, and -3 will
be the end of the system's day.
A value of -4 will activate the info no matter what. For StartTime, it
will use the current time minus 1 in order to say the info has just started.
For EndTime, it will use 24 hours into the future. This isn't foolproof, as
the info could expire and yield odd results, but if you have a player that
is playing for long enough for that to happen, you have a _different_ problem
on your hands.
*/

func (c ConfiguredInfo) ConstructParam() string {
	displayType := DisplayTypes[c.Data.DisplayType]
	message := c.Data.Message
	imageID := c.Data.ImageID
	infoType := InfoTypes[c.Data.InfoType]
	extra := c.Data.Extra
	param := ""
	if displayType == SkipValue {
		return param
	}
	param += displayType
	if message == SkipValue {
		return param
	}
	param += "_" + message
	if imageID == SkipValue {
		return param
	}
	param += "_" + imageID
	if infoType == SkipValue {
		return param
	}
	param += "_" + infoType
	if extra == SkipValue {
		return param
	}
	param += "_" + extra
	return param
}

func (c ConfiguredInfo) HasValidDisplayType() bool {
	_, ok := DisplayTypes[c.Data.DisplayType]
	return ok
}
func (c ConfiguredInfo) HasValidInfoType() bool {
	_, ok := InfoTypes[c.Data.InfoType]
	return ok
}

var CFile ConfigFile

type ConfigFile struct {
	EnableInfos         bool               `json:"enableInformation"`
	Infos               []ConfiguredInfo   `json:"infos"`
	EnableTickers       bool               `json:"enableTickers"`
	Tickers             []ConfiguredTicker `json:"tickers"`
	HideWatermarkTicker bool               `json:"hideWatermarkTicker"`
}

func Parse(filename string) error {
	CFile = ConfigFile{
		Defaults["DEnableInfos"].(bool),
		Defaults["DInfos"].([]ConfiguredInfo),
		Defaults["DEnableTickers"].(bool),
		Defaults["DTickers"].([]ConfiguredTicker),
		Defaults["DHideWatermarkTicker"].(bool),
	}
	file, err := loadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &CFile)
	if err != nil {
		return err
	}
	newInfos := []ConfiguredInfo{}
	for i, ci := range CFile.Infos {
		if !ci.HasValidDisplayType() {
			log.Printf("[WARN] Invalid information display type %s at index %v, ignoring\n", ci.Data.DisplayType, i)
			continue
		}
		if !ci.HasValidInfoType() {
			log.Printf("[WARN] Invalid information info type %s at index %v, ignoring\n", ci.Data.InfoType, i)
			continue
		}
		newInfos = append(newInfos, ci)
	}
	CFile.Infos = newInfos
	return nil
}

func loadFile(filename string) ([]byte, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return b, err
}
