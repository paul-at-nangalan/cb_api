package errorhandlers

import (
	"log"
	"runtime/debug"
)

func PanicOnError(err error){
	if err != nil{
		log.Panicln("ERROR ", err)
	}
}

func PanicHandler(){
	if r := recover(); r != nil{
		debug.PrintStack()
		log.Println("ERROR: ", r)
	}
}
