FROM golang:1.14.3-alpine AS build
WORKDIR /src
COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o /out/example "github.com/nazarov-pro/stock-exchange/services/account/cmd" 

FROM scratch AS bin
COPY --from=build /out/example /
CMD ["./example"]
