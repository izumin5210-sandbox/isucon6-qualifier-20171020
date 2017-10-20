package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func starKey(k string) string {
	return fmt.Sprintf("star:%s", k)
}

func saveStar(s *Star) {
	conn := pool.Get()
	defer conn.Close()
	conn.Do("LPUSH", starKey(s.Keyword), s.UserName)
}

func getStarsByKeyword(k string) []*Star {
	conn := pool.Get()
	defer conn.Close()
	usernames, err := redis.Strings(conn.Do("LRANGE", starKey(k)))
	panicIf(err)
	stars := make([]*Star, 0, len(usernames))
	for _, username := range usernames {
		stars = append(stars, &Star{
			Keyword:  k,
			UserName: username,
		})
	}
	return stars
}
