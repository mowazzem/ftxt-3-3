package candle

import (
	"bufio"
	"encoding/json"
	"fmt"
	"ftxt-3-3/model"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type candleHandler struct {
}

type RequestParam struct {
	Year  int `json:"hour"`
	Month int `json:"hour"`
	Day   int `json:"hour"`
	Hour  int `json:"hour"`
}

func NewCandleHandler() *candleHandler {
	return &candleHandler{}
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

	cm, err := getCandleMap()
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
	}
	candles := cm[code]
	respParams := ch.getResponseFromCandleMap(candles, rp)

	respBytes, err := json.Marshal(&respParams)
	if err != nil {
		w.Write([]byte("error occured: " + err.Error()))
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Header().Add("open", strconv.Itoa(respParams.Open))
	w.Header().Add("high", strconv.Itoa(respParams.High))
	w.Header().Add("low", strconv.Itoa(respParams.Low))
	w.Header().Add("close", strconv.Itoa(respParams.Close))

	w.Write(respBytes)
}

type responseParam struct {
	Open  int `json:closee"`
	High  int `json:closee"`
	Low   int `jsonclosese"`
	Close int `json:"close"`
}

func (ch candleHandler) getResponseFromCandleMap(candles []model.Candle, rp RequestParam) responseParam {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	requestTime := time.Date(rp.Year, time.Month(rp.Month), rp.Day, rp.Hour, 0, 0, 0, jst)
	fmt.Println("****", requestTime)

	priceSlice := make([]int, 0)
	fmt.Println("$$", len(candles))

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

	fmt.Println("$$", len(filteredCandles))

	var high, low, open, closePirce int
	if len(filteredCandles) > 0 {
		high = priceSlice[len(priceSlice)-1]
		low = priceSlice[0]
		open = filteredCandles[0].Price
		closePirce = filteredCandles[len(filteredCandles)-1].Price
	}

	respParams := responseParam{
		Open:  open,
		High:  high,
		Low:   low,
		Close: closePirce,
	}

	return respParams

}

func getCandleMap() (model.CandleMap, error) {
	readFile, err := os.Open("./order_books.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	const (
		layout = "2006-01-02 15:04:05"
	)

	jst, _ := time.LoadLocation("Asia/Tokyo")

	mp := make(model.CandleMap)
	for _, line := range fileLines {
		lineParts := strings.Split(line, ",")
		timePart := lineParts[0]
		indx := strings.Index(timePart, " +0900")
		timePart = timePart[0:indx]
		timeStamp, err := time.ParseInLocation(layout, timePart, jst)
		if err != nil {
			log.Fatal("##", err)
		}
		price, err := strconv.Atoi(lineParts[2])
		if err != nil {
			log.Fatal("##", err)
		}

		slc, ok := mp[lineParts[1]]
		if !ok {
			slc = make([]model.Candle, 0)
			slc = append(slc, model.Candle{
				Time:  timeStamp,
				Price: price,
			})
		} else {
			slc = append(slc, model.Candle{
				Time:  timeStamp,
				Price: price,
			})
		}
		mp[lineParts[1]] = slc
	}

	return mp, nil
}
