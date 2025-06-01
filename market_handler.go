package projectx

import (
	"log"
	"sync"
	"time"
)

type MarketDataCallback func(bar HistoryBar)

type MarketDataManager struct {
	mutex         sync.RWMutex
	currentBar    *HistoryBar
	lastTradeTime time.Time
	barPeriod     time.Duration
	callback      MarketDataCallback
	contractID    string
}

func NewMarketDataManager(contractID string, barPeriodMinutes int, callback MarketDataCallback) *MarketDataManager {
	return &MarketDataManager{
		barPeriod:  time.Duration(barPeriodMinutes) * time.Minute,
		callback:   callback,
		contractID: contractID,
	}
}

func (m *MarketDataManager) OnQuote(contractID string, data map[string]interface{}) {
	if contractID != m.contractID {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract quote data
	bid, ok1 := data["bid"].(float64)
	ask, ok2 := data["ask"].(float64)
	if !ok1 || !ok2 {
		log.Printf("Invalid quote data format")
		return
	}

	now := time.Now()

	// Initialize or update current bar
	if m.currentBar == nil {
		m.initializeNewBar(now, (bid+ask)/2)
		return
	}

	// Update current bar
	price := (bid + ask) / 2
	if price > m.currentBar.High {
		m.currentBar.High = price
	}
	if price < m.currentBar.Low {
		m.currentBar.Low = price
	}
	m.currentBar.Close = price

	// Check if it's time to close the bar
	if now.Sub(m.currentBar.Time) >= m.barPeriod {
		m.closeCurrentBar()
		m.initializeNewBar(now, price)
	}
}

func (m *MarketDataManager) OnTrade(contractID string, data map[string]interface{}) {
	if contractID != m.contractID {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract trade data
	price, ok1 := data["price"].(float64)
	size, ok2 := data["size"].(float64)
	if !ok1 || !ok2 {
		log.Printf("Invalid trade data format")
		return
	}

	now := time.Now()
	m.lastTradeTime = now

	// Initialize or update current bar
	if m.currentBar == nil {
		m.initializeNewBar(now, price)
		return
	}

	// Update current bar
	if price > m.currentBar.High {
		m.currentBar.High = price
	}
	if price < m.currentBar.Low {
		m.currentBar.Low = price
	}
	m.currentBar.Close = price
	m.currentBar.Vol += int(size)

	// Check if it's time to close the bar
	if now.Sub(m.currentBar.Time) >= m.barPeriod {
		m.closeCurrentBar()
		m.initializeNewBar(now, price)
	}
}

func (m *MarketDataManager) OnDepth(contractID string, data map[string]interface{}) {
	// Market depth data is not used for bar construction
}

func (m *MarketDataManager) initializeNewBar(t time.Time, price float64) {
	barStartTime := t.Truncate(m.barPeriod)
	m.currentBar = &HistoryBar{
		Time:  barStartTime,
		Open:  price,
		High:  price,
		Low:   price,
		Close: price,
		Vol:   0,
	}
}

func (m *MarketDataManager) closeCurrentBar() {
	if m.currentBar != nil && m.callback != nil {
		m.callback(*m.currentBar)
	}
}
