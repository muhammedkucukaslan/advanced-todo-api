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
