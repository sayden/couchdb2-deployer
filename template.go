package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TemplateExecutor interface {
	Do(coord string) error
}

// CouchDbHTTP is a Template pattern to allow easy communication with the HTTP socket of a CouchDB
// instance.
type CouchDbHTTP struct {
	ActionPerformer ClusterActionPerformer
	HTTPRequester
	Finish bool
}

func (te *CouchDbHTTP) SetActionPerformer(t ClusterActionPerformer) {
	te.ActionPerformer = t
}

func (t *CouchDbHTTP) Do(coord string) error {
	enabler := t.ActionPerformer.GetRequest()

	req, err := t.getRequest(coord, enabler)
	if err != nil {
		return err
	}

	client := t.getClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	ok := ClusterResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ok)
	if err != nil {
		return err
	}

	t.ActionPerformer.HandleResponse(&ok)

	return nil
}

type RetryTemplate struct {
	ActionPerfomer ClusterActionPerformer
	Template       TemplateExecutor
	Coordinators   []string
	CandidateNode  Node
	Wait           int
}

func (r *RetryTemplate) Do(candidateNode string) error {
	r.Template = &CouchDbHTTP{
		ActionPerformer: &EnableHostAction{
			Node: r.CandidateNode,
		},
	}

	for {
		for _, coord := range r.Coordinators {
			if r.CandidateNode.host == coord {
				//Don't try to connect to itself
				continue
			}

			if err := Remove(coord, r.CandidateNode); err != nil {
				fmt.Println(err)
			}

			if err := JoinAllClusterAction(r.CandidateNode.user, r.CandidateNode.pass,
				coord, []string{candidateNode}, r.CandidateNode.port); err != nil {
				fmt.Println(err)
				continue
			} else {
				//template := &CouchDbHTTP{}
				//node := r.CandidateNode
				//node.SetHost(r.Coordinators[0])
				//template.SetActionPerformer(&FinishAction{node})
				//
				//return template.Do(r.Coordinators[0])
				return nil
			}

			time.Sleep(300 * time.Millisecond)
		}

		fmt.Printf("All coordination nodes have been tried. Waiting %d seconds before trying again.\n\n\n", r.Wait)
		time.Sleep(time.Duration(r.Wait) * time.Second)
	}
}

type HTTPRequester struct{}

func (h HTTPRequester) getClient() *http.Client {
	client := &http.Client{
		Timeout: time.Second * 3,
	}

	return client
}

func (h HTTPRequester) getRequest(coord string, c *ClusterRequest) (*http.Request, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	body, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%d/_cluster_setup", coord, c.Port), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = header

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	return req, nil
}
