package projectx

import "time"

type LoginRequest struct {
	UserName string `json:"userName"`
	APIKey   string `json:"apiKey"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	Success      bool   `json:"success"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type AccountSearchRequest struct {
	OnlyActiveAccounts bool `json:"onlyActiveAccounts"`
}

type Account struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CanTrade  bool   `json:"canTrade"`
	IsVisible bool   `json:"isVisible"`
}

type AccountSearchResponse struct {
	Accounts     []Account `json:"accounts"`
	Success      bool      `json:"success"`
	ErrorCode    int       `json:"errorCode"`
	ErrorMessage string    `json:"errorMessage"`
}

type ContractSearchRequest struct {
	Live       bool   `json:"live"`
	SearchText string `json:"searchText"`
}

type Contract struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	TickSize       float64 `json:"tickSize"`
	TickValue      float64 `json:"tickValue"`
	ActiveContract bool    `json:"activeContract"`
}

type ContractSearchResponse struct {
	Contracts    []Contract `json:"contracts"`
	Success      bool       `json:"success"`
	ErrorCode    int        `json:"errorCode"`
	ErrorMessage string     `json:"errorMessage"`
}

type ContractSingleResponse struct {
	Contract     Contract `json:"contract"`
	Success      bool     `json:"success"`
	ErrorCode    int      `json:"errorCode"`
	ErrorMessage string   `json:"errorMessage"`
}

type OrderRequest struct {
	AccountID     int      `json:"accountId"`
	ContractID    string   `json:"contractId"`
	Type          int      `json:"type"`
	Side          int      `json:"side"`
	Size          int      `json:"size"`
	LimitPrice    *float64 `json:"limitPrice"`
	StopPrice     *float64 `json:"stopPrice"`
	TrailPrice    *float64 `json:"trailPrice"`
	CustomTag     *string  `json:"customTag"`
	LinkedOrderID *int     `json:"linkedOrderId"`
}

type OrderResponse struct {
	OrderID      int    `json:"orderId"`
	Success      bool   `json:"success"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type OrderInfo struct {
	ID                int        `json:"id"`
	AccountID         int        `json:"accountId"`
	ContractID        string     `json:"contractId"`
	CreationTimestamp time.Time  `json:"creationTimestamp"`
	UpdateTimestamp   *time.Time `json:"updateTimestamp,omitempty"`
	Status            int        `json:"status"`
	Type              int        `json:"type"`
	Side              int        `json:"side"`
	Size              int        `json:"size"`
	LimitPrice        *float64   `json:"limitPrice,omitempty"`
	StopPrice         *float64   `json:"stopPrice,omitempty"`
}

type OrderSearchRequest struct {
	AccountID      int        `json:"accountId"`
	StartTimestamp time.Time  `json:"startTimestamp"`
	EndTimestamp   *time.Time `json:"endTimestamp,omitempty"`
}

type OrderSearchResponse struct {
	Orders       []OrderInfo `json:"orders"`
	Success      bool        `json:"success"`
	ErrorCode    int         `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
}

type HistoryRequest struct {
	ContractID        string    `json:"contractId"`
	Live              bool      `json:"live"`
	StartTime         time.Time `json:"startTime"`
	EndTime           time.Time `json:"endTime"`
	Unit              int       `json:"unit"`
	UnitNumber        int       `json:"unitNumber"`
	Limit             int       `json:"limit"`
	IncludePartialBar bool      `json:"includePartialBar"`
}

type HistoryBar struct {
	Time  time.Time `json:"t"`
	Open  float64   `json:"o"`
	High  float64   `json:"h"`
	Low   float64   `json:"l"`
	Close float64   `json:"c"`
	Vol   int       `json:"v"`
}

type HistoryResponse struct {
	Bars         []HistoryBar `json:"bars"`
	Success      bool         `json:"success"`
	ErrorCode    int          `json:"errorCode"`
	ErrorMessage string       `json:"errorMessage"`
}

type Trade struct {
	ID                int       `json:"id"`
	AccountID         int       `json:"accountId"`
	ContractID        string    `json:"contractId"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
	Price             float64   `json:"price"`
	ProfitAndLoss     *float64  `json:"profitAndLoss"`
	Fees              float64   `json:"fees"`
	Side              int       `json:"side"`
	Size              int       `json:"size"`
	Voided            bool      `json:"voided"`
	OrderID           int       `json:"orderId"`
}

type OpenPosition struct {
	ID                int     `json:"id"`
	AccountID         int     `json:"accountId"`
	ContractID        string  `json:"contractId"`
	CreationTimestamp string  `json:"creationTimestamp"`
	Type              int     `json:"type"`
	Size              int     `json:"size"`
	AveragePrice      float64 `json:"averagePrice"`
}

type OpenPositionResponse struct {
	Positions    []OpenPosition `json:"positions"`
	Success      bool           `json:"success"`
	ErrorCode    int            `json:"errorCode"`
	ErrorMessage string         `json:"errorMessage"`
}

// Time unit constants
const (
	TimeUnitSecond = 1
	TimeUnitMinute = 2
	TimeUnitHour   = 3
	TimeUnitDay    = 4
	TimeUnitWeek   = 5
	TimeUnitMonth  = 6
)

var TimeUnitName = map[int]string{
	TimeUnitSecond: "Second",
	TimeUnitMinute: "Minute",
	TimeUnitHour:   "Hour",
	TimeUnitDay:    "Day",
	TimeUnitWeek:   "Week",
	TimeUnitMonth:  "Month",
}
