package coingecko

import (
	"encoding/json"
	"fmt"

	"github.com/namefreezers/genesis-ses-assignment/fetchbtcrate/util"
)

const coingecko_api_url = "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah&precision=2"

// Struct to parse response from Coingecko api
type coingecko_parsed_response struct {
	Data coingecko_json_bitcoin_field `json:"bitcoin"`
}

// inner field `bitcoin` of response from Coingecko api
type coingecko_json_bitcoin_field struct {
	Btc_uah float64 `json:"uah"`
}

func coingecko_parse_price_from_response_body(body_bytes []byte) (float64, error) {
	var coingecko_resp_struct coingecko_parsed_response = coingecko_parsed_response{}
	err := json.Unmarshal(body_bytes, &coingecko_resp_struct)
	if err != nil {
		return 0, fmt.Errorf("can't unmarshall: %v", string(body_bytes))
	}

	return coingecko_resp_struct.Data.Btc_uah, nil
}

func FetchBtcUahPrice() (float64, error) {
	resp_body, err := util.Request_and_get_resp_body(coingecko_api_url)
	// if there is an error while requesting url or getting body from response, return error
	if err != nil {
		return 0, err
	}

	res_price_float, err := coingecko_parse_price_from_response_body(resp_body)
	// Return error if response's body doesn't match scheme
	if err != nil {
		return 0, err
	}

	return res_price_float, nil
}
