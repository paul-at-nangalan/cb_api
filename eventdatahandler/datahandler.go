package eventdatahandler

import (
	"cb_api/cloudbet"
	"cb_api/errorhandlers"
	"encoding/csv"
	"log"
	"os"
	"sync"
	"time"
)

type EventHolder struct{
	ev cloudbet.Event
	lastpubtime time.Time
	lastreporttime time.Time
	cutofftime time.Time
	clear bool
}

type DataHandler struct{
	events map[string]*EventHolder ///if num events became excessively large, we'd need to fragment this
									/// across services
	csvwriter *csv.Writer
	logintrvl time.Duration
	alertintrvl time.Duration
	writelock sync.Locker
	maplock sync.Locker
	f *os.File
}

func NewDataHandler(outputfile string, logintrvl time.Duration, alertintrvl time.Duration)*DataHandler{
	f, err := os.Create(outputfile)
	errorhandlers.PanicOnError(err)
	csvwriter := csv.NewWriter(f)
	return &DataHandler{
		events: make(map[string]*EventHolder),
		csvwriter: csvwriter,
		logintrvl: logintrvl,
		alertintrvl: alertintrvl,
		f: f,
	}
}

func (p *DataHandler)Close(){
	p.f.Close()
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
		if ev.clear{
			////We have a problem ... we've refound and event or something
			log.Panicln("Refound a cleared event ", event.Key)
		}
		if ev.lastpubtime.After(time.Now().Add(p.logintrvl)){
			p.writeEventData(event, "normal")
			ev.lastreporttime = time.Now()
		}
	}else{
		cutofftime, err := time.Parse("2006-01-02T15:04:05Z07:00", event.CutoffTime)
		errorhandlers.PanicOnError(err)
		p.events[event.Key] = &EventHolder{
			ev: *event,
			lastreporttime: time.Now(),
			lastpubtime: time.Now(),
			cutofftime: cutofftime,
		}
		p.writeEventData(event, "normal")
	}
}

//// because this can delete entries from the map, it must NOT be called while any event
//// processors are still running ....
//// and locks are bad for performance
//// also, the logic is messed up if any processors are still running
func (p *DataHandler)Check(){
	for key, ev := range p.events{
		if ev.clear{
			///deleting from a map inside an iteration is apparrently safe ... according to stack overflow
			/// I'd dig a bit deeper before using it in prod
			delete(p.events, key)
			continue
		}
		if time.Now().After(ev.cutofftime){
			p.events[key].clear = true
			continue
		}
		if ev.lastreporttime.Add(p.alertintrvl).Before(time.Now()) {
			p.writeEventData(&ev.ev, "alert")
			///already reported, clear it out
			p.events[key].clear = true
		}
	}
}



