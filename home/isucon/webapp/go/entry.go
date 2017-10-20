package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	keyEntryFreshness = "entryfreshness"
	keyEntryLength    = "entrylength"
)

func entryKey(k string) string {
	return fmt.Sprintf("entry:%s", k)
}

func saveEntry(e *Entry) {
	conn := pool.Get()
	defer conn.Close()
	conn.Do("HMSET", entryKey(e.Keyword), "description", e.Description)
	conn.Do("ZADD", keyEntryFreshness, time.Now().UTC().UnixNano(), e.Keyword)
	conn.Do("ZADD", keyEntryLength, len(e.Keyword), e.Keyword)
}

func deleteEntry(e *Entry) {
	conn := pool.Get()
	defer conn.Close()
	conn.Do("DEL", entryKey(e.Keyword))
	conn.Do("ZREM", keyEntryLength, e.Keyword)
	conn.Do("ZREM", keyEntryFreshness, e.Keyword)
}

func getEntryByKeyword(k string) *Entry {
	conn := pool.Get()
	defer conn.Close()
	m, err := redis.StringMap(conn.Do("HGETALL", entryKey(k)))
	panicIf(err)
	if m == nil {
		return nil
	}
	return &Entry{
		Description: m["description"],
	}
}

func getEntries(limit int, offset int) []*Entry {
	conn := pool.Get()
	defer conn.Close()
	keywords, err := redis.Strings(conn.Do("ZREVRANGE", keyEntryFreshness, offset, limit+offset))
	panicIf(err)
	entries := make([]*Entry, 0, len(keywords))
	for _, k := range keywords {
		m, err := redis.StringMap(conn.Do("HGETALL", entryKey(k)))
		panicIf(err)
		entries = append(entries, &Entry{
			Description: m["description"],
		})
	}
	return entries
}

func getEntryCount() int {
	conn := pool.Get()
	defer conn.Close()
	cnt, err := redis.Int(conn.Do("ZCARD", keyEntryFreshness))
	panicIf(err)
	return cnt
}

func getKeywordsOrderByLength() []string {
	conn := pool.Get()
	defer conn.Close()
	keywords, err := redis.Strings(conn.Do("ZREVRANGE", keyEntryLength, 0, -1))
	panicIf(err)
	return keywords
}
