package sources

import (
	"fmt"
	"sync"
	"testing"

	"github.com/subfinder/research/core"
)

func TestThreatCrowd(t *testing.T) {
	domain := "google.com"
	source := ThreatCrowd{}
	results := []*core.Result{}

	for result := range source.ProcessDomain(domain) {
		//t.Log(result)
		fmt.Println(result)
		results = append(results, result)
		// Not waiting around to iterate all the possible pages.
		if len(results) >= 20 {
			break
		}
	}

	if !(len(results) >= 20) {
		t.Errorf("expected more than 20 result(s), got '%v'", len(results))
	}
}
