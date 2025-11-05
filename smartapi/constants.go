package smartapi

import "time"

const (
	BASE_URL              = "https://apiconnect.angelone.in"
	REQUEST_TIMEOUT       = 30 * time.Second
	LOGIN_ENDPOINT        = "/rest/auth/angelbroking/user/v1/loginByPassword"
	USER_PROFILE_ENDPOINT = "/rest/secure/angelbroking/user/v1/getProfile"
	MARGIN_ENDPOINT       = "/rest/secure/angelbroking/user/v1/getRMS"
)
