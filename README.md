# Invoice Project

โปรเจกต์ตัวอย่างสร้างระบบ Invoice ด้วย GoFiber + GORM + PostgreSQL ครับ

## Requirements

- Go 1.20+
- PostgreSQL
- `go install github.com/gofiber/fiber/v2@latest`
- `go install gorm.io/gorm@latest`
- `go install gorm.io/driver/postgres@latest`
- (ถ้าใช้ Wire) `go install github.com/google/wire/cmd/wire@latest`

## วิธีใช้งาน

1. สร้างฐานข้อมูล PostgreSQL ชื่อ `invoice_db` แล้วตั้งค่า user/password ให้ตรงกับ `configs/config.yaml`
2. รันคำสั่ง:
   ```bash
   go mod tidy
   go run cmd/main.go
# back-end-invoice
