package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/rs/xid"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

const expiration = 24 * time.Hour

var notfound = errors.New("not found")

type saveable struct {
	Content string
}

func key() string {
	return xid.New().String()
}

func saveContent(r *http.Request, c string) (string, error) {
	ctx := appengine.NewContext(r)
	item := memcache.Item{
		Key:        key(),
		Value:      []byte(c),
		Expiration: expiration,
	}

	if err := memcache.Set(ctx, &item); err != nil {
		return "", err
	}

	return item.Key, nil
}

func getContent(r *http.Request, key string) (string, error) {
	ctx := appengine.NewContext(r)

	item, err := memcache.Get(ctx, key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return "", notfound
		}

		return "", err
	}

	return string(item.Value), nil
}
