package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EnvType struct {
	Url      string `json:"url"`
	FileName string `json:"fileName"`
}

func loadEnv() EnvType {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("can't read .env: %v", err)
	}
	url := os.Getenv("URL")
	fileName := os.Getenv("FILE")
	res := EnvType{Url: url, FileName: fileName}
	return res
}
