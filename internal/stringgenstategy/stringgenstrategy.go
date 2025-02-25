package stringgenstategy

import "github.com/physicist2018/url-shortener-go/internal/ports/randomstring"

// Контекст, который использует стратегию генерации строк
type StringGeneratorContext struct {
	strategy randomstring.StringGenerator
}

// Установка стратегии
func (c *StringGeneratorContext) SetStrategy(strategy randomstring.StringGenerator) {
	c.strategy = strategy
}

// Генерация строки с использованием текущей стратегии
func (c *StringGeneratorContext) GenerateString() string {
	return c.strategy.Generate()
}
