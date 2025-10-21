# ---- Build Stage ----
# ใช้ Go เวอร์ชั่นล่าสุดเป็น base image สำหรับ build
FROM golang:1.24.2-alpine AS builder

# Install timezone data
RUN apk add --no-cache tzdata

# ตั้งค่า Working Directory ภายใน Image
WORKDIR /app

# Copy ไฟล์ที่ใช้จัดการ dependencies และดาวน์โหลด
COPY go.mod go.sum ./
RUN go mod download

# Copy โค้ดทั้งหมดของโปรเจกต์เข้ามา
COPY . .

# Build โปรแกรม Go ให้เป็นไฟล์ binary ชื่อ server
# CGO_ENABLED=0 ทำให้ได้ไฟล์แบบ static ที่ไม่ต้องพึ่งพา C libraries
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./main.go

# ---- Final Stage ----
# ใช้ base image ขนาดเล็กสำหรับ run application จริง
FROM alpine:latest

WORKDIR /app

# Copy เฉพาะไฟล์ binary ที่ build เสร็จแล้วจาก Stage แรก
COPY --from=builder /app/server .

# เปิด Port ที่แอปพลิเคชันของคุณรัน (จากไฟล์ .env คือ 8000)
EXPOSE 8000

# คำสั่งที่จะรันเมื่อ Container เริ่มทำงาน
CMD ["/app/server"]