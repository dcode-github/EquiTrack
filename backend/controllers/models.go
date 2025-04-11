package controllers

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Investment struct {
	UserId           int     `json:"user_id"`
	Instrument       string  `json:"instrument"`
	Qty              int     `json:"qty"`
	Avg              float64 `json:"avg"`
	Price            float64 `json:"ltp"`
	TotInvestment    float64 `json:"tot_invest"`
	CurVal           float64 `json:"currVal"`
	PNL              float64 `json:"pnl"`
	NetChg           float64 `json:"netChng"`
	PercentageChange float64 `json:"dayChng"`
	Date             string  `json:"date"`
}

type StockData struct {
	Price            float64 `json:"price"`
	PercentageChange float64 `json:"per_change"`
}

type TotalInvestmentData struct {
	TotalInvestment float64 `json:"total_investment"`
	TotalCurrentVal float64 `json:"total_currVal"`
	TotalPNL        float64 `json:"total_pnl"`
	TotalPNLPercent float64 `json:"total_pnl_percent"`
}

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
