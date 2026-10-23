# smartapi-go

Unofficial Angel One SmartAPI SDK for internal use, written in Go.

This SDK provides REST clients for authentication, orders, portfolio data,
historical candles, and option Greeks; a WebSocket client for market ticks; and
helpers for downloading and searching the instrument master.

## Installation

Requires **Go 1.25.3 or newer**, as declared in `go.mod`.

```sh
go get github.com/suyotech/smartapi-go
```

| Package | Responsibility |
| --- | --- |
| `github.com/suyotech/smartapi-go/smartapi` | Sessions, REST requests, orders, portfolio, and historical data |
| `github.com/suyotech/smartapi-go/smartapiws` | WebSocket connections, subscriptions, and decoded ticks |
| `github.com/suyotech/smartapi-go/instruments` | Instrument master download, cache, filtering, expiry dates, and strikes |

## Authentication quick start

Provide your client ID, MPIN, SmartAPI API key, and TOTP **secret key**. The SDK
generates the current TOTP code from the secret when `GenerateSession` runs.
The environment variable names below are application conventions; the SDK does
not read environment variables itself.

```go
package main

import (
	"log"
	"os"

	"github.com/suyotech/smartapi-go/smartapi"
)

func main() {
	client := smartapi.NewClient(
		os.Getenv("SMARTAPI_CLIENT_ID"),
		os.Getenv("SMARTAPI_MPIN"),
		os.Getenv("SMARTAPI_API_KEY"),
		os.Getenv("SMARTAPI_TOTP_SECRET"),
	)
	if err := client.GenerateSession(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Logout(); err != nil {
			log.Printf("logout: %v", err)
		}
	}()

	profile, err := client.UserProfile()
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Signed in as %s (%s)", profile.Name, profile.ClientCode)
}
```

`NewClient` collects local IP, public IP, and MAC metadata. Public IP discovery
makes outbound requests to IP lookup services during construction. Local IP
discovery calls `log.Fatal` if its UDP dial fails, which exits the process.

`GenerateSession` sends the credentials and generated TOTP to the login endpoint,
then stores `JWTToken`, `RefreshToken`, and `FeedToken` on the client. Subsequent
REST calls attach the JWT as a bearer token. To restore an existing session:

```go
client.SetTokens(jwtToken, refreshToken, feedToken)
```

There is no automatic token refresh or refresh-token exchange method. The
application must obtain replacement tokens when needed. `Logout` sends a logout
request but does not clear the token fields; use `client.SetTokens("", "", "")`
to clear them locally. Configure clients before use and coordinate token changes
with active requests; REST client fields are not protected by a mutex.

## REST API reference

All methods below belong to `*smartapi.Client`. Returned values are decoded from
the API response's `data` field; the outer response wrapper is not returned.

| Method | Return values | Purpose |
| --- | --- | --- |
| `GenerateSession()` | `error` | Generate TOTP, log in, and store tokens |
| `SetTokens(jwtToken, refreshToken, feedToken string)` | None | Assign existing tokens |
| `UserProfile()` | `*UserProfileResponse, error` | Read account profile |
| `UserMargin()` | `*UserMarginResponse, error` | Read funds and margin |
| `Logout()` | `error` | Request session logout |
| `PlaceOrder(req *PlaceOrderRequest)` | `*PlaceOrderResponse, error` | Submit an order |
| `ModifyOrder(req *ModifyOrderRequest)` | `*ModifyOrderResponse, error` | Modify an order; see current limitations |
| `CancelOrder(req *CancelOrderRequest)` | `error` | Cancel one order |
| `GetOrderBook()` | `[]OrderBook, error` | Read orders |
| `GetTradeBook()` | `[]TradeBook, error` | Read executed trades |
| `CancelAllOrders()` | `error` | Cancel orders with status `open` or `trigger pending` |
| `GetHoldings()` | `[]Holding, error` | Read holdings |
| `GetAllTotalHoldings()` | `AllTotalHolding, error` | Read holdings and aggregate totals |
| `GetPositions()` | `[]Positions, error` | Read positions |
| `ConvertPosition(req ConvertPositionRequest)` | `error` | Request a product conversion |
| `ClosePosition(position Positions)` | `error` | Submit an opposite-side market order; see current limitations |
| `CloseAllPositions()` | `error` | Fetch positions and call `ClosePosition` sequentially |
| `GetCandleData(req *CandleDataRequest)` | `[]Candle, error` | Read historical OHLC, volume, and optional OI |
| `GetOptionGreeks(name, expirydate string)` | `[]OptionGreek, error` | Read option Greeks for a name and expiry |
| `SetProxy(proxyURL string)` | `error` | Configure this REST client's proxy |
| `SetDebug(debug bool)` | None | Enable or disable verbose request/response logging |

Request and response fields are defined in [user.go](smartapi/user.go),
[orders.go](smartapi/orders.go), [portfolio.go](smartapi/portfolio.go), and
[data.go](smartapi/data.go). Some broker numeric fields remain strings, such as
margin values and several position fields; use the declared Go field types.

### Orders

The following snippets assume `client` is authenticated and run inside a function.
Order calls submit real broker requests. Resolve the current symbol and token
from the instrument master before submitting.

```go
order, err := client.PlaceOrder(&smartapi.PlaceOrderRequest{
	Variety:         smartapi.VARIETY_NORMAL,
	TradingSymbol:   instrument.Symbol,
	SymbolToken:     instrument.Token,
	Exchange:        instrument.Exchange,
	TransactionType: smartapi.TRANSACTION_TYPE_BUY,
	OrderType:       smartapi.ORDER_TYPE_LIMIT,
	Quantity:        quantity,
	ProductType:     smartapi.PRODUCT_TYPE_INTRADAY,
	Price:           limitPrice,
	Duration:        smartapi.DURATION_DAY,
})
if err != nil {
	return err
}

// An order ID acknowledges submission; inspect the order/trade book for fills.
fmt.Println(order.OrderID, order.UniqueOrderID)

// Cancel only when your application intends to cancel this order.
err = client.CancelOrder(&smartapi.CancelOrderRequest{
	Variety: smartapi.VARIETY_NORMAL,
	OrderID: order.OrderID,
})
if err != nil {
	return err
}
```

`PlaceOrderRequest` also accepts `TriggerPrice`, `SquareOff`, `StopLoss`,
`TrailingStopLoss`, `DisclosedQuantity`, and `OrderTag`. These optional fields are
omitted from JSON when zero or empty. The SDK does not validate lot sizes, tick
sizes, required field combinations, or broker eligibility before submission.

Order constants in [constants.go](smartapi/constants.go) include:

| Category | Constants (prefix followed by suffixes) |
| --- | --- |
| Variety | `VARIETY_`: `NORMAL`, `STOPLOSS`, `ROBO` |
| Side | `TRANSACTION_TYPE_`: `BUY`, `SELL` |
| Order type | `ORDER_TYPE_`: `MARKET`, `LIMIT`, `STOPLOSS_LIMIT`, `STOPLOSS_MARKET` |
| Product | `PRODUCT_TYPE_`: `INTRADAY`, `DELIVERY`, `CARRYFORWARD`, `MARGIN`, `BO` |
| Duration | `DURATION_`: `DAY`, `IOC` |
| REST exchange | `EXCHANGE_`: `NSE`, `BSE`, `CDS`, `MCX`, `NFO` |

### Portfolio

```go
positions, err := client.GetPositions()
if err != nil {
	return err
}
for _, position := range positions {
	fmt.Println(position.TradingSymbol, position.NetQty, position.NetValue)
}

margin, err := client.UserMargin()
if err != nil {
	return err
}
fmt.Println("Available cash:", margin.AvailableCash)
```

`ConvertPosition` takes a value rather than a pointer. Populate its
`ConvertPositionRequest` with the instrument metadata, old and new product types,
transaction side, quantity, and conversion type required by the broker.

### Historical candles and option Greeks

```go
candles, err := client.GetCandleData(&smartapi.CandleDataRequest{
	Exchange:    smartapi.EXCHANGE_NSE,
	SymbolToken: instrument.Token,
	Interval:    smartapi.INTERVAL_5MIN,
	FromDate:    "2026-10-22 09:15",
	ToDate:      "2026-10-22 15:30",
})
if err != nil {
	return err
}
for _, candle := range candles {
	fmt.Println(candle.Time, candle.Open, candle.High, candle.Low, candle.Close)
}

greeks, err := client.GetOptionGreeks(underlyingName, expiryDate)
if err != nil {
	return err
}
for _, greek := range greeks {
	fmt.Println(greek.StrikePrice, greek.Optiontype, greek.Delta, greek.IV)
}
```

Use request dates in `YYYY-MM-DD HH:MM` format. Candle timestamps are parsed into
`time.Time` using the response's numeric timezone offset, with `Asia/Kolkata`
loaded as the location. Missing volume or OI columns become zero. Option expiry
strings are passed to the broker unchanged; use the relevant instrument expiry.

Available interval constants are `INTERVAL_1MIN`, `INTERVAL_3MIN`, `INTERVAL_5MIN`,
`INTERVAL_10MIN`, `INTERVAL_15MIN`, `INTERVAL_30MIN`, `INTERVAL_1HOUR`, and
`INTERVAL_1DAY`.

## Instrument master

The instrument package works independently of REST authentication. This example
uses a custom cache path, refreshes it if needed, and finds an equity instrument:

```go
if err := instruments.CheckInstrumentsUpdateAt("data/instruments.json"); err != nil {
	return err
}
master, err := instruments.LoadInstrumentsAt("data/instruments.json")
if err != nil {
	return err
}
matches, err := instruments.FindInstruments(instruments.FindInstrumentRequest{
	Exchange: "NSE",
	Symbol:   "SBIN-EQ",
}, master)
if err != nil {
	return err
}
if len(matches) == 0 {
	return fmt.Errorf("instrument not found")
}
instrument := matches[0]
fmt.Println(instrument.Symbol, instrument.Token, instrument.LotSize)
```

Import `github.com/suyotech/smartapi-go/instruments` for these helpers.

| Function | Behavior |
| --- | --- |
| `DownloadInstruments() error` | Download and replace `instruments.json` in the working directory |
| `DownloadInstrumentsAt(filePath string) error` | Download to a chosen path, creating parent directories |
| `CheckInstrumentsUpdate() error` | Check the default file and download if needed |
| `CheckInstrumentsUpdateAt(filePath string) error` | Check a chosen cache file |
| `LoadInstruments() ([]Instrument, error)` | Load the default file |
| `LoadInstrumentsAt(filePath string) ([]Instrument, error)` | Load a chosen cache file |
| `FindInstruments(req FindInstrumentRequest, master []Instrument) ([]Instrument, error)` | Filter instruments by exact, case-sensitive field matches |
| `GetExpiryDates(req FindInstrumentRequest, master []Instrument) ([]string, error)` | Return unique, parseable expiry dates in ascending order |
| `GetOptionStrikes(price float64, req FindInstrumentRequest, maxStrikes int, master []Instrument) (OptionStrikes, error)` | Return nearest strike and up to `maxStrikes` on either side |

Refresh checks use the latest **08:30 IST** cutoff: today's cutoff after 08:30,
otherwise yesterday's. Missing, invalid, empty, or older files trigger a download.
This uses file modification time and does not account for trading holidays.
Downloads have a two-minute HTTP timeout, validate a nonempty JSON instrument
list, and replace the destination through a temporary file in the same directory.

Search fields are `Token`, `Symbol`, `Name`, `Expiry`, `Strike`, `InstrumentType`,
and `Exchange`. Empty request fields are ignored, and all populated filters must
match. Consequently, `INTTYPE_EQ` (an empty string) does not restrict a search to
equities by itself. When expiry values are present, matches are sorted by expiry.
`GetExpiryDates` skips blank or unparseable expiry values.

Passing a nil master slice to a search helper loads the default
`instruments.json`; it does not download it. For a custom path, load the file and
pass the resulting slice explicitly.

Instrument numeric fields, including strike, lot size, and tick size, are strings
from the master. `GetOptionStrikes` compares the supplied price against master
strikes divided by 100, then returns `NEAR`, `UP`, and `DOWN` as strings in the
original master scale. `UP` is ascending and `DOWN` is descending. Ensure the
filter contains at least one numeric strike before calling it; an empty strike
set currently panics.

Instrument type constants include `INTTYPE_EQ`, `INTTYPE_INDEX`, `INTTYPE_FUTIDX`,
`INTTYPE_FUTSTX`, `INTTYPE_OPTIDX`, `INTTYPE_OPTSTX`, `INTTYPE_FUTCOM`, and
`INTTYPE_OPTFUT`.

## WebSocket market data

After authenticating the REST client, pass its JWT and feed token to the
WebSocket constructor. Register the callback before connecting:

```go
// Imports in addition to smartapi: os, os/signal, instruments, and smartapiws.
// Import paths for SDK packages are listed in the installation table above.
ws, err := smartapiws.NewWSClient(
	client.ClientId, client.APIKey, client.JWTToken, client.FeedToken,
)
if err != nil {
	return err
}
defer ws.Close()

ws.SetOnData(func(tick smartapiws.TickData) {
	fmt.Println(tick.Token, tick.LastTradedPrice, tick.VolumeTraded)
})
if err := ws.Connect(); err != nil {
	return err
}
if err := ws.Subscribe([]instruments.Instrument{instrument}); err != nil {
	return err
}

// Connect starts background goroutines and returns; keep the application alive.
interrupt := make(chan os.Signal, 1)
signal.Notify(interrupt, os.Interrupt)
defer signal.Stop(interrupt)
<-interrupt
```

`Connect` opens the market-data socket and starts the heartbeat and read loops.
The heartbeat interval is 30 seconds. `Subscribe` groups instrument tokens by
exchange and always requests **full mode (3)**, using correlation ID `abcd1234`.
It maps `NSE`, `NFO`, `BSE`, `BFO`, `MCX`, `NCDEX`, and `CDS` to exchange types;
unknown names are sent with exchange type zero, so supply a supported name.

Binary messages are decoded to `TickData`: token, exchange, sequence, timestamps,
LTP, quantity, average price, volume, buy/sell quantities, OHLC, OI, ten best-five
entries, circuit limits, and 52-week prices, depending on subscription mode.
The parser divides integer price fields by 100 for every exchange; timestamps
are exposed as raw `int64` values. Mode constants are `SOCKET_MODE_LTP` (1),
`SOCKET_MODE_LTP_OHLC` (2), and `SOCKET_MODE_FULL` (3), but `Subscribe` has no mode
argument.

Callbacks run synchronously on the read goroutine. Keep them short or hand work
off to your own queue. Text messages and parsing/connection errors are logged;
there is no separate error or connection-state callback.

A read failure starts automatic reconnect attempts, up to 300, with a delay
starting at five seconds and increasing to a maximum of 30 seconds. Reconnect
does **not** restore subscriptions or refresh authentication. The application
must arrange resubscription after reconnecting; there is no reconnect callback.
`Close` stops the heartbeat and closes the socket permanently for that instance;
create a new client to connect after closing.

`Connect` already starts `SetHeartBeat` and `ReadLoop`; do not start duplicate
loops manually. Advanced callers can use `SendMessage(messageType int, data
[]byte) error` for raw messages. There is no dedicated unsubscribe helper.
Application writes, subscription calls, and heartbeat writes are not serialized
by a shared write lock; do not assume concurrent writes are safe.

## Proxies and debug logging

REST and WebSocket proxies are configured independently, before requests or
connection establishment:

```go
if err := client.SetProxy("http://127.0.0.1:8080"); err != nil {
	return err
}
if err := ws.SetProxy("http://127.0.0.1:8080"); err != nil {
	return err
}
```

Pass an empty string to either `SetProxy` to force a direct connection. Without
an explicit override, the clients inherit the default HTTP/WebSocket transport's
environment proxy behavior. WebSocket proxy configuration accepts `http` and
`https`; HTTPS proxies use a TLS CONNECT tunnel and support credentials in the
proxy URL. A WebSocket proxy change applies to subsequent connections.

REST proxy settings do not configure instrument downloads or the public-IP
lookup performed during `NewClient`.

`client.SetDebug(true)` prints request bodies, headers, response bodies, and
decoded results to standard output. This includes credentials, TOTP codes, API
keys, and session tokens. Debug logging defaults to disabled; only enable it
where those logs can be handled securely.

## Errors and current limitations

REST requests have a 30-second timeout and return transport, JSON decoding,
HTTP, or broker API errors. Non-2xx responses return `*smartapi.HTTPError`, with
`StatusCode`, `Message`, and `RawBody`. A broker `AB1021` response also returns
this type, even when the HTTP status is successful. Other broker error codes
are returned as formatted errors.

```go
// Import "errors" for errors.Is and errors.As.
if errors.Is(err, smartapi.ErrRateLimited) {
	// Let the application schedule a later attempt.
}
var httpErr *smartapi.HTTPError
if errors.As(err, &httpErr) {
	fmt.Println(httpErr.StatusCode, httpErr.Message)
}
```

`ErrRateLimited` matches when an `HTTPError` message/body contains `AB1021`,
`too many requests`, or the recognized access-rate denial text. HTTP status 429
alone does not guarantee that match. There is no automatic REST retry or rate
limiter. Reconcile order state before retrying a submission whose result is
unknown, since a request error does not prove the broker rejected the order.

The following behaviors reflect the current implementation:

- `ModifyOrder` passes a nil result pointer internally, so its response remains
  nil even on success. Use the order book to inspect the result.
- `ClosePosition` chooses the opposite side but sets quantity to `-NetQty`.
  This produces a negative quantity for long positions. Prefer an explicit
  `PlaceOrder` request with a positive quantity until this helper is corrected.
  `CloseAllPositions` uses the same helper and stops at the first error.
- `CancelAllOrders` handles only exact statuses `open` and `trigger pending`,
  uses `VARIETY_NORMAL` for every cancellation, and stops at the first error.
- `GetHoldings` currently calls `/getAllHolding` and expects an array;
  `GetAllTotalHoldings` calls `/getAllHolding/total` and expects an object.
  These endpoint/response combinations have not been verified against the live
  broker here.
- `GetCandleData` assumes correctly shaped broker rows and uses unchecked type
  assertions. Unexpected row lengths or field types can panic.
- No paper-trading mode, fill tracking, or order-status stream is provided.
  REST submission and cancellation success do not establish final execution.

These examples document the repository implementation. They have not been
validated against a live brokerage account.
