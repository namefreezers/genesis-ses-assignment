package fetchbtcrate

import (
	"fmt"

	"github.com/namefreezers/genesis-ses-assignment/fetchbtcrate/coinbase"
	"github.com/namefreezers/genesis-ses-assignment/fetchbtcrate/coingecko"
)

func FetchBtcUahRateMain() (float64, error) {
	rate, err := coinbase.FetchBtcUahRate()
	if err == nil {
		return rate, nil
	}

	rate, err = coingecko.FetchBtcUahPrice()
	if err == nil {
		return rate, nil
	}

	return 0, fmt.Errorf("All btc rate services are unavailable. Last fetch error: %v", err.Error())
}
