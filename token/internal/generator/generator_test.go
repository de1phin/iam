package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	t.Parallel()

	var (
		g                 = NewGenerator(512)
		amountGenerations = 1_000_000
		set               = make(map[string]bool, amountGenerations)
	)

	for i := 0; i < amountGenerations; i++ {
		token := g.Generate()

		if set[token] {
			assert.Fail(t, "failed unique")
		}

		set[token] = true
	}
}
