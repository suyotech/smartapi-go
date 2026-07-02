package smartapiws

import "testing"

func TestCloseIsIdempotentAndReconnectAfterCloseReturns(t *testing.T) {
	ws, err := NewWSClient("client", "api", "jwt", "feed")
	if err != nil {
		t.Fatalf("NewWSClient returned error: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	if err := ws.Close(); err != nil {
		t.Fatalf("first Close returned error: %v", err)
	}
	if err := ws.Close(); err != nil {
		t.Fatalf("second Close returned error: %v", err)
	}

	ws.Reconnect()

	ws.mu.Lock()
	defer ws.mu.Unlock()
	if !ws.closed {
		t.Fatal("websocket should remain closed")
	}
	if ws.reconnecting {
		t.Fatal("Reconnect should not start after Close")
	}
	if ws.heartbeatChannel != nil {
		t.Fatal("heartbeatChannel should be nil after Close")
	}
}
