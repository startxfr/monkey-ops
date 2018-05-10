package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"os"
)

//Get all the pods running from a project
func GetPods(token string, project string, url string) []string {

	urlGetPods := url + "/api/v1/namespaces/" + project + "/pods"

	// Set up the HTTP request to get pods
	req, err := http.NewRequest("GET", urlGetPods, nil)
	req.Header.Add("Authorization", "Bearer "+token)

	if err != nil {
		panic(err)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{
		Transport: transport,
	}

	resp, err := cli.Do(req)

	if err != nil {
		log.Println("API_SERVER=" + url)
		log.Fatal("Fail getting Pods")
	}

	defer resp.Body.Close()

	pods, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	podsName := make([]string, 0)
	podsCustom := map[string]interface{}{}
	json.Unmarshal(pods, &podsCustom)

	if podsCustom != nil  && len(podsCustom)> 0 {
		items := podsCustom["items"].([]interface{})

		for _, item := range items {
			itemObject := item.(map[string]interface{})
			metadataMap := itemObject["metadata"].(map[string]interface{})
			statusMap := itemObject["status"].(map[string]interface{})
			status := statusMap["phase"].(string)
			if status == "Running" {
				podsName = append(podsName, metadataMap["name"].(string))
			}

		}
	}

	return podsName
}

//Delete a running pod
func DeletePod(pod string, chaosInput *ChaosInput) {
	start := time.Now()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{
		Transport: transport,
	}

	urlDeletePod := chaosInput.Url + "/api/v1/namespaces/" + chaosInput.Project + "/pods/" + pod

	// Set up the HTTP request to delete pod
	req, err := http.NewRequest("DELETE", urlDeletePod, nil)
	req.Header.Add("Authorization", "Bearer "+chaosInput.Token)
	resp, err := cli.Do(req)

	if err != nil {
		log.Println("API_SERVER=" + chaosInput.Url)
		log.Fatal("Fail deleting Pod")
	}

	defer resp.Body.Close()

	log.Printf(
		"%s\t%s",
		"deleted pod: "+pod,
		time.Since(start),
	)

}

//Get all the DeploymentConfig from a project
func GetDCs(chaosInput *ChaosInput) []DcObject {

	urlGetDCs := chaosInput.Url + "/oapi/v1/namespaces/" + chaosInput.Project + "/deploymentconfigs"

	// Set up the HTTP request to get DCs
	req, err := http.NewRequest("GET", urlGetDCs, nil)
	req.Header.Add("Authorization", "Bearer "+chaosInput.Token)

	if err != nil {
		panic(err)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{
		Transport: transport,
	}

	resp, err := cli.Do(req)

	if err != nil {
		log.Println("API_SERVER=" + chaosInput.Url)
		log.Fatal("Fail getting DeploymentConfigs")
	}

	defer resp.Body.Close()

	dcs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	dcsName := make([]DcObject, 0)
	dcsCustom := map[string]interface{}{}
	json.Unmarshal(dcs, &dcsCustom)

	if dcsCustom != nil && len(dcsCustom)> 0 {
		items := dcsCustom["items"].([]interface{})

		for _, item := range items {
			itemObject := item.(map[string]interface{})
			metadataMap := itemObject["metadata"].(map[string]interface{})
			specMap := itemObject["spec"].(map[string]interface{})
			dcsName = append(dcsName, DcObject{metadataMap["name"].(string), specMap["replicas"].(float64)})
		}
	}

	return dcsName
}

//Scale down a DC if the number of replicas > o or scale up a DC if number of replicas = 0
func scaleDC(dc string, chaosInput *ChaosInput, replicas float64) {

	start := time.Now()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cli := &http.Client{
		Transport: transport,
	}

	urlScaleDC := chaosInput.Url + "/oapi/v1/namespaces/" + chaosInput.Project + "/deploymentconfigs/" + dc + "/scale"

	// Set up the HTTP request to scale DC
	metadata := Metadata{
		Name:      dc,
		Namespace: chaosInput.Project}
	spec := Spec{
		Replicas: replicas}
	scale := Scale{
		Kind:       "Scale",
		ApiVersion: "extensions/v1beta1",
		Metadata:   metadata,
		Spec:       spec}

	body, err := json.Marshal(scale)

	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("PUT", urlScaleDC, bytes.NewReader(body))
	req.Header.Add("Authorization", "Bearer "+chaosInput.Token)

	resp, err := cli.Do(req)

	if err != nil {
		log.Println("API_SERVER=" + chaosInput.Url)
		log.Fatal("Fail deleting DeploymentConfig")
	}

	defer resp.Body.Close()

	log.Printf(
		"%s\t%s",
		"scaled dc "+dc+" to "+strconv.FormatFloat(replicas, 'f', -1, 64),
		time.Since(start),
	)

}

//method to launch the chaos
func ExecuteChaos(chaosInput *ChaosInput, mode string) {

	start := time.Now()
	rand.Seed(time.Now().Unix())

	for doChaos := (mode == "background" || (time.Since(start).Seconds() < chaosInput.TotalTime)); doChaos; doChaos = (mode == "background" || (time.Since(start).Seconds() < chaosInput.TotalTime)) {
		
		//Randomly choice if delete pod or scale a DC
		randComponent := random(1, 3)
		//Use selected method if given
        if (chaosInput.Method != "" && chaosInput.Method != "random")  {
          switch chaosInput.Method {
          case "pod":
              randComponent = 1
          case "dc":
              randComponent = 2
          }
        }

		switch randComponent {
		case 1:
			pods := GetPods(chaosInput.Token, chaosInput.Project, chaosInput.Url)
			if pods != nil && len(pods) > 0 {
				randPod := random(0, len(pods))
				DeletePod(pods[randPod], chaosInput)
			}
		case 2:
			dcs := GetDCs(chaosInput)
			if dcs != nil && len(dcs) > 0 {
				randDc := random(0, len(dcs))
				replicas := dcs[randDc].Replicas
				if replicas > 0 {
					replicas--
				} else {
					replicas++
				}
				//To avoid Monkey-ops atack itself
				if dcs[randDc].Name ==  os.Getenv("APP_NAME"){
					log.Println("Prevent Monkey-Ops from attacking itself")
					if randDc == 0 {
						randDc ++
					} else {
						randDc --
					}
				}
				scaleDC(dcs[randDc].Name, chaosInput, replicas)
			}
		}

		//Waiting for the next monkey action
		time.Sleep(time.Second * time.Duration(chaosInput.Interval))
	}

}
