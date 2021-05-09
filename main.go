package main

import (
	"cb_api/errorhandlers"
	"cb_api/eventdatahandler"
	"cb_api/processors/competition"
	"os"
	"strconv"
	"time"
)

func getIntVar(name string)int64{
	val := os.ExpandEnv(name)
	ival, err := strconv.ParseInt(val, 10, 64)
	errorhandlers.PanicOnError(err)
	return ival
}
func getDurVar(name string)time.Duration{
	val := os.ExpandEnv(name)
	ival, err := time.ParseDuration(val)
	errorhandlers.PanicOnError(err)
	return ival
}

func main(){
	defer errorhandlers.PanicHandler()

	outfile := os.ExpandEnv("STATFILE")
	interval := os.ExpandEnv("STATINTERVAL")
	apikey := os.ExpandEnv("APIKEY")
	numeventprocs := getIntVar("NUM_EVENTPROCESSORS")
	statintrvl := getDurVar("STATSINTERVAL")
	checkintrvl := getDurVar("CHECKINTERVAL")/// check for events that have gone missing

	dur, err := time.ParseDuration(interval)
	errorhandlers.PanicOnError(err)
	datahandler := eventdatahandler.NewDataHandler(outfile, dur)

	proc := competition.NewProcessor(int(numeventprocs), apikey, datahandler)
	proc.Run(statintrvl, checkintrvl) /// run this on the main thread
}
