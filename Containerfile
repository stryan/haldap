FROM golang:1.22-alpine as builder
WORKDIR /go/src/app
COPY . .
#RUN apt update && apt upgrade -y
RUN CGO_ENABLED=0 GOOS=linux go build -a

FROM homeassistant/home-assistant:stable
COPY --from=builder /go/src/app/haldap /root/

