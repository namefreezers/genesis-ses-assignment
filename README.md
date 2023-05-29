# Software Engineering School 3.0 Test Assignment 
_authored by Oleksandr Fedorov_

## There are after-deadline improved version in separate branch [`enhancements_after_deadline`](https://github.com/namefreezers/genesis-ses-assignment/tree/enhancements_after_deadline) !!!
I don't know if you consider after-deadline improvements, so I've pushed them to separate branch just in case, leaving branch `main` in before-the-deadline state, so please look into [`enhancements_after_deadline`](https://github.com/namefreezers/genesis-ses-assignment/tree/enhancements_after_deadline) branch if after-deadline improvements matter.

I had to work during workdays and had a plans for the weekend, so I pushed some enhancements here on Monday 29th.

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

We fetch these services asyncronously. Which service respond first, value of that service we return.

### `/api/subscription` endpoint
There was task not to use external full-featured DB. So we use in-memory `set` as the database. Email entries are quite small (not more that 50 bytes in average), therefore we can store around 20M of subsribed emails per each GB of memory.

In-memory database is backed-up by file-on-disk-database, to preserve subscribed emails between service launches.

(Possible step for enhancement) In current solution file-on-disk is flushed after every `/api/subscribe` request. But if we have a lot of `subscribe` requests, then we need to flush it when 1) `bufio.Writer`'s buffer is full (default functionality) or 2) every 10 seconds (to not to loose last subsribed emails before service shutdown).

### `/api/sendEmails` endpoint
We send emails separately (not in a batch) to have ability to set email header `To:` individually for each email.

We send all these emails asyncronously, so sending emails not in a batch doesn't affect perfomance a lot.

## Possible steps for enhancements
1. Implement some unit tests.

## Other ideas for enhancement
1. There is stated in the task description, that we need to implement _one_ service for all these api endpoints.
   But we can divide this api into two microservices: _BtcUahRateService_ and _SubscriptionService_ 
2. If `/api/subscribe` requests are frequent, then move filedb service writing to file new subscribed emails, to separate thread which will flush when 1) buffer of `bufio.Writer` is full and 2) every 10(?) seconds if there was new subscriptions in previous 10 seconds.

   In [current implementation](https://github.com/namefreezers/genesis-ses-assignment/blob/1cfc29b00f3749b5bda1849a9acf9141c33dd052/emailsdb/emailsdb.go#L72), there is flush to file upon every `/api/subsribe` request.
