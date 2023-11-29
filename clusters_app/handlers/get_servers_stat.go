package handlers

import (
	"bytes"
	"clusters_app/http_client"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type LoadResponse struct {
	Server  string `json:"server"`
	Loading string `json:"loading"`
}

type Result struct {
	Domain   string
	Response LoadResponse
	Err      error
}

var (
	mockServer = map[string]string{
		"a.domain.com": "1.1.1.1",
		"b.domain.com": "2.2.2.2",
		"c.domain.com": "3.3.3.3",
	}
	mockLoad = map[string]string{
		"a.domain.com": "65%",
		"b.domain.com": "15%",
		"c.domain.com": "98%",
	}
	clusters = []string{
		"a.domain.com",
		"b.domain.com",
		"c.domain.com",
	}
)

func GetServerStatHandler(w http.ResponseWriter, r *http.Request) {

	results := make(chan Result, len(clusters))
	var wg sync.WaitGroup

	for _, cluster := range clusters {
		wg.Add(1)
		// Start a goroutine for each domain
		go func(cluster string) {
			defer wg.Done()
			result := Result{Domain: cluster}

			// TODO: because it is not a real app we mock response
			r := LoadResponse{
				Server:  mockServer[cluster],
				Loading: mockLoad[cluster],
			}
			respBytes, err := json.Marshal(r)
			if err != nil {
				result.Err = err
				results <- result
				return
			}
			mockResp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(string(respBytes))),
			}
			mockGetter := &http_client.MockHTTPGetter{Resp: mockResp}
			resp, err := http_client.FetchData(mockGetter, "http://"+cluster+"/load")
			if err != nil {
				result.Err = err
				results <- result
				return
			}
			defer resp.Body.Close()

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				result.Err = err
				results <- result
				return
			}
			r = LoadResponse{}
			err = json.Unmarshal(b, &r)
			if err != nil {
				result.Err = err
				results <- result
				return
			}
			result.Response = r
			results <- result
		}(cluster)
	}

	// Close the results channel once all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	res := []LoadResponse{}
	// Collect results
	for result := range results {
		if result.Err != nil {
			fmt.Printf("Error fetching %s: %s\n", result.Domain, result.Err)
		} else {
			fmt.Printf("Response from %s: %s\n", result.Domain, result.Response)
			res = append(res, result.Response)
		}
	}
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
