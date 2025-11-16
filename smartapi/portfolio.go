package smartapi

type Holding struct {
	TradingSymbol      string  `json:"tradingsymbol"`
	Exchange           string  `json:"exchange"`
	ISIN               string  `json:"isin"`
	T1Quantity         int64   `json:"t1quantity,string"`
	RealisedQuantity   int64   `json:"realisedquantity,string"`
	Quantity           int64   `json:"quantity,string"`
	AuthorisedQuantity int64   `json:"authorisedquantity,string"`
	Product            string  `json:"product"`
	CollateralQuantity *int64  `json:"collateralquantity,string"`
	CollateralType     *string `json:"collateraltype"`
	Haircut            float64 `json:"haircut,string"`
	AveragePrice       float64 `json:"averageprice,string"`
	LTP                float64 `json:"ltp,string"`
	SymbolToken        string  `json:"symboltoken"`
	Close              float64 `json:"close,string"`
	ProfitAndLoss      float64 `json:"profitandloss,string"`
	PNLPercentage      float64 `json:"pnlpercentage,string"`
}

type TotalHolding struct {
	TotalHoldingValue  float64 `json:"totalholdingvalue"`
	TotalInvValue      float64 `json:"totalinvvalue"`
	TotalProfitAndLoss float64 `json:"totalprofitandloss"`
	TotalPNLPercentage float64 `json:"totalpnlpercentage"`
}

type AllTotalHolding struct {
	Holdings     []Holding `json:"holdings"`
	TotalHolding `json:"totalholding"`
}

type Positions struct {
	Exchange          string  `json:"exchange"`
	SymbolToken       string  `json:"symboltoken"`
	ProductType       string  `json:"producttype"`
	TradingSymbol     string  `json:"tradingsymbol"`
	SymbolName        string  `json:"symbolname"`
	InstrumentType    string  `json:"instrumenttype"`
	PriceDen          string  `json:"priceden"`
	PriceNum          string  `json:"pricenum"`
	GenDen            string  `json:"genden"`
	GenNum            string  `json:"gennum"`
	Precision         string  `json:"precision"`
	Multiplier        string  `json:"multiplier"`
	BoardLotSize      string  `json:"boardlotsize"`
	BuyQty            string  `json:"buyqty"`
	SellQty           string  `json:"sellqty"`
	BuyAmount         string  `json:"buyamount"`
	SellAmount        string  `json:"sellamount"`
	SymbolGroup       string  `json:"symbolgroup"`
	StrikePrice       string  `json:"strikeprice"`
	OptionType        string  `json:"optiontype"`
	ExpiryDate        string  `json:"expirydate"`
	LotSize           string  `json:"lotsize"`
	CfBuyQty          string  `json:"cfbuyqty"`
	CfSellQty         string  `json:"cfsellqty"`
	CfBuyAmount       string  `json:"cfbuyamount"`
	CfSellAmount      string  `json:"cfsellamount"`
	BuyAvgPrice       float64 `json:"buyavgprice,string"`
	SellAvgPrice      float64 `json:"sellavgprice,string"`
	AvgNetPrice       float64 `json:"avgnetprice,string"`
	NetValue          float64 `json:"netvalue,string"`
	NetQty            int64   `json:"netqty,string"`
	TotalBuyValue     float64 `json:"totalbuyvalue,string"`
	TotalSellValue    float64 `json:"totalsellvalue,string"`
	CfBuyAvgPrice     float64 `json:"cfbuyavgprice,string"`
	CfSellAvgPrice    float64 `json:"cfsellavgprice,string"`
	TotalBuyAvgPrice  float64 `json:"totalbuyavgprice,string"`
	TotalSellAvgPrice float64 `json:"totalsellavgprice,string"`
	NetPrice          float64 `json:"netprice,string"`
}

type ConvertPositionRequest struct {
	Exchange        string `json:"exchange"`
	SymbolToken     string `json:"symboltoken"`
	OldProductType  string `json:"oldproducttype"`
	NewProductType  string `json:"newproducttype"`
	TradingSymbol   string `json:"tradingsymbol"`
	SymbolName      string `json:"symbolname"`
	InstrumentType  string `json:"instrumenttype"`
	PriceDen        string `json:"priceden"`
	PriceNum        string `json:"pricenum"`
	GenDen          string `json:"genden"`
	GenNum          string `json:"gennum"`
	Precision       string `json:"precision"`
	Multiplier      string `json:"multiplier"`
	BoardLotSize    string `json:"boardlotsize"`
	BuyQty          string `json:"buyqty"`
	SellQty         string `json:"sellqty"`
	BuyAmount       string `json:"buyamount"`
	SellAmount      string `json:"sellamount"`
	TransactionType string `json:"transactiontype"`
	Quantity        int64  `json:"quantity"`
	Type            string `json:"type"`
}

func (c *Client) GetHoldings() (holdings []Holding, err error) {
	err = c.doRequest("GET", HOLDINGS_ENDPOINT, nil, &holdings)
	return holdings, err
}

func (c *Client) GetAllTotalHoldings() (allTotalHolding AllTotalHolding, err error) {
	err = c.doRequest("GET", HOLDINGS_ENDPOINT+"/total", nil, &allTotalHolding)
	return allTotalHolding, err
}

func (c *Client) GetPositions() (positions []Positions, err error) {

	err = c.doRequest("GET", POSITIONS_ENDPOINT, nil, &positions)
	return positions, err
}

func (c *Client) ConvertPosition(req ConvertPositionRequest) error {
	err := c.doRequest("POST", CONVERT_POSITION_ENDPOINT, req, nil)
	return err
}
