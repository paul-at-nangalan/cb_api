package competition

import (
	"cb_api/cloudbet"
	"net/http"
)

const(
	GET_SPORTS = "/v2/odds/sports"
)

/*
Notes on documentation:
Probably needs a bit of explanation about how to use protobuf - or a golang example
 */

type Processor struct {
	client *http.Client
	apikey string
}

func NewProcessor(numeventprocessors int, apikey string){

}

func GetSports()[]cloudbet.Sports{

}