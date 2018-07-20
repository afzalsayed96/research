package sources

import (
	"bufio"
	"net"
	"net/http"
	"time"

	core "github.com/subfinder/research/core"
)

type CertSpotter struct{}

type certspotterObject struct {
	DNSNames []string `json:"dns_names"`
}

func (source *CertSpotter) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		httpClient := &http.Client{
			Timeout: time.Second * 60,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		}

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- &core.Result{Type: "certspotter", Failure: err}
			return
		}

		uniqFilter := map[string]bool{}

		// get response from the API, optionally with an API key
		resp, err := httpClient.Get("https://certspotter.com/api/v0/certs?domain=" + domain)
		if err != nil {
			results <- &core.Result{Type: "certspotter", Failure: err}
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
				_, found := uniqFilter[str]
				if !found {
					uniqFilter[str] = true
					results <- &core.Result{Type: "threatminer", Success: str}
				}
			}
		}

	}(domain, results)
	return results
}
