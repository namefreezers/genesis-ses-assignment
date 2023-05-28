# Software Engineering School 3.0 Test Assignment 
_authored by Oleksandr Fedorov_

## How to run:
```
git clone https://github.com/namefreezers/genesis-ses-assignment.git
cd genesis-ses-assignment
docker build -t genesis-ses-assignment .
docker run -dp 5000:5000 --mount type=bind,src="$(pwd)/emails_data",target="/emails_data" --env-file "./env.list" genesis-ses-assignment
```

## Api examples:
Fetch BTC-UAH rate:
```
curl http://localhost:5000/api/rate
```

Subscribe new email:
```
curl http://localhost:5000/api/subscribe --request "POST" --data '{"email": "your@e.mail"}'
```

Request mailing to all subscribed emails:
```
curl http://localhost:5000/api/sendEmails --request "POST"
```

## Possible steps for enhancements
1. Make fetching btc-uah rate from few third party services _asyncronously_ and use value of that third-party api, whose answer comes first. 
   In [current solution](https://github.com/namefreezers/genesis-ses-assignment/blob/main/fetchbtcrate/fetch_rate_main.go) we fetch few third-party api's syncronously, 
   one-by-one and use answer from first available service among few defined services (currently there are implemetations for [_Coinbase_](https://github.com/namefreezers/genesis-ses-assignment/blob/main/fetchbtcrate/coinbase/fetch_price_coinbase.go) and [_Coingecko_](https://github.com/namefreezers/genesis-ses-assignment/blob/main/fetchbtcrate/coingecko/fetch_price_coingecko.go) third-party services).
2. Implement some unit tests.

I had to work during workdays and had a plans for the weekend, so if I won't have time to do it before the deadline, I will possibly implement it in the branch `enhancements_after_deadline` till Monday's 29th midday.

## Ideas for enhancement
1. There is stated in the task description, that we need to implement _one_ service for all these api endpoints.
   But we can divide this api into two microservices: _BtcUahRateService_ and _SubscriptionService_ 
