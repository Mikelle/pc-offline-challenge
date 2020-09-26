package main

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/text/language"
)

type translatorWrapper struct {
	translator Translator
}

func NewTranslatorWrapper(t Translator) *translatorWrapper {
	return &translatorWrapper{
		translator: t,
	}
}

func (tw translatorWrapper) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	return tw.translateWithRetry(ctx, from, to, data)
}

func (tw translatorWrapper) translateWithRetry(ctx context.Context, from, to language.Tag, data string) (string, error) {
	var translation string
	operation := func() error {
		var err error
		translation, err = tw.Translate(ctx, from, to, data)
		return err
	}
	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	return translation, err
}


