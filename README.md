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
export DB_PASSWORD=a26c375687
export DB_NAME=admin_tax
export JWT_SECRET=your-jwt-secret
```

The application reads `configs/config.yaml` for defaults but any environment variable above will override the values in the file.

To register, first request an OTP by POSTing to `/auth/send-otp` with `purpose` set to `verify_email` and the desired email. Then call `/auth/register` providing the email as `username`, the password, and the `otp_ref` and `otp_code` returned from the previous step. Only when the OTP is verified will the user account be created with the `user` role. Access to invoice endpoints requires either the `user` or `admin` role.
The `/me` endpoint returns the authenticated user's profile along with merchant details, including stores and whether the merchant is a person or company.

Run the project with:

```bash
go mod tidy
go run cmd/main.go
```

### Using Gmail OTP Service

To send OTP codes via Gmail, provide OAuth credentials and a token file
via the following environment variables or in `configs/config.yaml`:

```bash
export GMAIL_CREDENTIALS_FILE=/path/to/credentials.json
export GMAIL_TOKEN_FILE=/path/to/token.json
export GMAIL_FROM_EMAIL=you@example.com
```

When these values are set, the server will use Gmail to deliver OTP
codes. Otherwise an in-memory service is used.

### OTP Purposes

Requests to `/auth/send-otp` and `/auth/verify-otp` require a `purpose`
value to indicate why the OTP was issued. Valid values are:

- `verify_email` - send a code to verify a user's email address
- `reset_password` - send a code to allow password reset

Any other value will result in a `400 Bad Request` response.

Once the user receives a reset password code, submit it along with the
new password to `POST /auth/reset-password`.

### Using SMTP OTP Service

You can also send OTP codes through a generic SMTP server. Provide the
SMTP connection details via environment variables or `configs/config.yaml`:

```bash
export SMTP_HOST=smtp.example.com
export SMTP_PORT=587
export SMTP_USERNAME=user
export SMTP_PASSWORD=pass
export SMTP_FROM_EMAIL=noreply@example.com
```

If `SMTP_HOST` and `SMTP_FROM_EMAIL` are set, the application will use the
SMTP service for delivering OTP codes. Gmail settings take precedence
when provided; otherwise an in-memory service is used.
