package projectx

import (
	"fmt"
	"time"
)

func (c *Client) Login(username, apiKey string) error {
	req := LoginRequest{
		UserName: username,
		APIKey:   apiKey,
	}
	var resp LoginResponse
	if err := c.doRequest("POST", "/api/Auth/loginKey", req, &resp); err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("login failed: %s", resp.ErrorMessage)
	}
	c.Token = resp.Token
	return nil
}

func (c *Client) GetAccounts(onlyActive bool) ([]Account, error) {
	req := AccountSearchRequest{OnlyActiveAccounts: onlyActive}
	var resp AccountSearchResponse
	if err := c.doRequest("POST", "/api/account/search", req, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("account search failed: %s", resp.ErrorMessage)
	}
	return resp.Accounts, nil
}

func (c *Client) GetContracts(live bool, searchText string) ([]Contract, error) {
	req := ContractSearchRequest{Live: live, SearchText: searchText}
	var resp ContractSearchResponse
	if err := c.doRequest("POST", "/api/contract/search", req, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("contract search failed: %s", resp.ErrorMessage)
	}
	return resp.Contracts, nil
}

func (c *Client) GetContractByID(contractID string) (*Contract, error) {
	req := struct {
		ContractID string `json:"contractId"`
	}{ContractID: contractID}

	var resp ContractSingleResponse
	if err := c.doRequest("POST", "/api/contract/searchById", req, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("contract search by ID failed: %s", resp.ErrorMessage)
	}
	return &resp.Contract, nil
}

func (c *Client) PlaceOrder(order OrderRequest) (*OrderResponse, error) {
	var resp OrderResponse
	if err := c.doRequest("POST", "/api/order/place", order, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return &resp, fmt.Errorf("order failed: %s", resp.ErrorMessage)
	}
	return &resp, nil
}

func (c *Client) CancelOrder(accountId, orderId int) error {
	req := struct {
		AccountID int `json:"accountId"`
		OrderID   int `json:"orderId"`
	}{
		AccountID: accountId,
		OrderID:   orderId,
	}
	var resp struct {
		Success      bool   `json:"success"`
		ErrorCode    int    `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
	}
	if err := c.doRequest("POST", "/api/order/cancel", req, &resp); err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("order cancel failed: %s", resp.ErrorMessage)
	}
	return nil
}

func (c *Client) ModifyOrder(accountId, orderId int, size *int, limitPrice, stopPrice, trailPrice *float64) error {
	req := struct {
		AccountID  int      `json:"accountId"`
		OrderID    int      `json:"orderId"`
		Size       *int     `json:"size"`
		LimitPrice *float64 `json:"limitPrice"`
		StopPrice  *float64 `json:"stopPrice"`
		TrailPrice *float64 `json:"trailPrice"`
	}{
		AccountID:  accountId,
		OrderID:    orderId,
		Size:       size,
		LimitPrice: limitPrice,
		StopPrice:  stopPrice,
		TrailPrice: trailPrice,
	}
	var resp struct {
		Success      bool   `json:"success"`
		ErrorCode    int    `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
	}
	if err := c.doRequest("POST", "/api/order/modify", req, &resp); err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("order modify failed: %s", resp.ErrorMessage)
	}
	return nil
}

func (c *Client) GetOpenPositions(accountId int) ([]OpenPosition, error) {
	req := struct {
		AccountID int `json:"accountId"`
	}{
		AccountID: accountId,
	}
	var resp OpenPositionResponse
	if err := c.doRequest("POST", "/api/position/searchOpen", req, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("open position search failed: %s", resp.ErrorMessage)
	}
	return resp.Positions, nil
}

func (c *Client) ClosePosition(accountId int, contractId string, size int) error {
	req := struct {
		AccountID  int    `json:"accountId"`
		ContractID string `json:"contractId"`
		Size       int    `json:"size"`
	}{
		AccountID:  accountId,
		ContractID: contractId,
		Size:       size,
	}
	var resp struct {
		Success      bool   `json:"success"`
		ErrorCode    int    `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
	}
	if err := c.doRequest("POST", "/api/position/closeContract", req, &resp); err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("position close failed: %s", resp.ErrorMessage)
	}
	return nil
}

func (c *Client) PartialClosePosition(accountId int, contractId string, size int) error {
	req := struct {
		AccountID  int    `json:"accountId"`
		ContractID string `json:"contractId"`
		Size       int    `json:"size"`
	}{
		AccountID:  accountId,
		ContractID: contractId,
		Size:       size,
	}
	var resp struct {
		Success      bool   `json:"success"`
		ErrorCode    int    `json:"errorCode"`
		ErrorMessage string `json:"errorMessage"`
	}
	if err := c.doRequest("POST", "/api/position/partialCloseContract", req, &resp); err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("partial position close failed: %s", resp.ErrorMessage)
	}
	return nil
}

func (c *Client) GetHistoricalBars(req HistoryRequest) ([]HistoryBar, error) {
	var resp HistoryResponse
	if err := c.doRequest("POST", "/api/history/retrieveBars", req, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("historical data request failed: %s", resp.ErrorMessage)
	}
	return resp.Bars, nil
}

func (c *Client) SearchOrders(req OrderSearchRequest) ([]OrderInfo, error) {
	var resp OrderSearchResponse
	if err := c.doRequest("POST", "/api/order/search", req, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("order search failed: %s", resp.ErrorMessage)
	}
	return resp.Orders, nil
}

func (c *Client) SearchOpenOrders(accountId int) ([]OrderInfo, error) {
	req := struct {
		AccountID int `json:"accountId"`
	}{AccountID: accountId}
	var resp OrderSearchResponse
	if err := c.doRequest("POST", "/api/order/searchOpen", req, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("open order search failed: %s", resp.ErrorMessage)
	}
	return resp.Orders, nil
}

func (c *Client) SearchTrades(accountId int, start, end *time.Time) ([]Trade, error) {
	type request struct {
		AccountID      int        `json:"accountId"`
		StartTimestamp time.Time  `json:"startTimestamp"`
		EndTimestamp   *time.Time `json:"endTimestamp,omitempty"`
	}
	type response struct {
		Trades       []Trade `json:"trades"`
		Success      bool    `json:"success"`
		ErrorCode    int     `json:"errorCode"`
		ErrorMessage string  `json:"errorMessage"`
	}
	req := request{
		AccountID:      accountId,
		StartTimestamp: *start,
		EndTimestamp:   end,
	}
	var resp response
	if err := c.doRequest("POST", "/api/trade/search", req, &resp); err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, fmt.Errorf("trade search failed: %s", resp.ErrorMessage)
	}
	return resp.Trades, nil
}
