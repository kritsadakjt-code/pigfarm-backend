package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CheckJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึงค่าจะได้ Authorization : token
		authHeader := c.Get("Authorization")
		// check ถ้าไม่ขึ้นต้นด้วย Bearer
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "Missing or invalid Token"})
		}
		// ตัด Bearer ข้างหน้าออก
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		// check token ที่ ผูกกับ JWT_SECRET ตอน signing
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		// check เกิด error หรือ token หมดอายุ
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Missing or invalid Token"})
		}

		// ดึงข้อมูลที่เก็บไว้ใน claims
		claims := token.Claims.(jwt.MapClaims)
		//ดึงค่าจาก claims มาเก็บไว้ใน fiber เพื่อเอาไปใช้ต่อไป
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])
		return c.Next()

	}
}

func IsAdminOnly(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "owner" {
		return c.Status(403).JSON(fiber.Map{"error": "Only admin can access"})
	}
	return c.Next()
}
