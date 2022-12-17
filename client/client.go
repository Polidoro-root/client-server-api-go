package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type DollarExchange struct {
	Bid string `json:"bid"`
}

const GET_DOLLAR_EXCHANGE_REQUEST_DEADLINE = time.Millisecond * 300

func main() {
	ctx := context.Background()

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(GET_DOLLAR_EXCHANGE_REQUEST_DEADLINE))

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)

	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var dollarExchange DollarExchange

	err = json.Unmarshal(body, &dollarExchange)

	if err != nil {
		panic(err)
	}

	var file *os.File

	_, err = os.Stat("cotacao.txt")

	if err != nil {
		file, err = os.Create("cotacao.txt")

		if err != nil {
			panic(err)
		}

	} else {
		file, err = os.OpenFile("cotacao.txt", os.O_WRONLY, os.ModeAppend)

		if err != nil {
			panic(err)
		}

	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", dollarExchange.Bid))

	if err != nil {
		panic(err)
	}

}
