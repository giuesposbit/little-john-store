package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine{
    router := gin.Default()
    return router
}

func TestTickersUnauthorized(t *testing.T) {
	mockResponse, _ := PrettyString(`{"error":"Specify the user token with basic authentication"}`)
    r := SetUpRouter()
    r.GET("/tickers", getTickers)
    req, _ := http.NewRequest("GET", "/tickers", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    responseData, _ := io.ReadAll(w.Body)
    assert.Equal(t, mockResponse, string(responseData))
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTickers(t *testing.T) {
    r := SetUpRouter()
    r.GET("/tickers", getTickers)
    req, _ := http.NewRequest("GET", "/tickers", nil)
	req.Header.Add("Authorization","Basic " + basicAuth("pippo","")) 

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

   	// Convert the JSON response to a map
   	var response []user_ticker_value
   	err := json.Unmarshal(w.Body.Bytes(), &response)
   	// Grab the value & whether or not it exists

	assert.NotEmpty(t,response)

	assert.Nil(t, err)
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestFakeTickerHistory(t *testing.T) {
    r := SetUpRouter()
	r.GET("/tickers/:ticker/history", getTickerHistory)
    req, _ := http.NewRequest("GET", "/tickers/FAKE_TICKER/history", nil)
	req.Header.Add("Authorization","Basic " + basicAuth("pippo","")) 

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTickerHistory(t *testing.T) {
    r := SetUpRouter()
	r.GET("/tickers/:ticker/history", getTickerHistory)
    req, _ := http.NewRequest("GET", "/tickers/BABA/history", nil)
	req.Header.Add("Authorization","Basic " + basicAuth("pippo","")) 

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

	// Convert the JSON response to a map
	var response []ticker_daily_value
	err := json.Unmarshal(w.Body.Bytes(), &response)
	// Grab the value & whether or not it exists

	assert.True(t,len(response) == 90)

 	assert.Nil(t, err)
 	assert.Equal(t, http.StatusOK, w.Code)
}

func PrettyString(str string) (string, error) {
    var prettyJSON bytes.Buffer
    if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
        return "", err
    }
    return prettyJSON.String(), nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}