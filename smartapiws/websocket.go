package smartapiws

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/suyotech/smartapi-go/instruments"
)

const (
	DefaultEndpoint          = "wss://smartapisocket.angelone.in/smart-stream"
	DefaultHeartbeatInterval = 30 // in seconds
	SOCKET_MODE_LTP          = 1
	SOCKET_MODE_LTP_OHLC     = 2
	SOCKET_MODE_FULL         = 3
)

type SmartAPIWS struct {
	clientid          string
	apikey            string
	jwttoken          string
	feedtoken         string
	socket            *websocket.Conn
	heartbeatInterval int // in seconds
	retryCount        int // default 300
	retryInterval     int // in seconds
	heartbeatChannel  chan struct{}
	onDataCallback    func(TickData)
}

func NewWSClient(clientid, apikey, jwttoken, feedtoken string) (*SmartAPIWS, error) {

	return &SmartAPIWS{
		clientid:          clientid,
		apikey:            apikey,
		jwttoken:          jwttoken,
		feedtoken:         feedtoken,
		heartbeatInterval: DefaultHeartbeatInterval,
		retryCount:        300,
		retryInterval:     5,
		socket:            nil,
		heartbeatChannel:  make(chan struct{}),
	}, nil
}

func (s *SmartAPIWS) Connect() error {

	endpoint := url.URL{Scheme: "wss", Host: "smartapisocket.angelone.in", Path: "/smart-stream"}
	header := http.Header{}
	header.Set("Authorization", s.jwttoken)
	header.Set("x-api-key", s.apikey)
	header.Set("x-client-code", s.clientid)
	header.Set("x-feed-token", s.feedtoken)

	socket, resp, err := websocket.DefaultDialer.Dial(endpoint.String(), header)
	if err != nil {
		return err
	}

	errorMessage := resp.Header.Get("x-error-message")
	if errorMessage != "" {
		return fmt.Errorf("websocket connection error: %s", errorMessage)
	}

	s.socket = socket

	// reset old heartbeat goroutine
	if s.heartbeatChannel != nil {
		close(s.heartbeatChannel)
	}
	s.heartbeatChannel = make(chan struct{})

	log.Println("Socket Connected Successfully")

	s.SetHeartBeat()
	s.ReadLoop()
	s.SetOnData(s.onDataCallback)

	return nil
}

func (s *SmartAPIWS) Close() error {
	if s.socket != nil {
		close(s.heartbeatChannel) // stop ping goroutine
		err := s.socket.Close()
		s.socket = nil
		return err
	}
	return nil
}

func (s *SmartAPIWS) SetHeartBeat() {
	if s.socket == nil {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Duration(s.heartbeatInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if s.socket != nil {
					if err := s.socket.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
						log.Println("Ping error:", err)
					}
					log.Println("Heartbeat ping sent")
				}
			case <-s.heartbeatChannel:
				log.Println("Heartbeat stopped")
				return
			}
		}
	}()
}

func (s *SmartAPIWS) SetOnData(handler func(TickData)) {
	s.onDataCallback = handler
}

func (s *SmartAPIWS) ReadLoop() {
	go func() {
		for {
			mt, msg, err := s.socket.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				go s.Reconnect()
				return
			}

			switch mt {
			case websocket.PongMessage:
				log.Println("Pong received", string(msg))
			case websocket.TextMessage:
				log.Println("Text:", string(msg))
			case websocket.BinaryMessage:
				tick, err := s.parseBinaryMessage(msg)
				if err != nil {
					log.Println("Error parsing binary message:", err)
				} else {
					if s.onDataCallback != nil {
						s.onDataCallback(tick)
					}
				}
			default:
				log.Println("Unknown msg type:", mt)
			}
		}
	}()
}

func (s *SmartAPIWS) Reconnect() {
	for i := 0; i < s.retryCount; i++ {
		log.Println("Reconnecting...", i+1)
		if err := s.Connect(); err == nil {
			log.Println("Reconnected!")
			return
		}
		// Incremental backoff: min(s.retryInterval + i, 30)
		delay := s.retryInterval + i
		if delay > 30 {
			delay = 30
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func (s *SmartAPIWS) SendMessage(messageType int, data []byte) error {
	if s.socket == nil {
		return fmt.Errorf("socket is not connected")
	}
	return s.socket.WriteMessage(messageType, data)
}

func (s *SmartAPIWS) Subscribe(instruments []instruments.Instrument) error {

	cm := buildSubscribeContract(instruments)
	jsonData, err := json.Marshal(cm)
	if err != nil {
		return err
	}
	return s.SendMessage(websocket.TextMessage, jsonData)
}

type TokenListItem struct {
	ExchangeType int      `json:"exchangeType"`
	Tokens       []string `json:"tokens"`
}

type Params struct {
	Mode      int             `json:"mode"`
	TokenList []TokenListItem `json:"tokenList"`
}

type ContractMessage struct {
	CorrelationID string `json:"correlationID"`
	Action        int    `json:"action"`
	Params        Params `json:"params"`
}

func buildSubscribeContract(instruments []instruments.Instrument) ContractMessage {

	exchangeMap := make(map[string][]string)

	for _, inst := range instruments {
		exchangeMap[inst.Exchange] = append(exchangeMap[inst.Exchange], inst.Token)
	}

	tokenList := []TokenListItem{}

	for exchange, tokens := range exchangeMap {
		exchType := 0
		switch exchange {
		case "NSE":
			exchType = 1
		case "NFO":
			exchType = 2
		case "BSE":
			exchType = 3
		case "BFO":
			exchType = 4
		case "MCX":
			exchType = 5
		case "NCDEX":
			exchType = 7
		case "CDS":
			exchType = 13
		}
		tokenList = append(tokenList, TokenListItem{
			ExchangeType: exchType,
			Tokens:       tokens,
		})
	}

	contractMsg := ContractMessage{
		CorrelationID: "abcd1234",
		Action:        1,
		Params: Params{
			Mode:      3,
			TokenList: tokenList,
		},
	}
	return contractMsg
}

type TickData struct {
	SubscriptionMode  int8             // 1 byte
	ExchangeType      int8             // 1 byte
	Token             string           // 25 bytes (null-terminated string)
	SequenceNumber    int64            // 8 bytes
	ExchangeTimestamp int64            // 8 bytes
	LastTradedPrice   float64          // 8 bytes
	LastTradedQty     int64            // 8 bytes
	AvgTradedPrice    float64          // 8 bytes
	VolumeTraded      int64            // 8 bytes
	TotalBuyQty       float64          // 8 bytes
	TotalSellQty      float64          // 8 bytes
	OpenPrice         float64          // 8 bytes
	HighPrice         float64          // 8 bytes
	LowPrice          float64          // 8 bytes
	ClosePrice        float64          // 8 bytes
	LastTradedTS      int64            // 8 bytes
	OpenInterest      int64            // 8 bytes
	OIChangePercent   float64          // 8 bytes
	BestFive          [10]BestFiveData // 200 bytes (10 x 20 bytes)
	UpperCircuitLimit float64          // 8 bytes
	LowerCircuitLimit float64          // 8 bytes
	Week52High        float64          // 8 bytes
	Week52Low         float64          // 8 bytes
}

type BestFiveData struct {
	BuySellFlag int16   // 2 bytes
	Quantity    int64   // 8 bytes
	Price       float64 // 8 bytes
	NumOrders   int16   // 2 bytes
}

func scaledToFloat64(value int64, scale float64) float64 {
	return float64(value) / scale
}

func (s *SmartAPIWS) parseBinaryMessage(msg []byte) (tick TickData, err error) {
	reader := bytes.NewReader(msg)

	// 1. Subscription Mode (int8, 1 byte)
	if err = binary.Read(reader, binary.LittleEndian, &tick.SubscriptionMode); err != nil {
		return
	}
	// 2. Exchange Type (int8, 1 byte)
	if err = binary.Read(reader, binary.LittleEndian, &tick.ExchangeType); err != nil {
		return
	}
	// 3. Token (25 bytes, null-terminated string)
	tokenBytes := make([]byte, 25)
	if _, err = reader.Read(tokenBytes); err != nil {
		return
	}
	tick.Token = string(bytes.Trim(tokenBytes, "\x00"))
	// 4. Sequence Number (int64, 8 bytes)
	if err = binary.Read(reader, binary.LittleEndian, &tick.SequenceNumber); err != nil {
		return
	}
	// 5. Exchange Timestamp (int64, 8 bytes)
	if err = binary.Read(reader, binary.LittleEndian, &tick.ExchangeTimestamp); err != nil {
		return
	}
	// 6. Last Traded Price (int32, 4 bytes)
	var ltp int64
	if err = binary.Read(reader, binary.LittleEndian, &ltp); err != nil {
		return
	}
	tick.LastTradedPrice = scaledToFloat64(ltp, 100.0)

	if tick.SubscriptionMode == 1 {
		return tick, nil
	}

	// 7. Last traded quantity (int64, 8 bytes)
	if err = binary.Read(reader, binary.LittleEndian, &tick.LastTradedQty); err != nil {
		return
	}
	// 8. Average traded price (int64, 8 bytes)
	var atp int64
	if err = binary.Read(reader, binary.LittleEndian, &atp); err != nil {
		return
	}
	tick.AvgTradedPrice = scaledToFloat64(atp, 100.0)

	// 9. Volume traded for the day (int64, 8 bytes)
	if err = binary.Read(reader, binary.LittleEndian, &tick.VolumeTraded); err != nil {
		return
	}
	// 10. Total buy quantity (float64, 8 bytes)
	if err = binary.Read(reader, binary.LittleEndian, &tick.TotalBuyQty); err != nil {
		return
	}
	// 11. Total sell quantity (float64, 8 bytes)
	if err = binary.Read(reader, binary.LittleEndian, &tick.TotalSellQty); err != nil {
		return
	}
	// 12. Open price of the day (int64, 8 bytes)
	var openPrice int64
	if err = binary.Read(reader, binary.LittleEndian, &openPrice); err != nil {
		return
	}
	tick.OpenPrice = scaledToFloat64(openPrice, 100.0)
	// 13. High price of the day (int64, 8 bytes)
	var highPrice int64
	if err = binary.Read(reader, binary.LittleEndian, &highPrice); err != nil {
		return
	}
	tick.HighPrice = scaledToFloat64(highPrice, 100.0)
	// 14. Low price of the day (int64, 8 bytes)
	var lowPrice int64
	if err = binary.Read(reader, binary.LittleEndian, &lowPrice); err != nil {
		return
	}
	tick.LowPrice = scaledToFloat64(lowPrice, 100.0)
	// 15. Close price (int64, 8 bytes)
	var closePrice int64
	if err = binary.Read(reader, binary.LittleEndian, &closePrice); err != nil {
		return
	}
	tick.ClosePrice = scaledToFloat64(closePrice, 100.0)

	if tick.SubscriptionMode == 2 {
		return tick, nil
	}

	// 16. Last traded timestamp (int64, 8 bytes)
	if err = binary.Read(reader, binary.LittleEndian, &tick.LastTradedTS); err != nil {
		return
	}
	// 17. Open Interest (int64, 8 bytes)
	if err = binary.Read(reader, binary.LittleEndian, &tick.OpenInterest); err != nil {
		return
	}
	// 18. OI Change Percent (float64, 8 bytes)
	if err = binary.Read(reader, binary.LittleEndian, &tick.OIChangePercent); err != nil {
		return
	}
	// 19. Best Five Data (10 packets x 20 bytes)
	for i := 0; i < 10; i++ {
		var bestFive BestFiveData
		if err = binary.Read(reader, binary.LittleEndian, &bestFive.BuySellFlag); err != nil {
			return
		}
		if err = binary.Read(reader, binary.LittleEndian, &bestFive.Quantity); err != nil {
			return
		}
		var price int64
		if err = binary.Read(reader, binary.LittleEndian, &price); err != nil {
			return
		}
		bestFive.Price = scaledToFloat64(price, 100.0)
		if err = binary.Read(reader, binary.LittleEndian, &bestFive.NumOrders); err != nil {
			return
		}
		tick.BestFive[i] = bestFive
	}
	// 20. Upper circuit limit (int64, 8 bytes)
	var upperCircuitLimit int64
	if err = binary.Read(reader, binary.LittleEndian, &upperCircuitLimit); err != nil {
		return
	}
	tick.UpperCircuitLimit = scaledToFloat64(upperCircuitLimit, 100.0)
	// 21. Lower circuit limit (int64, 8 bytes)
	var lowerCircuitLimit int64
	if err = binary.Read(reader, binary.LittleEndian, &lowerCircuitLimit); err != nil {
		return
	}
	tick.LowerCircuitLimit = scaledToFloat64(lowerCircuitLimit, 100.0)
	// 22. Week52High (int64, 8 bytes)
	var week52High int64
	if err = binary.Read(reader, binary.LittleEndian, &week52High); err != nil {
		return
	}
	tick.Week52High = scaledToFloat64(week52High, 100.0)
	// 23. Week52Low (int64, 8 bytes)
	var week52Low int64
	if err = binary.Read(reader, binary.LittleEndian, &week52Low); err != nil {
		return
	}
	tick.Week52Low = scaledToFloat64(week52Low, 100.0)

	return tick, nil
}
