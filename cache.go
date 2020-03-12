package main

import (
	"github.com/allegro/bigcache"
	"log"
	"time"
)

// NewDumper - create new dumper
func NewCache(cacheConfig cacheConfig) (cache *bigcache.BigCache) {
	config := bigcache.Config{
		Shards:             cacheConfig.Shards,
		LifeWindow:         time.Duration(cacheConfig.LifeWindow) * time.Minute,
		CleanWindow:        time.Duration(cacheConfig.CleanWindow) * time.Minute,
		MaxEntriesInWindow: cacheConfig.MaxEntriesInWindow,
		MaxEntrySize:       cacheConfig.MaxEntrySize,
		Verbose:            cacheConfig.Verbose,
		HardMaxCacheSize:   cacheConfig.HardMaxCacheSize,
	}

	cache, initErr := bigcache.NewBigCache(config)
	if initErr != nil {
		log.Fatal(initErr)
	}
	return cache
}
