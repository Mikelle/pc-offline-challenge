package main

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
	"golang.org/x/text/language"
)

type translatorWrapper struct {
	translator   Translator
	cache        *cache.Cache
	requestGroup *singleflight.Group
}

func NewTranslatorWrapper(t Translator, c *cache.Cache, requestGroup *singleflight.Group) *translatorWrapper {
	return &translatorWrapper{
		translator:   t,
		cache:        c,
		requestGroup: requestGroup,
	}
}

// checking value from cache, if not exist send request
// deduplicate simultaneous queries for the same parameters
// also sending with exponential backoff
// if request succeeded, then save result to cache and return
func (tw translatorWrapper) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	hash := fmt.Sprintf("%v#%v#%s", from, to, data)
	tr, found := tw.cache.Get(hash)
	if found {
		return tr.(string), nil
	}

	v, err, _ := tw.requestGroup.Do(hash, func() (interface{}, error) {
		return tw.translateWithRetry(ctx, from, to, data)
	})

	if err != nil {
		return "", err
	}

	translation, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("failed to cast translation to string")
	}

	tw.cache.Set(hash, translation, cache.DefaultExpiration)
	return translation, nil
}

func (tw translatorWrapper) translateWithRetry(ctx context.Context, from, to language.Tag, data string) (string, error) {
	var translation string
	operation := func() error {
		var err error
		translation, err = tw.translator.Translate(ctx, from, to, data)
		return err
	}
	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		return "", err
	}

	return translation, nil
}
