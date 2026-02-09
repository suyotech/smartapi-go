package smartapi

import "time"

type CandleDataRequest struct {
	Exchange    string `json:"exchange"`
	SymbolToken string `json:"symboltoken"`
	Interval    string `json:"interval"`
	FromDate    string `json:"fromdate"` // Format: "YYYY-MM-DD HH:MM"
	ToDate      string `json:"todate"`   // Format: "YYYY-MM-DD HH:MM"
}

type Candle struct {
	Time   time.Time `json:"time"`
	Open   float64   `json:"open,string"`
	High   float64   `json:"high,string"`
	Low    float64   `json:"low,string"`
	Close  float64   `json:"close,string"`
	Volume int64     `json:"volume,string,omitempty"`
	OI     int64     `json:"oi,string,omitempty"`
}

func (c *Client) GetCandleData(req *CandleDataRequest) ([]Candle, error) {

	var candleData []Candle
	var rawCandles [][]interface{}
	err := c.doRequest("POST", CANDLE_DATA_ENDPOINT, req, &rawCandles)
	if err != nil {
		return nil, err
	}
	for _, rc := range rawCandles {
		length := len(rc)
		oi := int64(0)
		if length > 6 {
			oi = int64(rc[6].(float64))
		}
		volume := int64(0)
		if length > 5 {
			volume = int64(rc[5].(float64))
		}
		// Convert timestamp string to Go time.Time
		timestampStr := rc[0].(string)
		loc, _ := time.LoadLocation("Asia/Kolkata")
		parsedTime, err := time.ParseInLocation("2006-01-02T15:04:05-07:00", timestampStr, loc)
		if err != nil {
			return nil, err
		}
		candle := Candle{
			Time:   parsedTime,
			Open:   rc[1].(float64),
			High:   rc[2].(float64),
			Low:    rc[3].(float64),
			Close:  rc[4].(float64),
			Volume: volume,
			OI:     oi,
		}
		candleData = append(candleData, candle)
	}
	return candleData, nil
}

type OptionGreek struct {
	Name        string  `json:"name"`
	Expiry      string  `json:"expiry"`
	StrikePrice float64 `json:"strikePrice,string"`
	Optiontype  string  `json:"optionType"`
	Delta       float64 `json:"delta,string"`
	Gamma       float64 `json:"gamma,string"`
	Theta       float64 `json:"theta,string"`
	Vega        float64 `json:"vega,string"`
	IV          float64 `json:"impliedVolatility,string"`
	TradeVolume float64 `json:"tradeVolume,string"`
}

func (c *Client) GetOptionGreeks(name string, expirydate string) ([]OptionGreek, error) {

	req := struct {
		Name       string `json:"name"`
		ExpiryDate string `json:"expirydate"`
	}{
		Name:       name,
		ExpiryDate: expirydate,
	}

	var result []OptionGreek
	err := c.doRequest("POST", OPTION_GREEKS_ENDPOINT, req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
