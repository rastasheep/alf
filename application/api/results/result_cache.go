package results

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/golang/groupcache"
)

const (
	storeKey = "results"
)

type ResultCache struct {
	*groupcache.Group
}

func NewResultCache(logger *log.Logger, size int64, getterFunc resultGetterFunc) *ResultCache {
	getter := &resultGetter{
		logger:     logger,
		getterFunc: getterFunc,
	}

	cache := groupcache.GetGroup(storeKey)
	if cache == nil {
		cache = groupcache.NewGroup(storeKey, size, getter)
	}

	return &ResultCache{cache}
}

func (cache ResultCache) GetResults(key string, data interface{}) error {
	var res []byte
	if err := cache.Get(nil, key, groupcache.AllocatingByteSliceSink(&res)); err != nil {
		return err
	}

	return getInterface(res, data)
}

func getInterface(bts []byte, data interface{}) error {
	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)

	return dec.Decode(data)
}
