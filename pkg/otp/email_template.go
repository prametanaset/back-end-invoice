package otp

import "fmt"

func buildOTPEmail(to, code string) string {
	return fmt.Sprintf("To: %s\r\nSubject: Your OTP Code\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"+
		`<!DOCTYPE html>
<html>
<head>
<style>
body { font-family: Arial, sans-serif; background-color: #f7f7f7; padding: 20px; }
.container { max-width: 600px; margin: auto; background-color: #ffffff; padding: 20px; border-radius: 4px; text-align: center; }
.code { font-size: 32px; font-weight: bold; letter-spacing: 4px; }
</style>
</head>
<body>
  <div class="container">
    <p>Your OTP code is:</p>
    <p class="code">%s</p>
  </div>
</body>
</html>`, to, code)
}
