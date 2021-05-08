package eventdatahandler

import (
	"cb_api/cloudbet"
	"cb_api/errorhandlers"
	"encoding/csv"
	"os"
	"sync"
	"time"
)

type EventHolder struct{
	ev cloudbet.Event
	lastpubtime time.Time
	lastreporttime time.Time
}

type DataHandler struct{
	events map[string]*EventHolder ///if num events became excessively large, we'd need to fragment this
									/// across services
	csvwriter *csv.Writer
	logintrvl time.Duration
	writelock sync.Locker
}

func NewDataHandler(outputfile string, logintrvl time.Duration)*DataHandler{
	f, err := os.Create(outputfile)
	errorhandlers.PanicOnError(err)
	csvwriter := csv.NewWriter(f)
	return &DataHandler{
		events: make(map[string]*EventHolder),
		csvwriter: csvwriter,
	}
}

func (p *DataHandler)writeEventData(event *cloudbet.Event, logtype string){
	record := []string{
		logtype, event.Name, event.Key, event.CutoffTime, event.Status.String(),
	}
	p.writelock.Lock()
	defer p.writelock.Unlock()
	p.csvwriter.Write(record)
	p.events[event.Key].lastpubtime = time.Now()
}

func (p *DataHandler) Put(event *cloudbet.Event) {
	ev, ok := p.events[event.Key]
	if ok {
		if ev.lastpubtime.After(time.Now().Add(p.logintrvl)){
			p.writeEventData(event, "normal")
			ev.lastreporttime = time.Now()
		}
	}else{
		p.events[event.Key] = &EventHolder{
			ev: *event,
			lastreporttime: time.Now(),
			lastpubtime: time.Now(),
		}
		p.writeEventData(event, "normal")
	}
}

func (p *DataHandler)CheckMissing(){
	for _, ev := range p.events{
		if ev.lastreporttime.Add(time.Second).Before(time.Now()){
			p.writeEventData(&ev.ev, "alert")
		}
	}
}



