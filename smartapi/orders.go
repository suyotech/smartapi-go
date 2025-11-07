package smartapi



type PlaceOrderResponse struct {
	Script        string `json:"script"`
	OrderID       string `json:"orderid"`
	UniqueOrderID string `json:"uniqueorderid"`
}

type PlaceOrderRequest struct {
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

func (c *Client) PlaceOrder(orderReq *PlaceOrderRequest) (*PlaceOrderResponse, error) {

	var orderResp = &PlaceOrderResponse{}

	err := c.doRequest("POST", PLACE_ORDER_ENDPOINT, orderReq, orderResp)
	if err != nil {
		return nil, err
	}
	return orderResp, nil
}
