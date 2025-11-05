package smartapi

import (
	"fmt"
	"log"
	"time"

	"github.com/pquerna/otp/totp"
)

type loginRequest struct {
	ClientCode string `json:"clientcode"`
	Password   string `json:"password"`
	TOTP       string `json:"totp"`
	State      string `json:"state"`
}

type loginResponse struct {
	JWTToken     string `json:"jwtToken"`
	RefreshToken string `json:"refreshToken"`
	FeedToken    string `json:"feedToken"`
	State        string `json:"state"`
}

type UserProfileResponse struct {
	ClientCode    string   `json:"clientcode"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	MobileNo      string   `json:"mobileno"`
	Exchanges     []string `json:"exchanges"`
	Products      []string `json:"products"`
	LastLoginTime string   `json:"lastlogintime"`
	BrokerID      string   `json:"brokerid"`
}

type UserMarginResponse struct {
	Net                    string `json:"net"`
	AvailableCash          string `json:"availablecash"`
	AvailableIntradayPayin string `json:"availableintradaypayin"`
	AvailableLimitMargin   string `json:"availablelimitmargin"`
	Collateral             string `json:"collateral"`
	M2MUnrealized          string `json:"m2munrealized"`
	M2MRealized            string `json:"m2mrealized"`
	UtilisedDebits         string `json:"utiliseddebits"`
	UtilisedSpan           string `json:"utilisedspan"`
	UtilisedOptionPremium  string `json:"utilisedoptionpremium"`
	UtilisedHoldingSales   string `json:"utilisedholdingsales"`
	UtilisedExposure       string `json:"utilisedexposure"`
	UtilisedTurnover       string `json:"utilisedturnover"`
	UtilisedPayout         string `json:"utilisedpayout"`
}

func (c *Client) GenerateSession() error {

	totp, err := totp.GenerateCode(c.TOTPKey, time.Now())
	if err != nil {
		return fmt.Errorf("Error Generating TOTP %s", err.Error())
	}

	var loginReq = loginRequest{
		ClientCode: c.ClientId,
		Password:   c.Mpin,
		TOTP:       totp,
		State:      "1",
	}

	var lr loginResponse

	err = c.doRequest("POST", LOGIN_ENDPOINT, loginReq, &lr)
	if err != nil {
		return err
	}

	log.Printf("Login Response %+v", lr)

	c.SetTokens(lr.JWTToken, lr.RefreshToken, lr.FeedToken)

	return nil
}

func (c *Client) SetTokens(jwtToken, refreshToken, feedToken string) {
	c.JWTToken = jwtToken
	c.RefreshToken = refreshToken
	c.FeedToken = feedToken
}

func (c *Client) UserProfile() (*UserProfileResponse, error) {

	var profile = &UserProfileResponse{}

	err := c.doRequest("GET", USER_PROFILE_ENDPOINT, nil, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (c *Client) UserMargin() (*UserMarginResponse, error) {

	var margin = &UserMarginResponse{}

	err := c.doRequest("GET", MARGIN_ENDPOINT, nil, margin)
	if err != nil {
		return nil, err
	}

	return margin, nil
}
