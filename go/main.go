package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"time"

	"github.com/gocolly/colly"
)

type UserGPA struct {
	GPA float64 `json:"gpa"`
}

type Lab struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Capacity int       `json:"capacity"`
	Year     int       `json:"year"`
	UserGPAs []UserGPA `json:"userGPAs"`
}

type LabsData struct {
	Labs []Lab `json:"labs"`
}

type LabsDataHistory struct {
	LabsData LabsData `json:"labsData"`
	Time     string   `json:"time"`
}

type LabsDataHistorys struct {
	LabsDataHistorys []LabsDataHistory `json:"labsDataHistorys"`
}

func main() {
	env := loadEnv()
	url, fileName := env.Url, env.FileName

	now := time.Now()
	waitSeconds := 60 - now.Second()
	initialDelay := time.Duration(waitSeconds) * time.Second
	time.Sleep(initialDelay)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		updateLabData(url, fileName)
	}
}

func updateLabData(url string, fileName string) {
	c := colly.NewCollector()
	var data LabsData

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		err := json.Unmarshal(r.Body, &data)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		appendData(data, fileName)
	})
	c.Visit(url)

}

func appendData(newData LabsData, fileName string) {
	var history LabsDataHistorys
	now := time.Now().Format(time.RFC3339)
	newLabsDataHistory := LabsDataHistory{
		LabsData: newData,
		Time:     now,
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&history)
	if err != nil && err != io.EOF {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	if err == io.EOF {
		history = LabsDataHistorys{LabsDataHistorys: make([]LabsDataHistory, 0)}
	}

	history.LabsDataHistorys = append(history.LabsDataHistorys, newLabsDataHistory)

	file.Seek(0, 0)
	file.Truncate(0)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(history); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
	}
	fmt.Printf("write %v\n", now)
}
