package main

import (
	"context"
	"fmt"
	"github.com/jobunski/mpesa-go-sdk/internal"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
)

func Init() {

	var b string
	_, b, _, _ = runtime.Caller(0)
	basePath := filepath.Dir(b)

	fmt.Printf("BasePath : %s\n", basePath)

	err := godotenv.Load("development.env")
	if err != nil {
		fmt.Errorf("error while Loading env file %v", err)
	}
}

func main() {

	/**
	Loads the configuration file.
	Either copies the env file tco the root directory or moves the location of the env file
	to the new directory
	*/
	Init()
	mpesa.AssignConfigsToVariables(os.Getenv("APP_KEY"),
		os.Getenv("APP_SECRET"),
		os.Getenv("BASE_URL"),
		os.Getenv("CALLBACK_URL"))

	v := mpesa.StkPushRequest(context.Background(), "174379", "20191213105713", "254710119383",
		"174379", "254710119383", "Test", "Test", 1)
	if v != nil {
		fmt.Printf("Error getting Access Token: %v", v)
	}

}
