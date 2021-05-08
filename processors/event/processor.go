package event

import (
	"cb_api/cloudbet"
	"cb_api/errorhandlers"
	"cb_api/processors"
	"net/http"
	"path"
)

const(
	GET_EVENT = "/v2/odds/events/"
	API_KEY_HEADER = "X-API-Key"
)

type EventDataHandler interface {
	Put(event *cloudbet.Event)
}

type Processor struct{
	processors.Retriever
	writer     EventDataHandler
	fetchqueue chan string
}

///Fetch queue is a queue of event keys to fetch data for
func NewProcessor(statwriter EventDataHandler, apikey string)*Processor{
	proc := &Processor{
		writer: statwriter,
		fetchqueue: make(chan string, 10000),
	}
	proc.Setup(apikey, http.DefaultClient)
	return proc
}

func (p *Processor)Queue(key string){
	p.fetchqueue <- key
}

///Should be run in a seperate thread
func (p *Processor)Run(){
	for {
		select {
		case eventkey := <-p.fetchqueue:
			p.processEvent(eventkey)
		}

	}
}

func (p *Processor)processEvent(eventkey string){
	fullurl := path.Join(GET_EVENT, eventkey)
	req, err := http.NewRequest(http.MethodGet, fullurl, nil)
	errorhandlers.PanicOnError(err)

	event := cloudbet.Event{}
	p.GetData(req, &event)
	///Let the event handle decide whether it needs to log this
	// and check for alerts
	p.writer.Put(&event)
}
