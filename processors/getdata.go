package processors

import (
	"cb_api/errorhandlers"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
)

const (
	API_KEY_HEADER = "X-API-Key"
)

type Retriever struct{
	apikey string
	client *http.Client
}

func (p *Retriever)Setup(apikey string, client *http.Client){
	p.client = client
	p.apikey = apikey
}

func (p *Retriever)GetData(req *http.Request, out proto.Message){

	req.Header.Set(API_KEY_HEADER, p.apikey)
	req.Header.Set("Accept-Encoding", "application/x-protobuf")

	resp, err := p.client.Do(req)
	errorhandlers.PanicOnError(err)

	data, err := ioutil.ReadAll(resp.Body)
	errorhandlers.PanicOnError(err)
	err = proto.Unmarshal(data, out)
	errorhandlers.PanicOnError(err)
}
