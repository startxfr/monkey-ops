package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	
	"github.com/spf13/viper"
	flag "github.com/spf13/pflag"
)

func main() {
	
	flag.String("REGION_URL", "https://api.boae.paas.gsnetcloud.corp:8443", "Region URL")
	flag.String("PROJECT_NAME", "devstack-dev", "Project to get crazy")
	flag.Float64("INTERVAL", 20, "interval time in seconds")
	flag.Float64("TOTAL_TIME", 0, "total time of chaos monkey in seconds")
	flag.String("MODE", "background", "Execution mode: background or rest")
	
	
	//Binding flags and env vars
	viper.BindPFlag( "REGION_URL", flag.Lookup("REGION_URL") )
	viper.BindPFlag( "PROJECT_NAME", flag.Lookup("PROJECT_NAME") )
	viper.BindPFlag( "INTERVAL", flag.Lookup("INTERVAL") )
	viper.BindPFlag( "TOTAL_TIME", flag.Lookup("TOTAL_TIME") )
	viper.BindPFlag( "MODE", flag.Lookup("MODE") )
	
	viper.BindEnv("REGION_URL")
	viper.BindEnv("PROJECT_NAME")
	viper.BindEnv("INTERVAL")
	viper.BindEnv("TOTAL_TIME")
	viper.BindEnv("MODE")
	
	flag.Parse()
	
	//set configuration
	regionUrl := viper.GetString("REGION_URL")
	project := viper.GetString("PROJECT_NAME")
	interval := viper.GetFloat64("INTERVAL")
	totalTime := viper.GetFloat64("TOTAL_TIME")
	mode := viper.GetString("MODE")
	
	// read the service account secret token file at once
    tokenBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
    if err != nil {
        panic(err)
    }

    token := string(tokenBytes[:])
    fmt.Println(token);
	
	chaosInput:= ChaosInput{
		Url: regionUrl,
		Project: project,
		Token: token,
		Interval: interval,
		TotalTime: totalTime,
	}
	
	go ExecuteChaos(&chaosInput, mode)

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}