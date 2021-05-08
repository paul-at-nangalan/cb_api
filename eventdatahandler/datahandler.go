package eventdatahandler

import "cb_api/cloudbet"

type DataHandler struct{

}

func NewDataHandler()*DataHandler{
	return &DataHandler{}
}

func (d *DataHandler) Put(event *cloudbet.Event) {
	panic("implement me")
}



