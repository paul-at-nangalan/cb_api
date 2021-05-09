package competition

import (
	"cb_api/cloudbet"
	"cb_api/errorhandlers"
	"cb_api/processors"
	"cb_api/processors/event"
	"fmt"
	"net/http"
	url2 "net/url"
	"path"
	"strings"
	"sync"
	"time"
)

const(
	SCHEME = "https"
	SERVER = "sports-api.cloudbet.com"
	GET_SPORTS = "/pub/v2/odds/sports"
	GET_EVENTS_BYCOMP = "/pub/v2/odds/competitions"
)

/*
Notes on documentation:
Probably needs a bit of explanation about how to use protobuf - or a golang example
and also the Content-Type header to request for protobuf data;

Too much dependancy in the API, e.g. I just want to get a list of all events (paged) with
sport, category etc ...
I don't want to have to get sports , then competitions, then iterate and then get events

The documentation doesn't seem to mention the dns, I ran the example with inspect to see
the server address

The documentation doesn't seem to mention the header "accept": "application/x-protobuf"
 */

type Processor struct {
	processors.Retriever
	eventprocs []*event.Processor
	datahandler processors.EventDataHandler
	waitgroup *sync.WaitGroup
}

func NewProcessor(numeventprocessors int, apikey string, datahandler processors.EventDataHandler)*Processor{
	proc := &Processor{
		eventprocs: make([]*event.Processor, numeventprocessors),
	}
	proc.Setup(apikey, http.DefaultClient)

	proc.datahandler = datahandler

	waitgroup := &sync.WaitGroup{}
	proc.waitgroup = waitgroup
	//// setup event processors
	for i := 0; i < numeventprocessors; i++{
		eventproc := event.NewProcessor(datahandler, apikey, proc.waitgroup)
		proc.eventprocs[i] = eventproc
		go eventproc.Run()
	}
	return proc
}

func (p *Processor)Run(procintrvl time.Duration){
	procticker := time.NewTicker(procintrvl)
	for{
		select {
		case <-procticker.C:
			p.Process()
			fmt.Println("Waiting for get events to finish")
			p.waitgroup.Wait()
			fmt.Println("Done")
			p.datahandler.Check()
		}
	}
}

func (p *Processor)GetSports()*cloudbet.Sports{
	url := url2.URL{}
	url.Scheme = SCHEME
	url.Host = SERVER
	url.Path = GET_SPORTS
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	errorhandlers.PanicOnError(err)

	sports := &cloudbet.Sports{}

	p.GetData(req, sports)
	return sports
}


func (p *Processor)GetSportByKey(key string)*cloudbet.SportWithCategory{
	url := url2.URL{}
	url.Scheme = SCHEME
	url.Host = SERVER
	url.Path = path.Join(GET_SPORTS, key)
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	errorhandlers.PanicOnError(err)

	sport := &cloudbet.SportWithCategory{}

	p.GetData(req, sport)
	return sport
}

func (p *Processor)GetEventsByCompetition(key string)*cloudbet.Competition{
	url := url2.URL{}
	url.Scheme = SCHEME
	url.Host = SERVER
	url.Path = path.Join(GET_EVENTS_BYCOMP, key)
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
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
				for _, ev := range events.Events{
					isvalid := true
					for name, _ := range ev.Markets{
						///not sure this is correct ...
						if strings.Contains(name, ".outright"){
							isvalid = false
							break
						}
					}
					if !isvalid{
						continue
					}
					p.eventprocs[fetcher % len(p.eventprocs)].Queue(ev.Key)
					fetcher++
				}
			}
		}
	}
}

