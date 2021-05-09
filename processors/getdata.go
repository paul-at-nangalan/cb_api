package processors

import (
	"cb_api/errorhandlers"
	"fmt"
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

func (p *Retriever)GetClient()*http.Client{
	return p.client
}

func (p *Retriever)SetHeaders(req *http.Request){
	req.Header.Set(API_KEY_HEADER, p.apikey)
	req.Header.Set("accept", "application/x-protobuf")
}

func (p *Retriever)GetData(req *http.Request, out proto.Message){

	p.SetHeaders(req)

	resp, err := p.client.Do(req)
	errorhandlers.PanicOnError(err)

	data, err := ioutil.ReadAll(resp.Body)
	errorhandlers.PanicOnError(err)
	fmt.Println("RAW DATA ", string(data))
	err = proto.Unmarshal(data, out)
	errorhandlers.PanicOnError(err)
}
