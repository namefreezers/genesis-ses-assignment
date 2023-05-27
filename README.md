git clone https://github.com/namefreezers/genesis-ses-assignment.git

cd genesis-ses-assignment

docker build -t genesis-ses-assignment .

docker run -dp 5000:5000 --mount type=bind,src="$(pwd)/emails_data",target="/emails_data" --env-file "./env.list" genesis-ses-assignment



Api examples:

curl http://localhost:5000/api/rate

curl http://localhost:5000/api/subscribe --request "POST" --data '{"email": "your@e.mail"}'

curl http://localhost:5000/api/sendEmails --request "POST"
