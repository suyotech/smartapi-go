package instruments

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Instrument struct {
	Token          string `json:"token"`
	Symbol         string `json:"symbol"`
	Name           string `json:"name"`
	Expiry         string `json:"expiry"`
	Strike         string `json:"strike"`
	LotSize        string `json:"lotsize"`
	InstrumentType string `json:"instrumenttype"`
	Exchange       string `json:"exch_seg"`
	TickSize       string `json:"tick_size"`
}

func DownloadInstruments() error {
	resp, err := http.Get("https://margincalculator.angelone.in/OpenAPI_File/files/OpenAPIScripMaster.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Optionally, validate JSON
	var instruments []Instrument
	if err := json.Unmarshal(body, &instruments); err != nil {
		return err
	}

	// Save to file
	return os.WriteFile("instruments.json", body, 0644)
}

// CheckInstrumentsUpdate checks if instruments.json was modified before today 8:30 AM IST.
// If yes, it downloads a new file.
func CheckInstrumentsUpdate() error {
	const fileName = "instruments.json"

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return err
	}

	now := time.Now().In(loc)
	today830 := time.Date(now.Year(), now.Month(), now.Day(), 8, 30, 0, 0, loc)

	info, err := os.Stat(fileName)
	if err != nil {
		// If file does not exist, download it
		if os.IsNotExist(err) {
			return DownloadInstruments()
		}
		return err
	}

	modTime := info.ModTime().In(loc)
	if modTime.Before(today830) {
		return DownloadInstruments()
	}

	return nil
}

func LoadInstruments() ([]Instrument, error) {
	data, err := os.ReadFile("instruments.json")
	if err != nil {
		return nil, err
	}

	var instruments []Instrument
	if err := json.Unmarshal(data, &instruments); err != nil {
		return nil, err
	}

	return instruments, nil
}

type FindInstrumentRequest struct {
	Token          string `json:"token,omitempty"`
	Symbol         string `json:"symbol,omitempty"`
	Name           string `json:"name,omitempty"`
	Expiry         string `json:"expiry,omitempty"`
	Strike         string `json:"strike,omitempty"`
	InstrumentType string `json:"instrumenttype,omitempty"`
	Exchange       string `json:"exch_seg,omitempty"`
}

func FindInstruments(req FindInstrumentRequest, instruments []Instrument) ([]Instrument, error) {
	var err error
	if instruments == nil {
		instruments, err = LoadInstruments()
		if err != nil {
			return nil, err
		}
	}
	var results []Instrument
	for _, inst := range instruments {
		match := true
		if req.Token != "" && inst.Token != req.Token {
			match = false
		}
		if req.Name != "" && inst.Name != req.Name {
			match = false
		}
		if req.Symbol != "" && inst.Symbol != req.Symbol {
			match = false
		}
		if req.Expiry != "" && inst.Expiry != req.Expiry {
			match = false
		}
		if req.Strike != "" && inst.Strike != req.Strike {
			match = false
		}
		if req.InstrumentType != "" && inst.InstrumentType != req.InstrumentType {
			match = false
		}
		if req.Exchange != "" && inst.Exchange != req.Exchange {
			match = false
		}
		if match {
			results = append(results, inst)
		}
	}
	return results, nil
}

type OptionStrikes struct {
	ATM float64   `json:"atm"`
	ITM []float64 `json:"itm"`
	OTM []float64 `json:"otm"`
}

func GetOptionStrikes(price float64, req FindInstrumentRequest, instruments []Instrument) (strikes OptionStrikes, err error) {

	if instruments == nil {
		var err error
		instruments, err = LoadInstruments()
		if err != nil {
			return strikes, err
		}
	}

	foundInstruments, err := FindInstruments(req, instruments)
	if err != nil {
		return strikes, err
	}
	strikeSet := make(map[float64]Instrument)
	for _, inst := range foundInstruments {
		strikeValue, err := strconv.ParseFloat(inst.Strike, 64)
		strikeDiff := (strikeValue) / 100
		if err == nil {
			strikeSet[strikeDiff] = inst
		}
	}

	// Logic to determine ATM, ITM, and OTM strikes based on the price and request
	// This is a placeholder implementation

	return strikes, nil
}

func GetExpiryDates(req FindInstrumentRequest, instruments []Instrument) ([]string, error) {
	if instruments == nil {
		var err error
		instruments, err = LoadInstruments()
		if err != nil {
			return nil, err
		}
	}

	expiryMap := make(map[string]struct{})
	for _, inst := range instruments {
		if inst.Symbol == req.Symbol && inst.Exchange == req.Exchange {
			expiryMap[inst.Expiry] = struct{}{}
		}
	}

	var expiryDates []time.Time
	expiryStrToTime := make(map[time.Time]string)
	for expiry := range expiryMap {
		// Try to parse in "02Jan2006" format (e.g., 27JAN2026)
		t, err := time.Parse("02Jan2006", strings.ToUpper(expiry))
		if err == nil {
			expiryDates = append(expiryDates, t)
			expiryStrToTime[t] = expiry
		}
	}

	// Sort by date ascending
	sort.Slice(expiryDates, func(i, j int) bool {
		return expiryDates[i].Before(expiryDates[j])
	})

	// Convert back to original string format
	result := make([]string, len(expiryDates))
	for i, t := range expiryDates {
		result[i] = expiryStrToTime[t]
	}

	return result, nil
}
