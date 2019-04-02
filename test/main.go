package main

import "time"

func main() {

	//var quit chan bool
	//
	//quit := make(bool, chan)

	var quit chan bool
	quit = make(chan bool)

	go func() {
		time.Sleep(2 * time.Second)
		//quit<- true
		close(quit)

	}()

	for {
		select {
		case <-quit:
			println("quit")
			return
		default:
			//
		}
	}




}