package main

import (
	"log"
	"net/http"
	
	flag "github.com/spf13/pflag"
)

func main() {
	
	regionUrl := flag.String("regionUrl", "https://api.boae.paas.gsnetcloud.corp:8443", "Region URL")
	project := flag.String("project", "devstack-dev", "Project to get crazy")
	token := flag.String("token", "MD94AKuyJw8Is2IHJ3Ciktns3dnDcrorzhHJmwFNDQ4", "Service Account Token")
	interval := flag.Float64("interval", 20, "interval time in seconds")
	totalTime := flag.Float64("totalTime", 0, "total time of chaos monkey in seconds")
	mode := flag.String("mode", "background", "Execution mode: background or rest")
	
	flag.Parse()
	
	chaosInput:= ChaosInput{
		Url: *regionUrl,
		Project: *project,
		Token: *token,
		Interval: *interval,
		TotalTime: *totalTime,
	}
	
	go ExecuteChaos(&chaosInput, *mode)

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}