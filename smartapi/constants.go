package smartapi

import "time"

//Endpoints constants
const (
	BASE_URL              = "https://apiconnect.angelone.in"
	REQUEST_TIMEOUT       = 30 * time.Second
	LOGIN_ENDPOINT        = "/rest/auth/angelbroking/user/v1/loginByPassword"
	USER_PROFILE_ENDPOINT = "/rest/secure/angelbroking/user/v1/getProfile"
	MARGIN_ENDPOINT       = "/rest/secure/angelbroking/user/v1/getRMS"
	LOGOUT_ENDPOINT       = "/rest/secure/angelbroking/user/v1/logout"
	PLACE_ORDER_ENDPOINT  = "/rest/secure/angelbroking/order/v1/placeOrder"
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
)
