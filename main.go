package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

var min_ticker_per_user = 1
var max_ticker_per_user = 10

var min_ticker_price = 100.00
var max_ticker_price = 120.00

var history_day = 90

var possible_tickers = []string{
  "AAPL",
  "MSFT",
  "GOOG",
  "AMZN",
  "FB",
  "TSLA",
  "NVDA",
  "JPM",
  "BABA",
  "JNJ",
  "WMT",
  "PG",
  "PYPL",
  "DIS",
  "ADBE",
  "PFE",
  "V",
  "MA",
  "CRM",
  "NFLX",
}

type ticker_daily_value struct {
    Date    string `json:"date"`
    Price   string `json:"price"`
}

type user_ticker_value struct {
    Symbol	string `json:"symbol"`
    Price	string `json:"price"`
}

func main(){
	router := gin.Default()
	router.GET("/tickers", getTickers)
	router.GET("/tickers/:ticker/history", getTickerHistory)
    router.Run("localhost:8080")
}

func getTickers(c *gin.Context) {
	user_token := checkUser(c)
	if user_token != "" {
		user_tickers_list := GetUserTickers(user_token)
		fmt.Printf("Tickers for %v = %v\n",user_token, user_tickers_list)
		user_tickers_list_with_current_values := GetTickersCurrentValue(user_tickers_list)
	
		c.IndentedJSON(http.StatusOK, user_tickers_list_with_current_values)
	}
}

func getTickerHistory(c *gin.Context) {
	user_token := checkUser(c)
	if user_token != "" {
		ticker_name := c.Param("ticker")

		if contains(possible_tickers,ticker_name) {
			ticker_historycal_values:=GetTickerHistorycalValues(ticker_name)
			c.IndentedJSON(http.StatusOK, ticker_historycal_values)
		} else {
			c.IndentedJSON(http.StatusNotFound, map[string]interface{}{"error": fmt.Sprintf("%v not found!",ticker_name)})
		}
		
	}
}

func checkUser(c *gin.Context) string {
	user_token, _, _:= c.Request.BasicAuth()
	if user_token == "" {
		c.IndentedJSON(http.StatusUnauthorized,map[string]interface{}{"error": "Specify the user token with basic authentication"})
	}
	return user_token
}

func PrintTestData () {
	
	var test_users = []string{
		"Giuse",
		"Giorgia",
		"Diego",
	}
	
	for _, user := range test_users {
		user_tickers_list := GetUserTickers(user)
		fmt.Printf("Tickers for %v = %v\n",user, user_tickers_list)
		user_tickers_list_with_current_values := GetTickersCurrentValue(user_tickers_list)
		PrintTickerValues(user_tickers_list_with_current_values)
	}

	for _, ticker := range possible_tickers {
		ticker_historycal_values:=GetTickerHistorycalValues(ticker)
		fmt.Printf("Values for Ticker %v\n",ticker)
		PrintPrices(ticker_historycal_values)
	}

}

func GetUserTickers(user_token string) []string {
    
	rand.Seed(getSeedFromString(user_token))

	ticker_number := rand.Intn(max_ticker_per_user - min_ticker_per_user) + min_ticker_per_user
	fmt.Printf("%v has %v tickers\n", user_token, ticker_number)

	var owned_tickers []string

	for i := 1; i <= ticker_number ; i++ {
		
		element := possible_tickers[rand.Intn(len(possible_tickers))]
		
		if contains(owned_tickers,element) {
			i--
		} else {
			owned_tickers = append(owned_tickers, element)
		}
	}
	sort.Strings(owned_tickers)
    return owned_tickers
}

func GetTickersCurrentValue (tickers []string) []user_ticker_value {
	now := time.Now()
	current_minute_string := fmt.Sprintf("%d-%02d-%02dT%02d:%02d",now.Year(), now.Month(), now.Day(),now.Hour(), now.Minute())
	
	var user_tickers_list_with_current_values []user_ticker_value
	for _,single_ticker := range tickers {
		user_tickers_list_with_current_values = append(user_tickers_list_with_current_values, GetCurrentTickerValue(single_ticker,current_minute_string))
	}
	return user_tickers_list_with_current_values
}

func GetCurrentTickerValue (ticker_name string, current_minute string) user_ticker_value {

	seed_value := fmt.Sprintf("%v-%v",ticker_name, current_minute)
	rand.Seed(getSeedFromString(seed_value))

	daily_price := min_ticker_price + rand.Float64() * (max_ticker_price - min_ticker_price)
	ticker_price := user_ticker_value{Symbol: ticker_name, Price:fmt.Sprintf("%.2f", daily_price)}

	return ticker_price

}


func GetTickerHistorycalValues (ticker_name string) []ticker_daily_value {

	var ticker_prices []ticker_daily_value
	now := time.Now()

	rand.Seed(getSeedFromString(ticker_name))

	for i := -1; i >= -history_day; i-- {
		app_date := now.AddDate(0, 0, i)
		daily_price := min_ticker_price + rand.Float64() * (max_ticker_price - min_ticker_price)
		ticker_prices = append(ticker_prices, ticker_daily_value{Date: app_date.Format("2006-01-02"), Price:fmt.Sprintf("%.2f", daily_price)})
	}

	return ticker_prices
}

func PrintTickerValues(user_ticker_value_list []user_ticker_value) {

	for _, ticket_curr_value := range user_ticker_value_list {
		fmt.Printf("Tickers %v = %v\n",ticket_curr_value.Symbol, ticket_curr_value.Price)
	}
	
}

func PrintPrices(prices_list []ticker_daily_value) {

	for _, v := range prices_list {
		fmt.Println("Ticker ",v.Date, " price is ", v.Price)
	}

}

func getSeedFromString(token string) int64 {
	h := md5.New()
	io.WriteString(h, token)
	var seed uint64 = binary.BigEndian.Uint64(h.Sum(nil))
	//fmt.Printf("%v has seed %v\n", token, seed)
	return int64(seed)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
