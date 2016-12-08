package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Template struct {
	Template ClusterActionPerformer
}

func (te *Template) SetTemplate(t ClusterActionPerformer) {
	te.Template = t
}

func (t *Template) Do(coord string) {
	enabler := t.Template.GetRequest()

	req, err := t.getRequest(coord, enabler)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := t.getClient()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ok := ClusterResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ok)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t.Template.HandleResponse(&ok)
}

func (t *Template) getClient() *http.Client {
	client := &http.Client{
		Timeout: time.Second * 3,
	}

	return client
}

func (t *Template) getRequest(coord string, c *ClusterRequest) (*http.Request, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:5984/_cluster_setup", coord), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = header

	req.SetBasicAuth(c.Username, c.Password)

	return req, nil
}
