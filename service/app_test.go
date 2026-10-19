package service

import (
	"sync"
	"testing"
)

// TestGetAppVersion
func TestGetAppVersion(t *testing.T) {
	s := &AppService{}
	v := s.GetAppVersion()
	// print the result
	t.Logf("App Version: %s", v)
}

func TestMultipleGetAppVersion(t *testing.T) {
	s := &AppService{}
	//concurrency test
	// use a WaitGroup to wait for all goroutines to finish
	wg := sync.WaitGroup{}
	wg.Add(10) // start 10 goroutines
	// start 10 goroutines
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done() // decrement the counter when done
			v := s.GetAppVersion()
			// print the result
			t.Logf("App Version: %s", v)
		}()
	}
	// wait for all goroutines to finish
	wg.Wait()
}
