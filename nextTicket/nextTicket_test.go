package nextTicket

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func init() {
	go func() {
		http.ListenAndServe("0.0.0.0:10030", nil)
	}()
}

func TestTicket(t *testing.T) {
	ticket, err := NewTicket(time.Second * 3)
	if err != nil {
		panic(err)
	}
	defer ticket.Stop()
	i := 0
	// var init bool
	for {
		fmt.Println("1")
		select {
		case now := <-ticket.C:
			fmt.Println("2")
			if i >= 100 {
				return
			}
			i += 1
			if i%2 != 0 {
				ticket.Next(time.Second)
				fmt.Println("c")
				continue
			}
			fmt.Println(now, time.Now())
			ticket.Next(0)
		}
	}
}
