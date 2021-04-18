FROM alpine:3.13 AS builder
RUN apk add --no-cache go
WORKDIR /go
COPY . /go/src/github.com/denysvitali/elastically
WORKDIR /go/src/github.com/denysvitali/elastically
RUN go build -o /usr/bin/elastically ./cmd

FROM alpine:3.13
COPY --from=builder /usr/bin/elastically /usr/bin/elastically
ENTRYPOINT ["elastically"]