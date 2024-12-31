package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var path = ""

func main() {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))

	type CSVRecord struct {
		// 日付
		Date string
		// 店舗名
		StoreName string
		// 金額
		Amount string
	}

	records := make([]*CSVRecord, 0)

	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		records = append(records, &CSVRecord{
			Date:      record[0],
			StoreName: record[1],
			Amount:    record[2],
		})
	}

	amountByStoreName := make(map[string]int, len(records))
	for idx, record := range records {
		// csv の header は skip する
		if idx == 0 {
			continue
		}
		// 金額が空文字になっているものは無視（ポイント消費など）
		if record.Amount == "" {
			continue
		}

		i, err := strconv.Atoi(record.Amount)
		if err != nil {
			panic(err)
		}
		// 消費の分析が目的のため、割引は無視する
		if i < 0 {
			continue
		}

		amountByStoreName[record.StoreName] = amountByStoreName[record.StoreName] + i
	}

	sum := 0
	for storeName, totalAmount := range amountByStoreName {
		sum = sum + totalAmount
		fmt.Println(storeName, ": ", totalAmount)
	}
	fmt.Println("合計: ", sum)
}
