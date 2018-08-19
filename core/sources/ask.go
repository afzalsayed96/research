package sources

import (
	"bufio"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/subfinder/research/core"
)

// Ask is a source to process subdomains from https://ask.com
type Ask struct{}

// ProcessDomain takes a given base domain and attempts to enumerate subdomains.
func (source *Ask) ProcessDomain(domain string) <-chan *core.Result {
	results := make(chan *core.Result)
	go func(domain string, results chan *core.Result) {
		defer close(results)

		domainExtractor, err := core.NewSubdomainExtractor(domain)
		if err != nil {
			results <- core.NewResult("ask", nil, err)
			return
		}

		uniqFilter := map[string]bool{}

		for currentPage := 1; currentPage <= 750; currentPage++ {
			resp, err := core.HTTPClient.Get("https://www.ask.com/web?q=site%3A" + domain + "+-www.+&page=" + strconv.Itoa(currentPage) + "&o=0&l=dir&qsrc=998&qo=pagination")
			if err != nil {
				results <- core.NewResult("ask", nil, err)
				return
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				results <- core.NewResult("ask", nil, errors.New(resp.Status))
				return
			}

			scanner := bufio.NewScanner(resp.Body)

			scanner.Split(bufio.ScanWords)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			for scanner.Scan() {
				for _, str := range domainExtractor.FindAllString(scanner.Text(), -1) {
					_, found := uniqFilter[str]
					if !found {
						uniqFilter[str] = true
						select {
						case results <- core.NewResult("ask", str, nil):
							// move along
						case <-ctx.Done():
							resp.Body.Close()
							return
						}
					}
				}
			}

			resp.Body.Close()
		}

	}(domain, results)
	return results
}
