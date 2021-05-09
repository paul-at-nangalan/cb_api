package event

import (
	"cb_api/cloudbet"
	"cb_api/errorhandlers"
	"cb_api/processors"
	"net/http"
	url2 "net/url"
	"path"
	"sync"
)

const(
	SCHEME = "https"
	SERVER = "sports-api.cloudbet.com"
	GET_EVENT = "https://sports-api.cloudbet.com/pub/v2/odds/events/"
)


type Processor struct{
	processors.Retriever
	writer     processors.EventDataHandler
	fetchqueue chan string
	waitgroup *sync.WaitGroup
}

///Fetch queue is a queue of event keys to fetch data for
func NewProcessor(statwriter processors.EventDataHandler, apikey string,
	waitgroup *sync.WaitGroup)*Processor{
	proc := &Processor{
		writer: statwriter,
		fetchqueue: make(chan string, 10000),
		waitgroup: waitgroup,
	}
	proc.Setup(apikey, http.DefaultClient)
	return proc
}

func (p *Processor)Queue(key string){
	p.waitgroup.Add(1)
	p.fetchqueue <- key
}

///Should be run in a seperate thread
func (p *Processor)Run(){
	for {
		select {
		case eventkey := <-p.fetchqueue:
			p.processEvent(eventkey)
			////Process event has a panic handler ... so this should be safe ...
			p.waitgroup.Add(-1)
		}
	}
}

func (p *Processor)processEvent(eventkey string){
	defer errorhandlers.PanicHandler()

	url := url2.URL{}
	url.Scheme = SCHEME
	url.Host = SERVER
	url.Path = path.Join(GET_EVENT, eventkey)
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	errorhandlers.PanicOnError(err)

	event := cloudbet.Event{}
	p.GetData(req, &event)
	///Let the event handle decide whether it needs to log this
	// and check for alerts
	p.writer.Put(&event)
}
