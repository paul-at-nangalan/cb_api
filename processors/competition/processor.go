package competition

import (
	"cb_api/cloudbet"
	"cb_api/errorhandlers"
	"cb_api/processors"
	"net/http"
	"path"
)

const(
	GET_SPORTS = "/v2/odds/sports"
)

/*
Notes on documentation:
Probably needs a bit of explanation about how to use protobuf - or a golang example
and also the Content-Type header to request for protobuf data;
 */

type Processor struct {
	processors.Retriever
}

func NewProcessor(numeventprocessors int, apikey string)*Processor{
	proc := &Processor{
	}
	proc.Setup(apikey, http.DefaultClient)
	return proc
}

func (p *Processor)GetSports()*cloudbet.Sports{
	fullurl := GET_SPORTS
	req, err := http.NewRequest(http.MethodGet, fullurl, nil)
	errorhandlers.PanicOnError(err)

	sports := &cloudbet.Sports{}

	p.GetData(req, sports)
	return sports
}


func (p *Processor)GetSportByKey(key string)*cloudbet.Sport{
	fullurl := path.Join(GET_SPORTS, key)
	req, err := http.NewRequest(http.MethodGet, fullurl, nil)
	errorhandlers.PanicOnError(err)

	sport := &cloudbet.Sport{}

	p.GetData(req, sport)
	return sport
}

func (p *Processor)Process(){
	sports := p.GetSports()

	for _, sport := range sports.Sports{
		fullsport := p.GetSportByKey(sport.Key)
	}
}
