package domain

import (
	"fmt"
	"os"
)

func NewWelcomeEmailBody(name string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4; }
        .container { max-width: 600px; margin: 0 auto; background: white; }
        .header { background: #2c3e50; color: white; padding: 30px; text-align: center; }
        .content { padding: 30px; }
        .welcome { font-size: 24px; color: #2c3e50; margin-bottom: 20px; }
        .text { font-size: 16px; line-height: 1.6; margin-bottom: 20px; }
        .button { display: inline-block; background: #3498db; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { background: #34495e; color: white; padding: 20px; text-align: center; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Foundation Name</h1>
        </div>
        <div class="content">
            <div class="welcome">Welcome %s!</div>
            <div class="text">
                Thank you for joining our foundation. Together we will make a difference.
            </div>
            <div class="text">
                You can follow our projects and see where your donations are being used.
            </div>
            <a href="#" class="button">View Projects</a>
            <div class="text">
                For questions: info@foundation.org
            </div>
        </div>
        <div class="footer">
            <p>© 2025 Foundation Name - All rights reserved</p>
        </div>
    </div>
</body>
</html>`, name)
}

func NewVerificationEmailLink(token string) string {

	return fmt.Sprintf("https://%s/verify-email?token=%s", os.Getenv("DOMAIN"), token)
}

func NewVerificationEmailBody(url string) string {
	return fmt.Sprintf("<p>Hello,</p><p>Please click the link below to verify your email address:</p><a href='%s'>Verify Email</a>", url)
}

var (
	WelcomeEmailSubject             = "Welcome to Advanced Todo API"
	SuccessfullyDeletedEmailSubject = "Your Account Has Been Successfully Deleted"
	VerificationEmailSubject        = "Please Verify Your Email Address"
	ForgotPasswordEmailSubject      = "Password Reset"
	EnglishSuccessfullyDeletedEmail = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Account Deletion Confirmation</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo {
            font-size: 24px;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .title {
            color: #e74c3c;
            font-size: 22px;
            font-weight: bold;
            margin-bottom: 20px;
        }
        .content {
            margin-bottom: 25px;
        }
        .info-box {
            background-color: #fff3cd;
            border: 1px solid #ffeaa7;
            border-radius: 5px;
            padding: 15px;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
            font-size: 14px;
        }
        .highlight {
            color: #e74c3c;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">{{FOUNDATION_NAME}}</div>
            <h1 class="title">Your Account Has Been Successfully Deleted</h1>
        </div>
        
        <div class="content">
            <p>Hello,</p>
            
            <p>Your account deletion request has been successfully processed. Your account and all personal information have been permanently removed from our system.</p>
            
            <div class="info-box">
                <h3>Important Information:</h3>
                <ul>
                    <li><span class="highlight">All your active donations have been cancelled</span></li>
                    <li>Your personal data has been deleted from our system</li>
                    <li>Your past donation records are kept anonymously</li>
                    <li>This action cannot be undone</li>
                </ul>
            </div>
            
            <p>If you wish to make donations in the future, you will need to create a new account.</p>
            
            <p>Thank you for your support until today. We wish you success in your charitable endeavors.</p>
        </div>
        
        <div class="footer">
            <p>This email was sent automatically. Please do not reply.</p>
            <p>{{FOUNDATION_NAME}} • {{FOUNDATION_ADDRESS}}</p>
            <p>{{FOUNDATION_PHONE}} • {{FOUNDATION_EMAIL}}</p>
        </div>
    </div>
</body>
</html>
        `
)

func NewForgotPasswordLink(token string) string {
	return fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("CLIENT_URL"), token)
}

func NewForgotPasswordEmail(url string) string {
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
