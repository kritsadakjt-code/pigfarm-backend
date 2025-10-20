package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// รับ จํานวน byte ที่ส่งมาคือ 32
func GenerateRandomToken(n int) string {
	// create slice 32 ช่อง
	b := make([]byte, n)
	// add ข้อมูลสุ่มเข้าไป
	rand.Read(b)
	// เเปลงเป็น string
	return hex.EncodeToString(b)
}
