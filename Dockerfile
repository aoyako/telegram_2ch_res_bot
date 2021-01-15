FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk --no-cache add ca-certificates

WORKDIR /build

COPY . .
RUN go mod download

RUN go build -o main cmd/main.go


WORKDIR /dist

RUN cp /build/main .
RUN cp -r /build/configs .
RUN cp -r /build/src .


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /dist /


ENTRYPOINT ["/main"]