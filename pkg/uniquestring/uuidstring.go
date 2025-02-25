package uniquestring

import "github.com/google/uuid"

// Стратегия для генерации строки с использованием UUID
type UUIDString struct{}

func NewUUIDString() *UUIDString {
	return &UUIDString{}
}

// Реализация метода интерфейса для генерации UUID строки
func (u *UUIDString) Generate() string {
	return uuid.New().String()
}
