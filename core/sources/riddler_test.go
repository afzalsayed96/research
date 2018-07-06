package sources

import core "github.com/subfinder/research/core"
import "testing"
import "sync"
import "fmt"

func TestRiddler(t *testing.T) {
	domain := "bing.com"
	source := Riddler{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		results = append(results, result)
	}

	if !(len(results) >= 9) {
		t.Errorf("expected more than 9 result(s), got '%v'", len(results))
	}
}

func TestRiddlerMultiThreaded(t *testing.T) {
	domains := []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com"}
	source := Riddler{}
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

	if len(results) < 30 {
		t.Errorf("expected more than 30 results, got '%v'", len(results))
	}
}

