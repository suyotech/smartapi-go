package smartapiws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	DefaultEndpoint          = "wss://smartapisocket.angelone.in/smart-stream"
	DefaultHeartbeatInterval = 30 // in seconds
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

	header := http.Header{}
	header.Set("Authorization", s.jwttoken)
	header.Set("x-api-key", s.apikey)
	header.Set("x-client-code", s.clientid)
	header.Set("x-feed-token", s.feedtoken)

	socket, _, err := websocket.DefaultDialer.Dial(DefaultEndpoint, header)
	if err != nil {
		return err
	}

	s.socket = socket
	s.SetHeartBeat()

	return nil
}

func (s *SmartAPIWS) Close() error {
	if s.socket != nil {
		return s.socket.Close()
	}
	return nil
}

func (s *SmartAPIWS) SetHeartBeat() {

	go func() {
		ticker := time.NewTicker(time.Duration(s.heartbeatInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if s.socket != nil {
					err := s.socket.WriteMessage(websocket.PingMessage, []byte{})
					if err != nil {
						// Handle error (e.g., log it, attempt reconnection, etc.)
					}
				}
			case <-s.heartbeatChannel:
				return
			}
		}
	}()
}
