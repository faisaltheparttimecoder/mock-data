package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ielizaga/mockd/core"
)

func insert(key string, jsonData string) error {
	url := fmt.Sprintf("http://%v:%v/gemfire-api/v1/%v/%v", Connector.HostName, Connector.Port, Connector.RegionName, key)
	//    log.Infof("URL: %v", url)

	//    log.Info(fmt.Sprintf("Inserting Object with key '%v' into Region '%v'...", key, Connector.RegionName))
	var jsonStr = []byte(jsonData)
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return err
	}

	defer response.Body.Close()

	//    log.Info(fmt.Sprintf("Inserting Object with key '%v' into Region '%v'... Done!. Response[Status=%v]", key, Connector.RegionName, response.Status))
	return nil
}

// GemFire Data Mock
func MockGemFire() error {
	log.Info(fmt.Sprintf("Inserting Mock Objects into Region '%v'...", Connector.RegionName))
	core.ProgressBar(Connector.RowCount, "GemFire")
	for i := 0; i < Connector.RowCount; i++ {
		v, _ := core.RandomInt(1, 10000)
		d, _ := core.RandomDate(-15, 50)

		insert(fmt.Sprintf("key_%v", i), fmt.Sprintf("{ \"@type\": \"io.pivotal.support.model.DomainClass\", \"key\": \"key_%v\", \"amount\": %v, \"dateTime\":  \"%v\" }", i, v, d))
		core.IncrementBar()
	}

	log.Info(fmt.Sprintf("Inserting Mock Objects into Region '%v'... Done!", Connector.RegionName))
	return nil
}
