# Software Engineering School 3.0 Test Assignment 
_authored by Oleksandr Fedorov_

## There are after-deadline improved version in separate branch [`enhancements_after_deadline`](https://github.com/namefreezers/genesis-ses-assignment/tree/enhancements_after_deadline) !!!
I don't know if you consider after-deadline improvements, so I've pushed them to separate branch just in case, leaving branch `main` in before-the-deadline state, so please look into [`enhancements_after_deadline`](https://github.com/namefreezers/genesis-ses-assignment/tree/enhancements_after_deadline) branch if after-deadline improvements matter.

## How to run:
```
git clone https://github.com/namefreezers/genesis-ses-assignment.git
cd genesis-ses-assignment
```
then:
```
docker compose up -d
```
or
```
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

## Description of the solution

### `/api/rate` endpoint
We use few third-party api's for reliability. There are implementation for two third-party services: [_Coinbase_](https://github.com/namefreezers/genesis-ses-assignment/blob/main/fetchbtcrate/coinbase/fetch_price_coinbase.go) and [_Coingecko_](https://github.com/namefreezers/genesis-ses-assignment/blob/main/fetchbtcrate/coingecko/fetch_price_coingecko.go).

In _submitted-before-the-deadline_ solution (branch `main`) we request them one-by-one and return first answered service's value.

And I implemented asyncronous solution after the deadline, in [`enhancements_after_deadline`](https://github.com/namefreezers/genesis-ses-assignment/blob/enhancements_after_deadline/fetchbtcrate/fetch_rate_main.go) branch.

### `/api/subscription` endpoint
There was task not to use external full-featured DB. So we use in-memory `set` as the database. Email entries are quite small (not more that 50 bytes in average), therefore we can store around 20M of subsribed emails per each GB of memory.

In-memory database is backed-up by file-on-disk-database, to preserve subscribed emails between service launches.

(Possible step for enhancement) In current solution file-on-disk is flushed after every `/api/subscribe` request. But if we have a lot of `subscribe` requests, then we need to flush it when 1) `bufio.Writer`'s buffer is full (default functionality) or 2) every 10 seconds (to not to loose last subsribed emails before service shutdown).

### `/api/sendEmails` endpoint
We send emails separately (not in a batch) to have ability to set email header `To:` individually for each email.

In _submitted-before-the-deadline_ solution (branch `main`) we send them one-by-one.

And I implemented asyncronous solution after the deadline, in [`enhancements_after_deadline`](https://github.com/namefreezers/genesis-ses-assignment/blob/enhancements_after_deadline/sendemail/sendemail.go) branch. **_It was crucial for performance._**

## Possible steps for enhancements
1. ([**_Done!_**](https://github.com/namefreezers/genesis-ses-assignment/blob/enhancements_after_deadline/fetchbtcrate/fetch_rate_main.go)) Make fetching btc-uah rate from few third party services _asyncronously_ and use value of that third-party api, whose answer comes first. 
   
   In [current solution](https://github.com/namefreezers/genesis-ses-assignment/blob/main/fetchbtcrate/fetch_rate_main.go) we fetch few third-party api's syncronously, 
   one-by-one and use answer from first available service among few defined services (currently there are implemetations for [_Coinbase_](https://github.com/namefreezers/genesis-ses-assignment/blob/main/fetchbtcrate/coinbase/fetch_price_coinbase.go) and [_Coingecko_](https://github.com/namefreezers/genesis-ses-assignment/blob/main/fetchbtcrate/coingecko/fetch_price_coingecko.go) third-party services).
2. ([**_Done!_**](https://github.com/namefreezers/genesis-ses-assignment/blob/enhancements_after_deadline/sendemail/sendemail.go)) Replase syncronous sending of each emails one-by-one by asyncronous sending. **_Crucial for performance._**

   We need to send all emails one-by-one instead of "batch-sending" because we need to set header `To:` individually.
3. (Still isn't done) Implement some unit tests.

I had to work during workdays and had a plans for the weekend, so if I won't have time to do it before the deadline, I will possibly implement it in the branch `enhancements_after_deadline` till Monday's 29th midday.

## Other ideas for enhancement
1. There is stated in the task description, that we need to implement _one_ service for all these api endpoints.
   But we can divide this api into two microservices: _BtcUahRateService_ and _SubscriptionService_ 
2. If `/api/subscribe` requests are frequent, then move filedb service writing to file new subscribed emails, to separate thread which will flush when 1) buffer of `bufio.Writer` is full and 2) every 10(?) seconds if there was new subscriptions in previous 10 seconds.

   In [current implementation](https://github.com/namefreezers/genesis-ses-assignment/blob/1cfc29b00f3749b5bda1849a9acf9141c33dd052/emailsdb/emailsdb.go#L72), there is flush to file upon every `/api/subsribe` request.
