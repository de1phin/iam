package generator

import (
	"crypto/rand"
	"encoding/hex"
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
	bytes := make([]byte, g.length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
