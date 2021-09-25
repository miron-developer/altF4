###################################################################################################
#                                                                                                 #
#                                   Miron-developer                                               #
#                                                                                                 #
###################################################################################################

FROM golang:1.16

COPY . .
WORKDIR /pkg
RUN go mod download; go build -o ./cmd/wnet cmd/main.go

LABEL description="This is the binance project." \
    authors="Miron-developer" \
    contacts="https://github.com/miron-developer"

CMD ["cmd/binance"]

EXPOSE 4430