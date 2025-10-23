package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// ฟังก์ชันกลางสำหรับส่งอีเมลผ่าน Mailjet API
func sendEmailWithMailjet(toEmail, subject, htmlContent string) error {
	apiKey := os.Getenv("MAILJET_API_KEY")
	apiSecret := os.Getenv("MAILJET_API_SECRET")
	senderEmail := os.Getenv("MAILJET_SENDER_EMAIL")
	senderName := os.Getenv("MAILJET_SENDER_NAME")

	if apiKey == "" || apiSecret == "" || senderEmail == "" {
		return fmt.Errorf("MAILJET_API_KEY, MAILJET_API_SECRET or MAILJET_SENDER_EMAIL not set")
	}

	payload := map[string]interface{}{
		"Messages": []map[string]interface{}{
			{
				"From": map[string]string{
					"Email": senderEmail,
					"Name":  senderName,
				},
				"To": []map[string]string{
					{"Email": toEmail},
				},
				"Subject":  subject,
				"HTMLPart": htmlContent,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.mailjet.com/v3.1/send", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mailjet API returned status: %s", resp.Status)
	}

	fmt.Printf("✅ Email sent successfully to %s via Mailjet\n", toEmail)
	return nil
}

// --- ฟังก์ชันสำหรับ Reset Password ---
func SentResetPasswordFromEmail(toEmail, resetURL string) error {
	subject := "คำขอรีเซ็ตรหัสผ่านสำหรับ Pig Farm"
	emailBody := fmt.Sprintf(`
	<h3>คำขอรีเซ็ตรหัสผ่าน</h3>
	<p>กรุณาคลิกลิงก์ด้านล่างเพื่อตั้งรหัสผ่านใหม่ (ภายใน 30 นาที):</p>
	<p><a href="%s" style="background-color:#4CAF50;color:white;padding:14px 25px;text-align:center;text-decoration:none;display:inline-block;border-radius:8px;">ตั้งรหัสผ่านใหม่</a></p>
	<p>หากคุณไม่ได้ร้องขอ กรุณาไม่ต้องดำเนินการใดๆ</p>
	`, resetURL)

	return sendEmailWithMailjet(toEmail, subject, emailBody)
}

// --- ฟังก์ชันสำหรับ Email Verification ---
func SendEmailVerification(toEmail, verificationLink string) error {
	subject := "ยืนยันบัญชีของคุณสำหรับ Pig Farm Management"
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
		<p>กรุณาคลิกปุ่มด้านล่างเพื่อยืนยันอีเมลของคุณ:</p>
		<a href="%s" class="button">ยืนยันบัญชี</a>
		<p>หากคุณไม่ได้ร้องขอ กรุณาไม่ต้องดำเนินการใดๆ</p>
		<p>ขอบคุณครับ,<br>ทีมงาน Pig Farm</p>
	</body>
	</html>
	`, verificationLink)

	return sendEmailWithMailjet(toEmail, subject, emailBody)
}

// package utils

// import (
// 	"fmt"
// 	"os"

// 	postmark "github.com/keighl/postmark"
// )

// // ฟังก์ชันกลางสำหรับส่งอีเมลผ่าน Postmark
// func sendEmailWithPostmark(toEmail, subject, htmlContent string) error {

// 	postmarkServerToken := os.Getenv("POSTMARK_SERVER_TOKEN")
// 	senderSignatureEmail := os.Getenv("POSTMARK_SENDER_EMAIL")
// 	fmt.Println(postmarkServerToken, "TEST")

// 	if postmarkServerToken == "" || senderSignatureEmail == "" {
// 		return fmt.Errorf("POSTMARK_SERVER_TOKEN or POSTMARK_SENDER_EMAIL not set")
// 	}

// 	client := postmark.NewClient(postmarkServerToken, "")

// 	email := postmark.Email{
// 		From:       senderSignatureEmail,
// 		To:         toEmail,
// 		Subject:    subject,
// 		HtmlBody:   htmlContent, // แก้จาก HTMLBody เป็น HtmlBody
// 		TrackOpens: true,
// 	}

// 	// SDK version นี้ไม่รองรับ context
// 	_, err := client.SendEmail(email)
// 	if err != nil {
// 		fmt.Printf("Error sending email via Postmark: %v\n", err)
// 		return fmt.Errorf("could not send email via Postmark: %w", err)
// 	}

// 	fmt.Printf("Email sent successfully to %s via Postmark\n", toEmail)
// 	return nil
// }

// // --- ฟังก์ชันตัวอย่างสำหรับ Reset Password ---
// func SentResetPasswordFromEmail(toEmail, resetURL string) error {
// 	subject := "คำขอรีเซ็ตรหัสผ่านสำหรับ Pig Farm"
// 	emailBody := fmt.Sprintf(`
// 	<h3>คำขอรีเซ็ตรหัสผ่าน</h3>
// 	<p>กรุณาคลิกที่ลิงก์ด้านล่างเพื่อตั้งรหัสผ่านใหม่ ลิงก์นี้จะหมดอายุใน 30 นาที:</p>
// 	<p><a href="%s" style="background-color: #4CAF50; color: white; padding: 14px 25px; text-align: center; text-decoration: none; display: inline-block; border-radius: 8px;"><strong>ตั้งรหัสผ่านใหม่</strong></a></p>
// 	<p>หากคุณไม่ได้ร้องขอ กรุณาไม่ต้องดำเนินการใดๆ</p>
// 	`, resetURL)

// 	return sendEmailWithPostmark(toEmail, subject, emailBody)
// }

// // --- ฟังก์ชันตัวอย่างสำหรับ Email Verification ---
// func SendEmailVerification(toEmail, verificationLink string) error {
// 	subject := "ยืนยันบัญชีของคุณสำหรับ Pig Farm Management"
// 	emailBody := fmt.Sprintf(`
// 	<!DOCTYPE html>
// 	<html>
// 	<head>
// 		<style>
// 			.button {
// 				background-color: #1ed837ff;
// 				border: none;
// 				color: white;
// 				padding: 15px 32px;
// 				text-align: center;
// 				text-decoration: none;
// 				display: inline-block;
// 				font-size: 16px;
// 				margin: 4px 2px;
// 				cursor: pointer;
// 				border-radius: 8px;
// 			}
// 		</style>
// 	</head>
// 	<body>
// 		<h2>สวัสดีครับ,</h2>
// 		<p>กรุณาคลิกที่ปุ่มด้านล่างเพื่อยืนยันอีเมลของคุณ:</p>
// 		<a href="%s" class="button">ยืนยันบัญชี</a>
// 		<p>หากคุณไม่ได้ร้องขอ กรุณาไม่ต้องดำเนินการใดๆ</p>
// 		<p>ขอบคุณครับ,<br>ทีมงาน Pig Farm</p>
// 	</body>
// 	</html>
// 	`, verificationLink)

// 	return sendEmailWithPostmark(toEmail, subject, emailBody)
// }

// package utils

// import (
// 	"fmt"
// 	"os"
// 	"strconv"

// 	"gopkg.in/gomail.v2"
// )

// func SentResetPasswordFromEmail(toEmail, resetURL string) error {
// 	smtpHost := os.Getenv("SMTP_HOST")
// 	smtpPortStr := os.Getenv("SMTP_PORT")
// 	smtpEmail := os.Getenv("SMTP_EMAIL")
// 	smtpPassword := os.Getenv("SMTP_PASSWORD")

// 	smtpPort, err := strconv.Atoi(smtpPortStr)
// 	if err != nil {
// 		return fmt.Errorf("invalid SMTP_PORT in .env file: %w", err)
// 	}

// 	m := gomail.NewMessage()
// 	m.SetHeader("From", smtpEmail)
// 	m.SetHeader("To", toEmail)
// 	m.SetHeader("Subject", "คำขอรีเซ็ตรหัสผ่านสำหรับ Pig Farm")

// 	emailBody := fmt.Sprintf(`
//     <h3>คำขอรีเซ็ตรหัสผ่าน</h3>
//     <p>กรุณาคลิกที่ลิงก์ด้านล่างเพื่อตั้งรหัสผ่านใหม่ ลิงก์นี้จะหมดอายุใน 30 นาที:</p>
//     <p><a href="%s"><strong>ตั้งรหัสผ่านใหม่</strong></a></p>
//     <p>หากคุณไม่ได้ร้องขอ กรุณาไม่ต้องดำเนินการใดๆ</p>
// `, resetURL)

// 	m.SetBody("text/html", emailBody)

// 	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)

// 	if err := d.DialAndSend(m); err != nil {
// 		return fmt.Errorf("could not send email: %w", err)
// 	}

// 	fmt.Printf("Password reset email sent successfully to %s\n", toEmail)
// 	return nil

// }

// func SendEmailVerification(toEmail, verificationLink string) error {

// 	smtpHost := os.Getenv("SMTP_HOST")
// 	smtpPortStr := os.Getenv("SMTP_PORT")
// 	smtpEmail := os.Getenv("SMTP_EMAIL")
// 	smtpPassword := os.Getenv("SMTP_PASSWORD")

// 	// แปลง Port จาก string เป็น integer
// 	smtpPort, err := strconv.Atoi(smtpPortStr)
// 	if err != nil {
// 		return fmt.Errorf("invalid SMTP_PORT in .env file: %w", err)
// 	}

// 	// 2. สร้างเนื้อหาอีเมล
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", smtpEmail)
// 	m.SetHeader("To", toEmail)
// 	m.SetHeader("Subject", "ยืนยันบัญชีของคุณสำหรับ Pig Farm Management")

// 	// เนื้อหาอีเมลแบบ HTML ที่มีปุ่ม Call-to-action ชัดเจน
// 	emailBody := fmt.Sprintf(`
// 	<!DOCTYPE html>
// 	<html>
// 	<head>
// 		<style>
// 			.button {
// 				background-color: #1ed837ff;
// 				border: none;
// 				color: white;
// 				padding: 15px 32px;
// 				text-align: center;
// 				text-decoration: none;
// 				display: inline-block;
// 				font-size: 16px;
// 				margin: 4px 2px;
// 				cursor: pointer;
// 				border-radius: 8px;
// 			}
// 		</style>
// 	</head>
// 	<body>
// 		<h2>สวัสดีครับ,</h2>
// 		<p>กรุณาคลิกที่ปุ่มด้านล่างเพื่อยืนยันอีเมลของคุณ:</p>
// 		<a href="%s" class="button">ยืนยันบัญชี</a>
// 		<p>หากคุณไม่ได้ร้องขอ กรุณาไม่ต้องดำเนินการใดๆ</p>
// 		<p>ขอบคุณครับ,<br>ทีมงาน Pig Farm</p>
// 	</body>
// 	</html>
// 	`, verificationLink)

// 	m.SetBody("text/html", emailBody)

// 	// 3. ตั้งค่าการเชื่อมต่อกับ SMTP Server (ในที่นี้คือ Gmail)
// 	// Dialer จะใช้ "App Password" ที่เราสร้างไว้ในการยืนยันตัวตน
// 	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)

// 	// 4. ส่งอีเมล
// 	if err := d.DialAndSend(m); err != nil {
// 		// ถ้าส่งไม่สำเร็จ จะคืนค่า error กลับไป
// 		return fmt.Errorf("could not send verification email: %w", err)
// 	}

// 	fmt.Printf("Verification email sent successfully to %s\n", toEmail)
// 	return nil
// }
