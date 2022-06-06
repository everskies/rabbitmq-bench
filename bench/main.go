package main

import (
	"fmt"
	"os"
	"time"
	"strconv"
	"runtime"

	"github.com/streadway/amqp"
)

var msg = []byte("Hello, world!")

func publish() int {
	start := time.Now()
	connection, err := amqp.Dial(os.Args[1:][0])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer connection.Close()

	channel, _ := connection.Channel()

	channel.Publish("", "test", false, false, amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            msg,
		DeliveryMode:    amqp.Transient,
		Priority:        0,
	})

	return int(time.Since(start).Milliseconds())
}

func publishLoop(c chan<- int) {
	for {
		c <- publish()
	}
}

func poll(input <-chan int) {
	cache := make([]int, 0, 10000)
	tick := time.NewTicker(1000 * time.Millisecond)

	for {
		select {
		case m := <-input:
			cache = append(cache, m)
		case <-tick.C:
			messages := len(cache)

			total := 0
			max := 0
			for i := 0; i < messages; i++ {
				total += cache[i]
				if cache[i] > max {
					max = cache[i]
				}
			}

			fmt.Printf("%d/sec avg: %dms max: %dms\n", messages, total/messages, max)
			cache = cache[:0]
		}
	}

}

func main() {
	messages := make(chan int)
	threadsArg, _ := strconv.ParseInt(os.Args[1:][1], 10, 32)
	threads := int(threadsArg)

	go poll(messages)
	for i := 0; i < threads; i++ {
		go publishLoop(messages)
	}

	runtime.Goexit()
}
