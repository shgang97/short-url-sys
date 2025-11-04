package link

import (
	"crypto/rand"
	"generate-service/internal/pkg/base62"
	"generate-service/internal/pkg/errors"
	"math/big"
	"regexp"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// ShortCodeGenerator 短码生成器
type ShortCodeGenerator struct {
	customCodeRegex *regexp.Regexp
	minCustomLength int
	maxCustomLength int
}

// NewShortCodeGenerator 创建短码生成器
func NewShortCodeGenerator() *ShortCodeGenerator {
	return &ShortCodeGenerator{
		customCodeRegex: regexp.MustCompile("^[A-Za-z0-9_-]+$"),
		minCustomLength: 3,
		maxCustomLength: 20,
	}
}

// GenerateFromID 从ID生成短码
func (g *ShortCodeGenerator) GenerateFromID(id uint64) string {
	return base62.Encode(id)
}

// ValidateCustomCode 验证自定义短码
func (g *ShortCodeGenerator) ValidateCustomCode(code string) error {
	if len(code) < g.minCustomLength || len(code) > g.maxCustomLength {
		return errors.NewBusinessError("custom code length out of range")
	}
	if !g.customCodeRegex.MatchString(code) {
		return errors.NewBusinessError("custom code can only contain letters, numbers, hyphens and underscores")
	}
	return nil
}

// GenerateRandomCode 生成随机短码（用于自定义短码冲突时）
func (g *ShortCodeGenerator) GenerateRandomCode(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
