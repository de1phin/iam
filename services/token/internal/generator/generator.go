package generator

import (
	"math/rand"
)

const charset = "ABCDEFGHIJKLMOPQRSTUVWXZabcedefghijklmnopqrstuvwxz0123456789+_-="

type Generator struct {
	length int
}

func NewGenerator(length int) *Generator {
	return &Generator{
		length: length,
	}
}

func (g *Generator) Generate() string {
	id := ""
	for i := 0; i < g.length; i++ {
		id += string(charset[rand.Intn(len(charset))])
	}
	return id
}
