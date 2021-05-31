FROM golang:1.14.3-alpine AS build

RUN apk update && apk add --no-cache git tzdata
# RUN apk add --no-cache --virtual protobuf-compiler protoc
# RUN go get google.golang.org/protobuf/cmd/protoc-gen-go
# RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

ENV USER=appuser
ENV UID=1000
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR $GOPATH/github.com/nazarov-pro/stock-exchange

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

# COPY services/account/pkg/domain/pb services/account/pkg/domain/pb
# RUN protoc -I="./services/account/pkg/domain/pb" --go_out=./ ./services/account/pkg/domain/pb/email.proto

COPY services/migrate services/migrate
COPY configs/apps/wallet.yaml /app/config.yaml
COPY configs/db/migrations/wallet /db-migrations
ARG DB_HOST_ARG="127.0.0.1"
ENV __DB_MIGRATION_SOURCE_URL="file:///db-migrations"
ENV __DB_HOST=$DB_HOST_ARG
RUN CONFIG_FILE="/app/config.yaml" go run "github.com/nazarov-pro/stock-exchange/services/migrate/cmd"

COPY pkg/ pkg/
COPY services/wallet services/wallet
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/wallet github.com/nazarov-pro/stock-exchange/services/wallet/cmd

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /app /app

EXPOSE 8080
ENV CONFIG_FILE="/app/config.yaml"

USER appuser:appuser

ENTRYPOINT ["/app/wallet"]
