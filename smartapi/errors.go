package smartapi

import (
	"errors"
	"fmt"
	"strings"
)

var ErrRateLimited = errors.New("SmartAPI rate limit exceeded")

type HTTPError struct {
	StatusCode int
	Message    string
	RawBody    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("SmartAPI HTTP %d: %s", e.StatusCode, e.Message)
}

func (e *HTTPError) Unwrap() error {
	text := strings.ToLower(e.Message + " " + e.RawBody)
	if strings.Contains(text, "access denied because of exceeding access rate") ||
		strings.Contains(text, "too many requests") ||
		strings.Contains(text, "ab1021") {
		return ErrRateLimited
	}
	return nil
}

var apiErrorMessages = map[string]string{
	"AG8001": "Invalid Token",
	"AG8002": "Token Expired",
	"AG8003": "Token missing",
	"AB8050": "Invalid Refresh Token",
	"AB8051": "Refresh Token Expired",
	"AB1000": "Invalid Email Or Password",
	"AB1001": "Invalid Email",
	"AB1002": "Invalid Password Length",
	"AB1003": "Client Already Exists",
	"AB1004": "Something Went Wrong, Please Try After Sometime",
	"AB1005": "User Type Must Be USER",
	"AB1006": "Client Is Block For Trading",
	"AB1007": "AMX Error",
	"AB1008": "Invalid Order Variety",
	"AB1009": "Symbol Not Found",
	"AB1010": "AMX Session Expired",
	"AB1011": "Client not login",
	"AB1012": "Invalid Product Type",
	"AB1013": "Order not found",
	"AB1014": "Trade not found",
	"AB1015": "Holding not found",
	"AB1016": "Position not found",
	"AB1017": "Position conversion failed",
	"AB1018": "Failed to get symbol details",
	"AB1050": "Invalid TOTP",
	"AB2000": "Error not specified",
	"AB2001": "Internal Error, Please try after sometime",
	"AB1031": "Old Password Mismatch",
	"AB1032": "User Not Found",
	"AB2002": "ROBO order is block",
	"AB4008": "ordertag length should be less than 20 characters",
	"AG7002": "Unregistered IP Address",
}

func getAPIErrorMessage(errorCode string) string {
	if msg, exists := apiErrorMessages[errorCode]; exists {
		return msg
	}
	return "Unknown error"
}
