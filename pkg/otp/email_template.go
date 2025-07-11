package otp

import "fmt"

func buildOTPEmail(to, code string) string {
	return fmt.Sprintf("To: %s\r\nSubject: OTP for E-mail Verification on ScaleTax by Sunscaleup\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"+
		`<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body {
      margin: 0;
      padding: 0;
      background-color: #f4f4f4;
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    }
    .container {
      max-width: 600px;
      margin: 40px auto;
      background-color: #ffffff;
      padding: 30px;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }
    .header {
      text-align: center;
      padding-bottom: 20px;
    }
    .header h1 {
      margin: 0;
      font-size: 24px;
      color: #af38ff;
    }
    .content {
      text-align: center;
      font-size: 16px;
      color: #555555;
    }
    .otp-code {
      margin: 20px 0;
      font-size: 36px;
      font-weight: bold;
      color: #2b2b2b;
      letter-spacing: 6px;
      background-color:rgba(175, 56, 255, 0.22);
      display: inline-block;
      padding: 10px 20px;
      border-radius: 6px;
    }
    .footer {
      margin-top: 30px;
      font-size: 13px;
      color: #999999;
      text-align: center;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>รหัสยืนยัน (OTP)</h1>
    </div>
    <div class="content">
      <p>กรุณาใช้รหัสนี้เพื่อดำเนินการยืนยันตัวตนให้เสร็จสิ้น:</p>
      <div class="otp-code">%s</div>
      <p>รหัสยืนยันจะมีอายุเพียง 5 นาที ห้ามเปิดเผยรหัสนี้กับบุคคลอื่น</p>
    </div>
    <div class="footer">
      <p>&copy; 2025 sunscaleup Ltd. All rights reserved.</p>
    </div>
  </div>
</body>
</html>`, to, code)
}
