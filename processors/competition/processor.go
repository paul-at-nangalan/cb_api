package competition

import (
	"cb_api/cloudbet"
	"cb_api/errorhandlers"
	"cb_api/eventdatahandler"
	"cb_api/processors"
	"cb_api/processors/event"
	"net/http"
	"os"
	"path"
	"time"
)

const(
	GET_SPORTS = "/v2/odds/sports"
	GET_EVENTS_BYCOMP = "/v2/odds/competitions"
)

/*
Notes on documentation:
Probably needs a bit of explanation about how to use protobuf - or a golang example
and also the Content-Type header to request for protobuf data;

Too much dependancy in the API, e.g. I just want to get a list of all events (paged) with
sport, category etc ...
I don't want to have to get sports , then competitions, then iterate and then get events
 */

type Processor struct {
	processors.Retriever
	eventprocs []*event.Processor
}

func NewProcessor(numeventprocessors int, apikey string)*Processor{
	proc := &Processor{
		eventprocs: make([]*event.Processor, numeventprocessors),
	}
	proc.Setup(apikey, http.DefaultClient)

	outfile := os.ExpandEnv("STATFILE")
	interval := os.ExpandEnv("STATINTERVAL")
	dur, err := time.ParseDuration(interval)
	errorhandlers.PanicOnError(err)
	datahandler := eventdatahandler.NewDataHandler(outfile, dur)

	//// setup event processors
	for i := 0; i < numeventprocessors; i++{
		eventproc := event.NewProcessor(datahandler, apikey)
		proc.eventprocs[i] = eventproc
		go eventproc.Run()
	}
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


func (p *Processor)GetSportByKey(key string)*cloudbet.SportWithCategory{
	fullurl := path.Join(GET_SPORTS, key)
	req, err := http.NewRequest(http.MethodGet, fullurl, nil)
	errorhandlers.PanicOnError(err)

	sport := &cloudbet.SportWithCategory{}

	p.GetData(req, sport)
	return sport
}

func (p *Processor)GetEventsByCompetition(key string)*cloudbet.Competition{
	fullurl := path.Join(GET_EVENTS_BYCOMP, key)
	req, err := http.NewRequest(http.MethodGet, fullurl, nil)
	errorhandlers.PanicOnError(err)

	events := &cloudbet.Competition{}

	p.GetData(req, events)

	return events
}

func (p *Processor)Process(){
	sports := p.GetSports()

	for _, sport := range sports.Sports{
		fullsport := p.GetSportByKey(sport.Key)
		for _, cat := range fullsport.Categories{
			for _, comp := range cat.Competitions{
				events := p.GetEventsByCompetition(comp.Key)
				fetcher := 0 ///does not give perfect distribution, but should be ok
				for ev := range events.Events{
					p.eventprocs[fetcher % len(p.eventprocs)].Queue(ev.Key)
					fetcher++
				}
			}
		}
	}
}
