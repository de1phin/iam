package generator

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	t.Parallel()

	var (
		g                 = NewGenerator(512)
		amountGenerations = 1_000_000
		amountTests       = 5
		globalSet         = make(map[string]bool, amountGenerations*amountTests)
		mu                = &sync.Mutex{}
		wg                = &sync.WaitGroup{}
	)

	for i := 0; i < amountTests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var localSet = make(map[string]bool, amountGenerations)
			for j := 0; j < amountGenerations; j++ {
				token := g.Generate()

				if localSet[token] {
					assert.Fail(t, "failed unique")
				}

				localSet[token] = true
			}

			mu.Lock()
			for token := range localSet {
				if _, ok := globalSet[token]; ok {
					assert.Fail(t, "failed unique")
				}
				globalSet[token] = true
			}
			mu.Unlock()
		}()
	}

	wg.Wait()
}
