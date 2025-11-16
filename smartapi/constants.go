package smartapi

import "time"

//Endpoints constants
const (
	BASE_URL                  = "https://apiconnect.angelone.in"
	REQUEST_TIMEOUT           = 30 * time.Second
	LOGIN_ENDPOINT            = "/rest/auth/angelbroking/user/v1/loginByPassword"
	USER_PROFILE_ENDPOINT     = "/rest/secure/angelbroking/user/v1/getProfile"
	MARGIN_ENDPOINT           = "/rest/secure/angelbroking/user/v1/getRMS"
	LOGOUT_ENDPOINT           = "/rest/secure/angelbroking/user/v1/logout"
	PLACE_ORDER_ENDPOINT      = "/rest/secure/angelbroking/order/v1/placeOrder"
	MODITY_ORDER_ENDPOINT     = "/rest/secure/angelbroking/order/v1/modifyOrder"
	CANCEL_ORDER_ENDPOINT     = "/rest/secure/angelbroking/order/v1/cancelOrder"
	ORDER_BOOK_ENDPOINT       = "/rest/secure/angelbroking/order/v1/getOrderBook"
	TRADE_BOOK_ENDPOINT       = "/rest/secure/angelbroking/order/v1/getTradeBook"
	POSITIONS_ENDPOINT        = "/rest/secure/angelbroking/portfolio/v1/positions"
	HOLDINGS_ENDPOINT         = "/rest/secure/angelbroking/portfolio/v1/holdings"
	CONVERT_POSITION_ENDPOINT = "/rest/secure/angelbroking/portfolio/v1/convertPosition"
	CANDLE_DATA_ENDPOINT      = "/rest/secure/angelbroking/historical/v1/getCandleData"
)

//Order related constants
const (
	VARIETY_NORMAL             = "NORMAL"
	VARIETY_STOPLOSS           = "STOPLOSS"
	VARIETY_ROBO               = "ROBO"
	TRANSACTION_TYPE_BUY       = "BUY"
	TRANSACTION_TYPE_SELL      = "SELL"
	ORDER_TYPE_MARKET          = "MARKET"
	ORDER_TYPE_LIMIT           = "LIMIT"
	ORDER_TYPE_STOPLOSS_LIMIT  = "STOPLOSS_LIMIT"
	ORDER_TYPE_STOPLOSS_MARKET = "STOPLOSS_MARKET"
	PRODUCT_TYPE_INTRADAY      = "INTRADAY"
	PRODUCT_TYPE_DELIVERY      = "DELIVERY"
	PRODUCT_TYPE_CARRYFORWARD  = "CARRYFORWARD"
	PRODUCT_TYPE_MARGIN        = "MARGIN"
	PRODUCT_TYPE_BO            = "BO"
	DURATION_DAY               = "DAY"
	DURATION_IOC               = "IOC"
	EXCHANGE_NSE               = "NSE"
	EXCHANGE_BSE               = "BSE"
	EXCHANGE_CDS               = "CDS"
	EXCHANGE_MCX               = "MCX"
	EXCHANGE_NFO               = "NFO"
	INTERVAL_1MIN              = "ONE_MINUTE"
	INTERVAL_3MIN              = "THREE_MINUTE"
	INTERVAL_5MIN              = "FIVE_MINUTE"
	INTERVAL_10MIN             = "TEN_MINUTE"
	INTERVAL_15MIN             = "FIFTEEN_MINUTE"
	INTERVAL_30MIN             = "THIRTY_MINUTE"
	INTERVAL_1HOUR             = "ONE_HOUR"
	INTERVAL_1DAY              = "ONE_DAY"
)
