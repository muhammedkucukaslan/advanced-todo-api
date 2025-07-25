package domain

import (
	"fmt"
	"os"
)

func NewForgotPasswordLink(token string) string {
	return fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("CLIENT_URL"), token)
}

func NewForgotPasswordSubject(language string) string {
	switch language {
	case "tr":
		return "Şifre Sıfırlama"
	case "en":
		return "Password Reset"
	case "ar":
		return "إعادة تعيين كلمة المرور"
	default:
		return ""
	}
}

func NewForgotPasswordEmail(url, language string) string {
	switch language {
	case "tr":
		return NewTurkishForgotPasswordEmail(url)
	case "en":
		return NewEnglishForgotPasswordEmail(url)
	case "ar":
		return NewArabicForgotPasswordEmail(url)
	default:
		return ""
	}
}

func NewTurkishForgotPasswordEmail(url string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Şifre Sıfırlama</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            background-color: #f4f4f4;
            padding: 20px;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background: #fff;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
        }
        p {
            color: #555;
        }
        a {
            color: #e74c3c;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Şifre Sıfırlama Talebi</h1>
        <p>Merhaba,</p>
        <p>Şifrenizi sıfırlamak için lütfen aşağıdaki bağlantıya tıklayın:</p>
        <p><a href="%s">Şifreyi Sıfırla</a></p>
        <p>Bu bağlantı, talep edildiği tarihten itibaren 30 dakika boyunca geçerlidir.</p>
        <p>Teşekkürler,</p>
        <p>Ekibimiz</p>
    </div>
</body>
</html>
`, url)
}

func NewArabicForgotPasswordEmail(url string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="ar" dir="rtl">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>إعادة تعيين كلمة المرور</title>
	<style>
	body {
		font-family: 'Segoe UI', Tahoma, Arial, sans-serif;
		line-height: 1.6;
		background-color: #f4f4f4;
		padding: 20px;
		direction: rtl;
	}
	.container {
		max-width: 600px;
		margin: 0 auto;
		background: #fff;
		padding: 20px;
		border-radius: 5px;
		box-shadow: 0 2px 5px rgba(0,0,0,0.1);
	}
	h1 {
		color: #333;
	}
	p {
		color: #555;
	}
	a {
		color: #e74c3c;
		text-decoration: none;
		font-weight: bold;
	}
	.btn {
		display: inline-block;
		background-color: #e74c3c;
		color: white;
		padding: 10px 20px;
		border-radius: 5px;
		text-decoration: none;
		margin: 10px 0;
	}
	</style>
</head>
<body>
	<div class="container">
		<h1>طلب إعادة تعيين كلمة المرور</h1>
		<p>مرحباً،</p>
		<p>لإعادة تعيين كلمة المرور الخاصة بك، يرجى النقر على الرابط التالي:</p>
		<p><a href="%s" class="btn">إعادة تعيين كلمة المرور</a></p>
		<p>أو يمكنك نسخ الرابط التالي ولصقه في المتصفح:</p>
		<p><a href="%s">%s</a></p>
		<p>هذا الرابط صالح لمدة 30 دقيقة من تاريخ الطلب.</p>
		<p>إذا لم تطلب إعادة تعيين كلمة المرور، يرجى تجاهل هذه الرسالة.</p>
		<p>شكراً لك،</p>
		<p>فريقنا</p>
	</div>
</body>
</html>
`, url, url, url)
}

// English Password Reset Email Template
func NewEnglishForgotPasswordEmail(url string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Password Reset</title>
	<style>
	body {
		font-family: Arial, sans-serif;
		line-height: 1.6;
		background-color: #f4f4f4;
		padding: 20px;
		margin: 0;
	}
	.container {
		max-width: 600px;
		margin: 0 auto;
		background: #fff;
		padding: 20px;
		border-radius: 5px;
		box-shadow: 0 2px 5px rgba(0,0,0,0.1);
	}
	h1 {
		color: #333;
		margin-bottom: 20px;
	}
	p {
		color: #555;
		margin-bottom: 15px;
	}
	a {
		color: #e74c3c;
		text-decoration: none;
	}
	.btn {
		display: inline-block;
		background-color: #e74c3c;
		color: white;
		padding: 12px 25px;
		border-radius: 5px;
		text-decoration: none;
		margin: 15px 0;
		font-weight: bold;
	}
	.btn:hover {
		background-color: #c0392b;
	}
	.footer {
		margin-top: 30px;
		padding-top: 20px;
		border-top: 1px solid #eee;
		font-size: 12px;
		color: #999;
	}
	</style>
</head>
<body>
	<div class="container">
		<h1>Password Reset Request</h1>
		<p>Hello,</p>
		<p>We received a request to reset your password. Click the button below to create a new password:</p>
		<p><a href="%s" class="btn">Reset Password</a></p>
		<p>Or copy and paste this link into your browser:</p>
		<p><a href="%s">%s</a></p>
		<p>This link will expire in 30 minutes from the time of request.</p>
		<p>If you didn't request a password reset, please ignore this email. Your password will remain unchanged.</p>
		<p>Best regards,</p>
		<p>Our Team</p>
		<div class="footer">
			<p>This is an automated email. Please do not reply to this message.</p>
		</div>
	</div>
</body>
</html>
`, url, url, url)
}
