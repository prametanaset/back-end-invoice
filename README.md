# Invoice Project

This is a sample invoice system built with GoFiber, GORM and PostgreSQL.

## Requirements

- Go 1.20+
- PostgreSQL
- `go install github.com/gofiber/fiber/v2@latest`
- `go install gorm.io/gorm@latest`
- `go install gorm.io/driver/postgres@latest`
- (optional) `go install github.com/google/wire/cmd/wire@latest`

## Usage

Set the following environment variables to configure the database and JWT secret:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your-password
export DB_NAME=invoice_db
export JWT_SECRET=your-jwt-secret
export JWT_ISSUER=invoice_service
export JWT_AUDIENCE=invoice_api
```

The application reads `configs/config.yaml` for defaults but any environment variable above will override the values in the file.

Users registered via `/auth/register` are created with the `user` role by default. Access to invoice endpoints now requires either the `user` or `admin` role.
JWT tokens now include standard `iss` and `aud` claims derived from `JWT_ISSUER` and `JWT_AUDIENCE`.

Run the project with:

```bash
go mod tidy
go run cmd/main.go
```
