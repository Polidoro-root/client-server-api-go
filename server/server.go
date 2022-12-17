package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type UsdBrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type AwesomeapiUsdBrlDto struct {
	UsdBrl UsdBrl `json:"USDBRL"`
}

type GetDollarExchangeDto struct {
	Bid string `json:"bid"`
}

const GET_DOLLAR_EXCHANGE_HANDLER_DEADLINE = time.Millisecond * 200

const SAVE_DOLLAR_EXCHANGE_DEADLINE = time.Millisecond * 10

func main() {
	createDatabase()

	mux := http.NewServeMux()

	mux.HandleFunc("/cotacao", getDollarExchangeHandler)

	http.ListenAndServe(":8080", mux)
}

func getDollarExchangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(GET_DOLLAR_EXCHANGE_HANDLER_DEADLINE))

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		println(err)
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var awesomeapiUsdBrl AwesomeapiUsdBrlDto

	err = json.Unmarshal(body, &awesomeapiUsdBrl)

	if err != nil {
		panic(err)
	}

	getDollarExchange := &GetDollarExchangeDto{
		Bid: awesomeapiUsdBrl.UsdBrl.Bid,
	}

	err = saveDollarExchange(getDollarExchange.Bid)

	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(w).Encode(getDollarExchange)

	if err != nil {
		panic(err)
	}

}

func createDatabase() {
	db, err := sql.Open("sqlite3", "./currency_exchange.db")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS exchanges (id INTEGER PRIMARY KEY AUTOINCREMENT, usdbrl VARCHAR(10))")

	if err != nil {
		panic(err)
	}
}

func saveDollarExchange(dollarExchange string) error {
	ctx := context.Background()

	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(SAVE_DOLLAR_EXCHANGE_DEADLINE))

	defer cancel()

	db, err := sql.Open("sqlite3", "./currency_exchange.db")

	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)

	defer tx.Rollback()

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO exchanges (usdbrl) VALUES (?)", dollarExchange)

	if err != nil {

		return err
	}

	return tx.Commit()
}
