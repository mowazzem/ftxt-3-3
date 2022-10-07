package candle

import (
	"encoding/json"
	"ftxt-3-3/model"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type candleHandler struct {
	candleMap *model.CandleMap
}

type RequestParam struct {
	Year  int `json:"hour"`
	Month int `json:"hour"`
	Day   int `json:"hour"`
	Hour  int `json:"hour"`
}

func NewCandleHandler(cm *model.CandleMap) *candleHandler {
	return &candleHandler{
		candleMap: cm,
	}
}

func (ch *candleHandler) GetCandle(w http.ResponseWriter, r *http.Request) {
	var rp RequestParam

	code := r.URL.Query().Get("code")
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")
	dayStr := r.URL.Query().Get("day")
	hourStr := r.URL.Query().Get("hour")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
	}

	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
	}

	rp.Year = year
	rp.Month = month
	rp.Day = day
	rp.Hour = hour

	cm := *ch.candleMap
	candles := cm[code]
	respParams := ch.getResponseFromCandleMap(candles, rp)

	respBytes, err := json.Marshal(&respParams)
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
		return
	}

	w.Header().Add("Content-type", "application/json")

	w.Write(respBytes)
}

type responseParam struct {
	Open  int `json:closee"`
	High  int `json:closee"`
	Low   int `jsonclosese"`
	Close int `json:"close"`
}

func (ch candleHandler) getResponseFromCandleMap(candles []model.Candle, rp RequestParam) responseParam {
	requestTime := time.Date(rp.Year, time.Month(rp.Month), rp.Day, rp.Hour, 0, 0, 0, time.Local)

	priceSlice := make([]int, 0)

	filteredCandles := make([]model.Candle, 0)

	for _, c := range candles {
		if c.Time.After(requestTime) && c.Time.Before(requestTime.Add(1*time.Hour)) {
			filteredCandles = append(filteredCandles, c)
		}
	}

	for _, c := range filteredCandles {
		priceSlice = append(priceSlice, c.Price)
	}

	sort.Slice(priceSlice, func(i, j int) bool {
		return priceSlice[i] < priceSlice[j]
	})

	sort.Slice(filteredCandles, func(i, j int) bool {
		return filteredCandles[i].Time.Before(filteredCandles[j].Time)
	})

	respParams := responseParam{
		Open:  filteredCandles[0].Price,
		High:  priceSlice[len(priceSlice)-1],
		Low:   priceSlice[0],
		Close: filteredCandles[len(filteredCandles)-1].Price,
	}

	return respParams

}
