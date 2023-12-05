package generator

import (
	"crypto/rand"
)

type Generator struct {
	length int
}

func NewGenerator(length int) *Generator {
	return &Generator{
		length: length,
	}
}

func (g *Generator) Generate() string {
	token := make([]byte, g.length)

	_, _ = rand.Read(token)

	return string(token)
}
