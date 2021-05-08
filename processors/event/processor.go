package event

import (
	"cb_api/cloudbet"
	"cb_api/errorhandlers"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"path"
)

const(
	GET_EVENT = "/v2/odds/events/"
	API_KEY_HEADER = "X-API-Key"
)

type StatWriter interface {
	Write(event *cloudbet.Event)
	AlertEventRemoved(event *cloudbet.Event)
}

type Processor struct{
	writer StatWriter
	fetchqueue chan string
	client *http.Client
	apikey string
}

///Fetch queue is a queue of event keys to fetch data for
func NewProcessor(statwriter StatWriter, apikey string, fetchqueue chan string)*Processor{
	return &Processor{
		writer: statwriter,
		fetchqueue: fetchqueue,
		client: http.DefaultClient,
	}
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

	req.Header.Set(API_KEY_HEADER, p.apikey)
	req.Header.Set("Accept-Ancoding", "application/x-protobuf")

	resp, err := p.client.Do(req)
	errorhandlers.PanicOnError(err)

	data, err := ioutil.ReadAll(resp.Body)
	errorhandlers.PanicOnError(err)
	event := cloudbet.Event{}
	err = proto.Unmarshal(data, &event)
	errorhandlers.PanicOnError(err)

}
