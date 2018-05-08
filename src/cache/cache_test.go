package cache_test

import (
	"log"
	"testing"
	"time"

	"github.com/away-team/go-cache/src/cache"
	"github.com/away-team/go-cache/src/storage/memory"
)

func Test_GetAndLoad(t *testing.T) {
	c := cache.New(memory.NewStorage(time.Minute, time.Minute), time.Minute, false, &log.Logger{})

	var ret, ret2 string
	loader := func() (interface{}, error) {
		time.Sleep(time.Second * 2)
		return "foo", nil
	}
	start := time.Now()
	err := c.GetAndLoad("foobar", &ret, loader)
	if err != nil {
		t.Fatalf("Error in GetAndLoad: %v", err)
	}
	middle := time.Now()
	err = c.GetAndLoad("foobar", &ret2, loader)
	if err != nil {
		t.Fatalf("Error in GetAndLoad2: %v", err)
	}
	end := time.Now()

	if middle.Sub(start) < (time.Second * 2) {
		t.Fatalf("loading was shorter than possible")
	}
	if end.Sub(middle) > time.Second {
		t.Fatalf("fetching from cache was way too slow")
	}

}
