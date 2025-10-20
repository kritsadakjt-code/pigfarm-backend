package utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SentResetPasswordFromEmail(toEmail, resetURL string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT in .env file: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "คำขอรีเซ็ตรหัสผ่านสำหรับ Pig Farm")

	emailBody := fmt.Sprintf(`
    <h3>คำขอรีเซ็ตรหัสผ่าน</h3>
    <p>กรุณาคลิกที่ลิงก์ด้านล่างเพื่อตั้งรหัสผ่านใหม่ ลิงก์นี้จะหมดอายุใน 30 นาที:</p>
    <p><a href="%s"><strong>ตั้งรหัสผ่านใหม่</strong></a></p>
    <p>หากคุณไม่ได้ร้องขอ กรุณาไม่ต้องดำเนินการใดๆ</p>
`, resetURL)

	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	fmt.Printf("Password reset email sent successfully to %s\n", toEmail)
	return nil

}

func SendEmailVerification(toEmail, verificationLink string) error {

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// แปลง Port จาก string เป็น integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT in .env file: %w", err)
	}

	// 2. สร้างเนื้อหาอีเมล
	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "ยืนยันบัญชีของคุณสำหรับ Pig Farm Management")

	// เนื้อหาอีเมลแบบ HTML ที่มีปุ่ม Call-to-action ชัดเจน
	emailBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			.button {
				background-color: #1ed837ff; 
				border: none;
				color: white;
				padding: 15px 32px;
				text-align: center;
				text-decoration: none;
				display: inline-block;
				font-size: 16px;
				margin: 4px 2px;
				cursor: pointer;
				border-radius: 8px;
			}
		</style>
	</head>
	<body>
		<h2>สวัสดีครับ,</h2>
		<p>กรุณาคลิกที่ปุ่มด้านล่างเพื่อยืนยันอีเมลของคุณ:</p>
		<a href="%s" class="button">ยืนยันบัญชี</a>
		<p>หากคุณไม่ได้ร้องขอ กรุณาไม่ต้องดำเนินการใดๆ</p>
		<p>ขอบคุณครับ,<br>ทีมงาน Pig Farm</p>
	</body>
	</html>
	`, verificationLink)

	m.SetBody("text/html", emailBody)

	// 3. ตั้งค่าการเชื่อมต่อกับ SMTP Server (ในที่นี้คือ Gmail)
	// Dialer จะใช้ "App Password" ที่เราสร้างไว้ในการยืนยันตัวตน
	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)

	// 4. ส่งอีเมล
	if err := d.DialAndSend(m); err != nil {
		// ถ้าส่งไม่สำเร็จ จะคืนค่า error กลับไป
		return fmt.Errorf("could not send verification email: %w", err)
	}

	fmt.Printf("Verification email sent successfully to %s\n", toEmail)
	return nil
}
