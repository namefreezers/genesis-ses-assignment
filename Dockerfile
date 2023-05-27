FROM golang:1.20.4-alpine3.18 as build-stage

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/namefreezers/genesis-ses-assignment

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /genesis-ses-assignment .

#
# final build stage
#
FROM scratch

# Copy ca-certs for app web access
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /genesis-ses-assignment /genesis-ses-assignment

# app uses port 5000
EXPOSE 5000

WORKDIR /

ENTRYPOINT ["/genesis-ses-assignment"]




# FROM golang:1.20.4-alpine3.18

# RUN apk --no-cache add ca-certificates

# WORKDIR /go/src/github.com/namefreezers/genesis-ses-assignment

# COPY . .

# RUN CGO_ENABLED=0 GOOS=linux go build -a -o /genesis-ses-assignment .

# # app uses port 5000
# EXPOSE 5000

# WORKDIR /

# ENTRYPOINT ["/genesis-ses-assignment"]