package projectx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/philippseith/signalr"
)

// MarketDataHandler defines the interface for handling real-time market data events.
// Implementations should process quotes, trades, and market depth updates for a given contract.
type MarketDataHandler interface {
	OnQuote(contractID string, data map[string]interface{}) // Called when a new quote is received
	OnTrade(contractID string, data map[string]interface{}) // Called when a new trade is executed
	OnDepth(contractID string, data map[string]interface{}) // Called when market depth changes
}

// SignalRClient manages the WebSocket connection to the market data hub using SignalR.
// It handles connection lifecycle, subscription management, and message routing.
type SignalRClient struct {
	client         signalr.Client     // The underlying SignalR client
	mutex          sync.RWMutex       // Protects access to shared state
	subscriptions  map[string]bool    // Tracks active contract subscriptions
	marketHandler  MarketDataHandler  // Handles market data events
	isConnected    bool               // Current connection state
	reconnectCount int                // Number of reconnection attempts
	ctx            context.Context    // Context for cancellation
	cancel         context.CancelFunc // Function to cancel the context
}

// NewSignalRClient creates a new SignalR client with the given JWT token and market data handler.
// It establishes a WebSocket connection to the market data hub and sets up message handling.
func NewSignalRClient(jwtToken string, marketHandler MarketDataHandler) (*SignalRClient, error) {
	// Create a cancellable context for the client
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the client structure
	client := &SignalRClient{
		subscriptions: make(map[string]bool),
		marketHandler: marketHandler,
		ctx:           ctx,
		cancel:        cancel,
	}

	// Configure the SignalR hub URL
	hubURL := "wss://rtc.thefuturesdesk.projectx.com/hubs/market"
	parsedURL, err := url.Parse(hubURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hub URL: %v", err)
	}

	// Add JWT token to query parameters for authentication
	q := parsedURL.Query()
	q.Add("access_token", jwtToken)
	parsedURL.RawQuery = q.Encode()

	// Create HTTP connection with WebSocket transport
	// This sets up the underlying WebSocket connection with proper headers
	conn, err := signalr.NewHTTPConnection(ctx, parsedURL.String(),
		signalr.WithTransports(signalr.TransportWebSockets),
		signalr.WithHTTPHeaders(func() http.Header {
			h := http.Header{}
			h.Set("Authorization", "Bearer "+jwtToken)
			return h
		}))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create SignalR connection: %v", err)
	}

	// Create SignalR client with the HTTP connection and register this instance as the message receiver
	c, err := signalr.NewClient(ctx,
		signalr.WithConnection(conn),
		signalr.WithReceiver(client))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create SignalR client: %v", err)
	}

	client.client = c
	return client, nil
}

// OnConnected is called when the SignalR connection is established.
// It updates the connection state and resubscribes to all previously subscribed contracts.
func (c *SignalRClient) OnConnected(connectionID string) {
	c.mutex.Lock()
	c.isConnected = true
	c.mutex.Unlock()
	log.Printf("SignalR connected with ID: %s", connectionID)

	// Resubscribe to all contracts that were previously subscribed
	for contractID := range c.subscriptions {
		if err := c.Subscribe(contractID); err != nil {
			log.Printf("Failed to resubscribe to %s: %v", contractID, err)
		}
	}
}

// OnDisconnected is called when the SignalR connection is lost.
// It updates the connection state and increments the reconnection counter.
func (c *SignalRClient) OnDisconnected(connectionID string) {
	c.mutex.Lock()
	c.isConnected = false
	c.reconnectCount++
	c.mutex.Unlock()
	log.Printf("SignalR disconnected (attempt %d)", c.reconnectCount)
}

// OnGatewayQuote handles incoming quote messages from the SignalR hub.
// It forwards the quote data to the market data handler.
func (c *SignalRClient) OnGatewayQuote(contractID string, data map[string]interface{}) {
	c.marketHandler.OnQuote(contractID, data)
}

// OnGatewayTrade handles incoming trade messages from the SignalR hub.
// It forwards the trade data to the market data handler.
func (c *SignalRClient) OnGatewayTrade(contractID string, data map[string]interface{}) {
	c.marketHandler.OnTrade(contractID, data)
}

// OnGatewayDepth handles incoming market depth messages from the SignalR hub.
// It forwards the depth data to the market data handler.
func (c *SignalRClient) OnGatewayDepth(contractID string, data map[string]interface{}) {
	c.marketHandler.OnDepth(contractID, data)
}

// Start initiates the SignalR connection.
// This begins the WebSocket connection and message processing.
func (c *SignalRClient) Start() error {
	c.client.Start()
	return nil
}

// Stop gracefully shuts down the SignalR connection.
// It unsubscribes from all contracts and closes the connection.
func (c *SignalRClient) Stop() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Unsubscribe from all contracts before stopping
	for contractID := range c.subscriptions {
		if err := c.unsubscribe(contractID); err != nil {
			log.Printf("Failed to unsubscribe from %s: %v", contractID, err)
		}
	}

	c.cancel() // Cancel the context to stop all operations
	c.isConnected = false
	c.client.Stop()
	return nil
}

// Subscribe adds a subscription for the specified contract.
// It sends subscription requests for quotes, trades, and market depth.
func (c *SignalRClient) Subscribe(contractID string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		return fmt.Errorf("not connected to SignalR hub")
	}

	// Subscribe to quotes
	ch := c.client.Send("SubscribeContractQuotes", contractID)
	if err := <-ch; err != nil {
		return fmt.Errorf("failed to subscribe to quotes: %v", err)
	}

	// Subscribe to trades
	ch = c.client.Send("SubscribeContractTrades", contractID)
	if err := <-ch; err != nil {
		return fmt.Errorf("failed to subscribe to trades: %v", err)
	}

	// Subscribe to market depth
	ch = c.client.Send("SubscribeContractMarketDepth", contractID)
	if err := <-ch; err != nil {
		return fmt.Errorf("failed to subscribe to market depth: %v", err)
	}

	c.subscriptions[contractID] = true
	return nil
}

// unsubscribe removes a subscription for the specified contract.
// It sends unsubscribe requests for quotes, trades, and market depth.
func (c *SignalRClient) unsubscribe(contractID string) error {
	if !c.isConnected {
		return fmt.Errorf("not connected to SignalR hub")
	}

	// Unsubscribe from quotes
	ch := c.client.Send("UnsubscribeContractQuotes", contractID)
	if err := <-ch; err != nil {
		return fmt.Errorf("failed to unsubscribe from quotes: %v", err)
	}

	// Unsubscribe from trades
	ch = c.client.Send("UnsubscribeContractTrades", contractID)
	if err := <-ch; err != nil {
		return fmt.Errorf("failed to unsubscribe from trades: %v", err)
	}

	// Unsubscribe from market depth
	ch = c.client.Send("UnsubscribeContractMarketDepth", contractID)
	if err := <-ch; err != nil {
		return fmt.Errorf("failed to unsubscribe from market depth: %v", err)
	}

	delete(c.subscriptions, contractID)
	return nil
}

// Unsubscribe safely removes a subscription for the specified contract.
// It acquires a lock before calling unsubscribe to ensure thread safety.
func (c *SignalRClient) Unsubscribe(contractID string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.unsubscribe(contractID)
}

// IsConnected returns the current connection state.
// It uses a read lock to safely access the connection state.
func (c *SignalRClient) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}
