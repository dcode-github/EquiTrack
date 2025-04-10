package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type StockUpdate struct {
	Instrument       string  `json:"instrument"`
	Price            float64 `json:"price"`
	PercentageChange float64 `json:"per_change"`
}

type Subscription struct {
	Conn        *websocket.Conn
	Instruments map[string]bool
	Mutex       sync.Mutex
}

type Dispatcher struct {
	clients    map[*Subscription]bool
	register   chan *Subscription
	unregister chan *Subscription
	mu         sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var dispatcher = NewDispatcher()

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		clients:    make(map[*Subscription]bool),
		register:   make(chan *Subscription),
		unregister: make(chan *Subscription),
	}
	go d.run()
	return d
}

func (d *Dispatcher) run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case sub := <-d.register:
			d.mu.Lock()
			d.clients[sub] = true
			d.mu.Unlock()

		case sub := <-d.unregister:
			d.mu.Lock()
			if _, ok := d.clients[sub]; ok {
				delete(d.clients, sub)
				sub.Conn.Close()
			}
			d.mu.Unlock()

		case <-ticker.C:
			d.broadcastUpdates()
		}
	}
}

func (d *Dispatcher) broadcastUpdates() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for sub := range d.clients {
		for instrument := range sub.Instruments {
			go func(sub *Subscription, instrument string) {
				stock, err := fetchLivePrice(instrument)
				if err != nil {
					log.Println("Error fetching stock for", instrument)
					return
				}
				update := StockUpdate{
					Instrument:       instrument,
					Price:            stock.Price,
					PercentageChange: stock.PercentageChange,
				}
				data, _ := json.Marshal(update)

				sub.Mutex.Lock()
				defer sub.Mutex.Unlock()
				if err := sub.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Println("WebSocket write error:", err)
					d.unregister <- sub
				}
			}(sub, instrument)
		}
	}
}

func LivePriceWebSocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		instruments := r.URL.Query().Get("instrument")
		if instruments == "" {
			sendErrorMessage(conn, "instrument query param is required, e.g. ?symbols=TCS,INFY")
			conn.Close()
			return
		}

		instrumentList := strings.Split(instruments, ",")
		sub := &Subscription{
			Conn:        conn,
			Instruments: make(map[string]bool),
		}
		for _, symbol := range instrumentList {
			sub.Instruments[strings.ToUpper(symbol)] = true
		}

		dispatcher.register <- sub

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
		dispatcher.unregister <- sub
	}
}

func sendErrorMessage(conn *websocket.Conn, msg string) {
	errMsg := map[string]string{"error": msg}
	data, _ := json.Marshal(errMsg)
	conn.WriteMessage(websocket.TextMessage, data)
}
