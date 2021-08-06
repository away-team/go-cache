# go-cache

This is a cache package that enables uses to switch between backend storage solutions.  It provides Get, Set, Delete functionality as well as a Passthrough cache optimized GetAndLoad function.


## Usage

Get, Set and Delete are pretty self explanitory but GetAndLoad provides some unique benefits.  If a requested item is not in the cache it will use the provided closure to fetch that item, cache it and return it.  Additionally an concurrent requests for that same item will wait until the item is fetched and cached.  This prevents potential cache runs and gaurentees that the closure is only called once to repopulate the cache.

```golang
import (
    "log"
	"github.com/away-team/go-cache/src/cache"
	"github.com/away-team/go-cache/src/storage/memory"
	"github.com/away-team/go-cache/src/storage/redis"
)

...

    // memory storage setup
    defaultExpiration := time.Second * 3600 //how long before items expire by default
    cleanupInterval := time.Second * 7200  //how frequently expired items are purged from the cache (see https://github.com/patrickmn/go-cache)
    storage := memory.NewStorage(defaultExpiration, cleanupInterval) // setup an in memory cache
    
    // redis storage setup   
    ring := redis.NewRing(&redis.RingOptions{ // see https://github.com/go-redis/cache readme for redis options
        Addrs: map[string]string{
            "server1": ":6379",
            "server2": ":6380",
        },
    })
    storage = rcache.NewStorage(ring)

    // setup the cache with memory or redis storage...
	c := cache.New(storage, defaultExpiration, false, &log.Logger{})

    //Example GetAndLoad
    loader := func() (interface{}, error) {
		// do work to fetch fresh data to cache

        // return that fresh data or an error if one occurred. 
		return "foo", nil
	}
    var result string
	err := c.GetAndLoad("foobar", &result, loader)
    if err != nil {
        ...
    }
    log.Printf(result) // prints "foo"
    ...

```
