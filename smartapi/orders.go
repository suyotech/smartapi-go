package smartapi

type PlaceOrderResponse struct {
	Script        string `json:"script"`
	OrderID       string `json:"orderid"`
	UniqueOrderID string `json:"uniqueorderid"`
}

type PlaceOrderRequest struct {
	Variety           string  `json:"variety"`
	TradingSymbol     string  `json:"tradingsymbol"`
	SymbolToken       string  `json:"symboltoken"`
	Exchange          string  `json:"exchange"`
	TransactionType   string  `json:"transactiontype"`
	OrderType         string  `json:"ordertype"`
	Quantity          int     `json:"quantity"`
	ProductType       string  `json:"producttype"`
	Price             float64 `json:"price,omitempty"`
	TriggerPrice      float64 `json:"triggerprice,omitempty"`
	SquareOff         float64 `json:"squareoff,omitempty"`
	StopLoss          float64 `json:"stoploss,omitempty"`
	TrailingStopLoss  float64 `json:"trailingStopLoss,omitempty"`
	DisclosedQuantity int     `json:"disclosedquantity,omitempty"`
	Duration          string  `json:"duration"`
	OrderTag          string  `json:"ordertag,omitempty"`
}

type ModifyOrderRequest struct {
	Variety       string  `json:"variety"`
	OrderID       string  `json:"orderid"`
	OrderType     string  `json:"ordertype"`
	ProductType   string  `json:"producttype"`
	Duration      string  `json:"duration"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	TradingSymbol string  `json:"tradingsymbol"`
	SymbolToken   string  `json:"symboltoken"`
	Exchange      string  `json:"exchange"`
}

type ModifyOrderResponse struct {
	OrderID       string `json:"orderid"`
	UniqueOrderID string `json:"uniqueorderid"`
}

type CancelOrderRequest struct {
	Variety string `json:"variety"`
	OrderID string `json:"orderid"`
}

type GetOrderBookResponse struct {
	Variety                 string  `json:"variety"`
	OrderType               string  `json:"ordertype"`
	ProductType             string  `json:"producttype"`
	Duration                string  `json:"duration"`
	Price                   float64 `json:"price"`
	TriggerPrice            float64 `json:"triggerprice,string"`
	Quantity                int64   `json:"quantity,string"`
	DisclosedQuantity       int64   `json:"disclosedquantity,string"`
	SquareOff               float64 `json:"squareoff,string"`
	StopLoss                float64 `json:"stoploss,string"`
	TrailingStopLoss        float64 `json:"trailingstoploss,string"`
	TradingSymbol           string  `json:"tradingsymbol"`
	TransactionType         string  `json:"transactiontype"`
	Exchange                string  `json:"exchange"`
	SymbolToken             string  `json:"symboltoken"`
	InstrumentType          string  `json:"instrumenttype"`
	StrikePrice             float64 `json:"strikeprice,string"`
	OptionType              string  `json:"optiontype"`
	ExpiryDate              string  `json:"expirydate"`
	LotSize                 int64   `json:"lotsize,string"`
	CancelSize              int64   `json:"cancelsize,string"`
	AveragePrice            float64 `json:"averageprice,string"`
	FilledShares            int64   `json:"filledshares,string"`
	UnfilledShares          int64   `json:"unfilledshares,string"`
	OrderID                 string  `json:"orderid"`
	Text                    string  `json:"text"`
	Status                  string  `json:"status"`
	OrderStatus             string  `json:"orderstatus"`
	UpdateTime              string  `json:"updatetime"`
	ExchangeTime            string  `json:"exchtime"`
	ExchangeOrderUpdateTime string  `json:"exchangeorderupdatetime"`
	FillID                  string  `json:"fillid"`
	FillTime                string  `json:"filltime"`
	ParentOrderID           string  `json:"parentorderid"`
	UniqueOrderID           string  `json:"uniqueorderid"`
	ExchangeOrderID         string  `json:"exchangeorderid"`
}

type TradeBookResponse struct {
	Exchange        string  `json:"exchange"`
	ProductType     string  `json:"producttype"`
	TradingSymbol   string  `json:"tradingsymbol"`
	InstrumentType  string  `json:"instrumenttype"`
	SymbolGroup     string  `json:"symbolgroup"`
	StrikePrice     float64 `json:"strikeprice,string"`
	OptionType      string  `json:"optiontype"`
	ExpiryDate      string  `json:"expirydate"`
	MarketLot       int     `json:"marketlot,string"`
	Precision       int     `json:"precision,string"`
	Multiplier      float64 `json:"multiplier,string"`
	TradeValue      float64 `json:"tradevalue,string"`
	TransactionType string  `json:"transactiontype"`
	FillPrice       float64 `json:"fillprice,string"`
	FillSize        int     `json:"fillsize,string"`
	OrderID         string  `json:"orderid"`
	FillID          string  `json:"fillid"`
	FillTime        string  `json:"filltime"`
}

type GetPositionsRequest struct {
}

func (c *Client) PlaceOrder(orderReq *PlaceOrderRequest) (*PlaceOrderResponse, error) {

	var orderResp = &PlaceOrderResponse{}

	err := c.doRequest("POST", PLACE_ORDER_ENDPOINT, orderReq, orderResp)
	if err != nil {
		return nil, err
	}
	return orderResp, nil
}

func (c *Client) ModifyOrder(req *ModifyOrderRequest) (*ModifyOrderResponse, error) {

	var modifyResp = &ModifyOrderResponse{}

	err := c.doRequest("POST", MODITY_ORDER_ENDPOINT, req, modifyResp)

	if err != nil {
		return nil, err
	}
	return modifyResp, nil
}

func (c *Client) CancelOrder(req *CancelOrderRequest) error {

	err := c.doRequest("POST", CANCEL_ORDER_ENDPOINT, req, nil)

	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetOrderBook() ([]GetOrderBookResponse, error) {

	var orderBook []GetOrderBookResponse

	err := c.doRequest("GET", ORDER_BOOK_ENDPOINT, nil, &orderBook)

	return orderBook, err
}
