package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"github.com/golang/groupcache"
	"strconv"
	"time"
)

type ResultsCache struct {
	*groupcache.Group
}

type resultsGetter struct {
	s *Server
}

func (s *Server) initCache(size int64) {
	cache := groupcache.GetGroup("results")
	if cache == nil {
		g := resultsGetter{s}
		cache = groupcache.NewGroup("results", size, g)
	}
	s.resultsCache = &ResultsCache{cache}
}

func (g resultsGetter) Get(ctx groupcache.Context, key string, dest groupcache.Sink) error {
	id, err := strconv.Atoi(key)
	if err != nil {
		return err
	}

	e := Execution{ID: id}
	db := ExecutionStore{g.s.dbStore}

	e, err = db.GetExecution(e)
	if err != nil {
		return err
	}

	g.s.logger.Printf("Executing query: %v", e.Query)

	tx, err := g.s.dbData.BeginTxx(context.Background(), &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelDefault,
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec("SET TRANSACTION READ ONLY;")
	if err != nil {
		return err
	}

	rows, err := tx.Queryx(e.Query)

	if err != nil {
		return err
	}

	data := make([]map[string]interface{}, 0)

	defer rows.Close()
	for rows.Next() {
		entry := make(map[string]interface{})

		err := rows.MapScan(entry)
		if err != nil {
			return err
		}

		data = append(data, entry)
	}
	tx.Commit()

	dataBytes, err := getBytes(data)
	if err != nil {
		return err
	}

	return dest.SetBytes(dataBytes)
}

func (g ResultsCache) GetResults(key string, data interface{}) error {
	var res []byte
	if err := g.Get(nil, key, groupcache.AllocatingByteSliceSink(&res)); err != nil {
		return err
	}

	return getInterface(res, data)
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

func getInterface(bts []byte, data interface{}) error {
	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)

	return dec.Decode(data)
}
