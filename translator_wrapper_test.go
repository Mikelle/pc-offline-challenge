package main

import (
	"context"
	"golang.org/x/sync/singleflight"
	"math/rand"
	"sync"
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

		var rg singleflight.Group
		translator := NewTranslatorWrapper(tr, c, &rg)
		res1, err := translator.Translate(ctx, from, to, data)
		require.Nil(t, err)
		res2, err := translator.Translate(ctx, from, to, data)
		require.Nil(t, err)
		require.Equal(t, res1, res2)
	})

	t.Run("test concurrent translate", func(t *testing.T) {
		ctx := context.Background()
		from := language.English
		to := language.Japanese
		data := "test"

		var wg sync.WaitGroup
		var res [10]string

		tr := newRandomTranslator(
			100*time.Millisecond,
			500*time.Millisecond,
			0,
		)

		c := cache.New(5 * time.Minute, 10 * time.Minute)

		var rg singleflight.Group
		translator := NewTranslatorWrapper(tr, c, &rg)

		for i := 0; i < 10; i++ {
			wg.Add(1)
			i := i
			go func(ctx context.Context, translator Translator, res [10]string) {
				defer wg.Done()
				r, err := tr.Translate(ctx, from, to, data)
				require.Nil(t, err)
				res[i] = r
			} (ctx, translator, res)
		}

		for i := 0; i < 10; i++ {
			require.Equal(t, res[0], res[i])
		}
	})
}
