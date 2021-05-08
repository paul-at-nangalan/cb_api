package errorhandlers

import (
	"log"
	"runtime/debug"
)

func PanicOnError(err error){
	if err != nil{
		debug.PrintStack()
		log.Panicln("ERROR ", err)
	}
}
