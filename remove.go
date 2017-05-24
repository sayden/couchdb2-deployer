package couchdb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
)

func Remove(coord string, n Node) error {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s:5986/_nodes/couchdb@%s", coord, n.host), nil)
	if err != nil {
		return err
	}
	req.Header = header

	req.SetBasicAuth(n.user, n.pass)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var clusterResp ClusterResponse
	if err := json.Unmarshal(byt, &clusterResp); err != nil {
		return err
	}

	if clusterResp.Error != "" {
		return errors.New(clusterResp.Reason)
	} else {
		return nil
	}
}
