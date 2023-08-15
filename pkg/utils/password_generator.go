package utils

import (
	"crypto/rand"
)

const (
	letters        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits         = "0123456789"
	specialChars   = "!@#$%^&*_-+"
	passwordLength = 10
)

func GenerateRandomPassword(length int) string {
	// สร้างเครื่องหมายระหว่างตัวอักษรและตัวเลขในรหัสผ่าน
	symbols := letters + digits + specialChars

	// สร้าง slice เพื่อเก็บรหัสผ่าน
	password := make([]byte, length)

	// สุ่มตัวเลขและตัวอักษรที่ใช้สร้างรหัสผ่าน
	rand.Read(password)

	// แปลงตัวเลขให้อยู่ในช่วง [0, len(symbols))
	for i := 0; i < length; i++ {
		password[i] = symbols[password[i]%byte(len(symbols))]
	}

	return string(password)
}
