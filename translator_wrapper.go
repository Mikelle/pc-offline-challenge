package main

import (
	"context"
	"fmt"

	"github.com/patrickmn/go-cache"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/text/language"
)

type translatorWrapper struct {
	translator Translator
	cache      *cache.Cache
}

func NewTranslatorWrapper(t Translator, c *cache.Cache) *translatorWrapper {
	return &translatorWrapper{
		translator: t,
		cache: c,
	}
}

func (tw translatorWrapper) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	hash := fmt.Sprintf("%v#%v#%s", from, to, data)
	tr, found := tw.cache.Get(hash)
	if found {
		return tr.(string), nil
	}

	translation, err := tw.translateWithRetry(ctx, from, to, data)
	if err != nil {
		return "", err
	}

	tw.cache.Set(hash, translation, cache.DefaultExpiration)
	return translation, err
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


