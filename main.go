package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"sync"
	"time"
)

type Message map[string]interface{}

func connectAndListen(url string, timeout time.Duration, messages chan<- Message, wg *sync.WaitGroup, print string, connectedChan chan<- bool) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Println("Error dialing websocket", err.Error())
		connectedChan <- false
	} else {
		connectedChan <- true
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	for { // loop through messages
		message := Message{}
		err = wsjson.Read(ctx, c, &message)
		if err != nil {
			break
		}
		fmt.Print(print)
		//messages <- message
	}

	c.Close(websocket.StatusNormalClosure, "")
	//log.Println("Disconnected gracefully")
	if messages != nil {
		close(messages)
	}
}

func main() {
	var duration string
	var url string
	var qty int
	var printEvery bool
	flag.BoolVar(&printEvery, "print", false, "Use --print to print a dot for each message received on each connection. False to only print for one channel.")
	flag.IntVar(&qty, "qty", 100, "Specify quantity of concurrent connections.")
	flag.StringVar(&duration, "duration", "10m", "Specify duration of test. Each connection will stay connected for this duration.")
	flag.StringVar(&url, "url", "wss://companies.stream/events", "Specify the url of the WebSocket. Should begin with ws:// or wss://.")
	flag.Parse()
	var listenFor, err = time.ParseDuration(duration)
	if err != nil {
		log.Fatal("Cannot parse duration:", duration, err)
	}

	log.Printf("Connecting %d clients to %s for %s", qty, url, listenFor.String())

	var wg sync.WaitGroup
	wg.Add(qty)
	messageChannels := make([]chan Message, qty) // messages are currently ignored
	connectedClients := 0
	for i := 0; i < qty; i++ {
		var printChar string = ""
		if i == 0 || printEvery {
			printChar = "."
		}
		connected := make(chan bool)
		go connectAndListen(url, listenFor, messageChannels[i], &wg, printChar, connected)
		success := <-connected
		if success {
			connectedClients++
			if printEvery {
				log.Printf("Client %d: Successfully connected to WebSocket", i)
			}
		} else {
			log.Printf("Client %d: Failed to connect to WebSocket", i)
		}
	}
	log.Printf("%d clients connected", connectedClients)
	wg.Wait()
	fmt.Println()
	log.Println("Finished receiving events")
}
