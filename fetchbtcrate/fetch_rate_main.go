package fetchbtcrate

import (
	"fmt"
	"time"

	"github.com/namefreezers/genesis-ses-assignment/fetchbtcrate/coinbase"
	"github.com/namefreezers/genesis-ses-assignment/fetchbtcrate/coingecko"
)

type btcUahFetcherFunc func() (float64, error)

func fetchBtcUahAsync(fetcher btcUahFetcherFunc, ch chan<- float64) {
	rate, err := fetcher()
	if err == nil {
		ch <- rate
	}
}

// Try to fetch btc-uah rate from few third-party api's one-by-one synchronously
// (There is an option for enhancement, to make request asyncronously,
// and return response from first answered third-party api.)
func FetchBtcUahRateMain() (float64, error) {

	ch := make(chan float64)

	// run few btc-uah rate fetches asyncronously
	go fetchBtcUahAsync(coinbase.FetchBtcUahRate, ch)
	go fetchBtcUahAsync(coingecko.FetchBtcUahPrice, ch)

	// 10 second timeout
	timeout := time.After(10 * time.Second)
	select {
	case rate := <-ch:
		// return first responded service's answer
		return rate, nil
	case <-timeout:
		// return error, if we reached the timeout
		fmt.Println("timed out all requests")
		return 0, fmt.Errorf("all btc rate services are unavailable")
	}
}
