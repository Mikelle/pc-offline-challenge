package main

import (
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

// Service is a Translator user.
type Service struct {
	translator Translator
}

func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)

	c := cache.New(5 * time.Minute, 10 * time.Minute)

	var rg singleflight.Group

	wt := NewTranslatorWrapper(t, c, &rg)
	return &Service{
		translator: wt,
	}
}
