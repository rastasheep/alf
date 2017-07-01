package results

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang/groupcache"
)

type resultGetterFunc func(int) ([]map[string]interface{}, error)

type resultGetter struct {
	logger     *log.Logger
	getterFunc resultGetterFunc
}

func (g *resultGetter) Get(ctx groupcache.Context, key string, dest groupcache.Sink) error {
	id, err := strconv.Atoi(key)
	if err != nil {
		return fmt.Errorf("failed to parse result cache key: %v", err)
	}
	g.logger.Printf("Filling cache for: %v", id)

	results, err := g.getterFunc(id)
	if err != nil {
		return err
	}

	resultsBytes, err := getBytes(results)
	if err != nil {
		return err
	}
	g.logger.Printf("Cache size for key: %v, %v Bytes", key, binary.Size(resultsBytes))

	return dest.SetBytes(resultsBytes)
}

func getBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	gob.Register(time.Time{})
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
