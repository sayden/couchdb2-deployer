package couchdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func CheckCluster(coord, pass string) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:5984/_membership", coord), nil)
	if err != nil {
		fmt.Printf("ERROR: Couldn't create request cluster membership\n%s\n", err.Error())
		os.Exit(1)
	}

	req.SetBasicAuth("admin", pass)

	client := http.Client{
		Timeout: time.Second * 3,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("ERROR doing request to check membership\n%s\n", err.Error())
		os.Exit(1)
	}

	var mem MembershipResponse
	err = json.NewDecoder(resp.Body).Decode(&mem)
	if err != nil {
		fmt.Println("Couldn't parse JSON response")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("All nodes:")
	for _, node := range mem.AllNodes {
		fmt.Printf("    %s\n", node)
	}

	fmt.Printf("\nCluster nodes:\n")
	for _, node := range mem.AllNodes {
		fmt.Printf("    %s\n", node)
	}
}
