package coinbase

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/namefreezers/genesis-ses-assignment/fetchbtcrate/util"
)

const coinbase_api_url = "https://api.coinbase.com/v2/prices/BTC-UAH/buy"

// Struct to parse response from Coinbase api
type coinbase_parsed_response struct {
	Data coinbase_json_data_field `json:"data"`
}

// inner field `data` of response from Coinbase api
type coinbase_json_data_field struct {
	Btc_uah string `json:"amount"`
}

func coinbase_parse_price_from_response_body(body_bytes []byte) (float64, error) {
	var coinbase_resp_struct coinbase_parsed_response = coinbase_parsed_response{}
	err := json.Unmarshal(body_bytes, &coinbase_resp_struct)
	if err != nil {
		return 0, fmt.Errorf("can't unmarshall: %v", string(body_bytes))
	}

	res_price_float, err := strconv.ParseFloat(coinbase_resp_struct.Data.Btc_uah, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse price from string: %v", coinbase_resp_struct.Data.Btc_uah)
	}

	return res_price_float, nil
}

func FetchBtcUahRate() (float64, error) {
	resp_body, err := util.Request_and_get_resp_body(coinbase_api_url)
	// if there is an error while requesting url or getting body from response, return error
	if err != nil {
		return 0, err
	}

	res_price_float, err := coinbase_parse_price_from_response_body(resp_body)
	// Return error if response's body doesn't match scheme
	if err != nil {
		return 0, err
	}

	return res_price_float, nil
}
