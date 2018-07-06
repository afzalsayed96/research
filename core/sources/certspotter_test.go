package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestCertSpotter(t *testing.T) {
	domain := "google.com"
	source := CertSpotter{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		if result.IsFailure() {
			t.Fatal(result.Failure)
		}
		results = append(results, result)
	}

	if !(len(results) >= 5000) {
		t.Errorf("expected more than 5000 results, got '%v'", len(results))
	}
}

func TestCertSpotter_MultiThreaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := CertSpotter{}
	results := []*core.Result{}

	wg := sync.WaitGroup{}
	mx := sync.Mutex{}

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			for result := range source.ProcessDomain(domain) {
				mx.Lock()
				results = append(results, result)
				mx.Unlock()
			}
		}(domain)
	}

	wg.Wait() // collect results

	if len(results) < 67000 {
		t.Errorf("expected more than 67000 results, got '%v'", len(results))
	}
}

