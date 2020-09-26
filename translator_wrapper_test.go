package main

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestTranslatorWrapper(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	t.Run("test translate with cache", func(t *testing.T) {
		ctx := context.Background()
		from := language.English
		to := language.Japanese
		data := "test"

		tr := newRandomTranslator(
			100*time.Millisecond,
			500*time.Millisecond,
			0.1,
		)

		c := cache.New(5 * time.Minute, 10 * time.Minute)

		translator := NewTranslatorWrapper(tr, c)
		res1, err := translator.Translate(ctx, from, to, data)
		require.Nil(t, err)
		res2, err := translator.Translate(ctx, from, to, data)
		require.Nil(t, err)
		require.Equal(t, res1, res2)
	})
}
