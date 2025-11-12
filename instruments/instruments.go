package instruments

import (
	"encoding/json"
	"io"
	"math"
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

const (
	INTTYPE_EQ     = ""       // Equity
	INTTYPE_INDEX  = "AMXIDX" // Index
	INTTYPE_FUTIDX = "FUTIDX" // Future Index
	INTTYPE_FUTSTX = "FUTSTX" // Future Stock
	INTTYPE_OPTIDX = "OPTIDX" // Option Index
	INTTYPE_OPTSTX = "OPTSTX" // Option Stock
	INTTYPE_FUTCOM = "FUTCOM" // Future Commodity
	INTTYPE_OPTFUT = "OPTFUT" // Option Future
)

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

	// If results have expiry, sort by expiry date ascending
	hasExpiry := false
	for _, inst := range results {
		if inst.Expiry != "" {
			hasExpiry = true
			break
		}
	}
	if hasExpiry {
		sort.Slice(results, func(i, j int) bool {
			ti, erri := time.Parse("02Jan2006", strings.ToUpper(results[i].Expiry))
			tj, errj := time.Parse("02Jan2006", strings.ToUpper(results[j].Expiry))
			if erri != nil && errj != nil {
				return results[i].Expiry < results[j].Expiry
			}
			if erri != nil {
				return false
			}
			if errj != nil {
				return true
			}
			return ti.Before(tj)
		})
	}

	return results, nil
}

type OptionStrikes struct {
	NEAR float64   `json:"near"`
	UP   []float64 `json:"up"`
	DOWN []float64 `json:"down"`
}

func GetOptionStrikes(price float64, req FindInstrumentRequest, maxStrikes int, instruments []Instrument) (strikes OptionStrikes, err error) {
	if instruments == nil {
		instruments, err = LoadInstruments()
		if err != nil {
			return strikes, err
		}
	}

	foundInstruments, err := FindInstruments(req, instruments)
	if err != nil {
		return strikes, err
	}

	var strikesSet = make(map[float64]bool)

	for _, inst := range foundInstruments {
		if inst.Strike != "" {
			strikeValue, err := strconv.ParseFloat(inst.Strike, 64)
			if err == nil {
				strikesSet[strikeValue] = true
			}
		}
	}

	var strikeValues []float64
	for strike := range strikesSet {
		strikeValues = append(strikeValues, strike/100.0) // Convert to actual price
	}

	sort.Float64s(strikeValues)

	//Find the nearest strike
	nearestStrike := strikeValues[0]
	minDiff := math.MaxFloat64
	for _, strike := range strikeValues {
		diff := math.Abs(price - strike)
		if diff < minDiff {
			minDiff = diff
			nearestStrike = strike
		}
	}

	strikes.NEAR = nearestStrike

	strikes.UP = []float64{}
	strikes.DOWN = []float64{}

	// Collect maxStrikes strikes above and below nearestStrike
	strikes.UP = []float64{}
	strikes.DOWN = []float64{}

	//Collect up strikes
	for _, strike := range strikeValues {
		if strike > nearestStrike && len(strikes.UP) < maxStrikes {
			strikes.UP = append(strikes.UP, strike)
		}
	}

	//Collect down strikes in reverse order
	for i := len(strikeValues) - 1; i >= 0; i-- {
		strike := strikeValues[i]
		if strike < nearestStrike && len(strikes.DOWN) < maxStrikes {
			strikes.DOWN = append(strikes.DOWN, strike)
		}
	}

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

	foundInstruments, err := FindInstruments(req, instruments)
	if err != nil {
		return nil, err
	}

	expiryMap := make(map[string]struct{})
	for _, inst := range foundInstruments {
		expiryMap[inst.Expiry] = struct{}{}
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
