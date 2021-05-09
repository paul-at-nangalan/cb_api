package processors

import "cb_api/cloudbet"

type EventDataHandler interface {
	Put(event *cloudbet.Event)
	Check() /// put this here so we can call it from the main processing thread (rather than
			//  having another thread on the data handler
}
