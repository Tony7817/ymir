package vars

import "fmt"

// captcha
const CaptchaCodeAnyWay = "293049"
const CaptchaMaxSendTimesPerDay int64 = 12

// template sms send api
const SmsUrl = "https://api-v4.mysubmail.com/internationalsms/xsend"

func GetCaptchaEmailLastRequestTimeKey(email string) string {
	return "cache:captcha:email:" + email + ":last_request_time"
}
func GetCaptchaEmailSendTimesKey(email string) string {
	return "cache:captcha:email:" + email + ":sendtimes"
}

func GetCaptchaPhonenumberSendTimesKey(phonenumber string) string {
	return "cache:captcha:phonenumber:" + phonenumber + ":sendtimes"
}

func GetCaptchaPhonenumberLastRequestTimeKey(phonenumber string) string {
	return "cache:captcha:phonenumber:" + phonenumber + ":last_request_time"
}

// sql

// cache
const CacheExpireIn5m = 5 * 60
const CacheExpireIn1d = 60 * 60 * 24
const CacheExpireIn1m = 60

func GetCaptchaEmailCacheKey(email string) string {
	return "cache:capcha:email:" + email
}

func GetCaptchaPhonenumberCacheKey(phonenumber string) string {
	return "cache:capcha:phonenumber:" + phonenumber
}

// site
const ResetPasswordSite = "www.lureros.com/user/forgetpass/reset"

// email
// captcha
var EmailNoReplySenderName = "noreply@mail.lureros.com"
var EmailAlias = "Miss Lover"
var EmailCaptchaSubJect = "Verify your Email"

func GetCaptchaEmailTemplate(captcha string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
        .header {
            color: aliceblue;
            height: 50px;
            background: linear-gradient(45deg, #cb2a45, #2a4a82);
            text-align: center;
            margin-bottom: 20px;
            display: flex;
            justify-content: center;
            align-items: center;
        }
        .header h1 {
            font-size: 24px;
            margin: 0;
        }
        .title {
            font-size: 20px;
            margin: 20px 0;
            text-align: center;
        }
        .content {
            font-size: 16px;
            line-height: 1.6;
            text-align: center;
        }
        .code {
            font-size: 28px;
            font-weight: bold;
            color: #CB2A45;
            margin: 20px 0;
        }
        .instructions {
            margin-top: 20px;
            font-size: 14px;
            line-height: 1.6;
        }
        .instructions a {
            color: #CB2A45;
        }
        .gradient-linear {
            /* background-color: #9A3E56; */
            /* background-color: #f6f7fc; */
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Lureros</h1>
        </div>
        <div class="title">
            <strong>Verify your email address</strong>
        </div>
        <div class="content">
            <p>You need to verify your email address to continue using your account.</p>
            <p>Enter the following code to verify your email address:</p>
            <div class="code">%s</div>
        </div>
        <div class="instructions">
            <p>In case you were not trying to access your Account & are seeing this email, please follow the instructions below:</p>
            <ul>
                <li><a href="%s">Reset your password</a></li>
                <li>Check if any changes were made to your account & user settings. If yes, revert them immediately.</li>
            </ul>
        </div>
    </div>
</body>
</html>	
	`, captcha, ResetPasswordSite)
}

// context key
type contextKey string

const RequestContextKey contextKey = "Request"

// product
func CacheProductInCartKey(userId int64) string {
    return fmt.Sprintf("cache:product:cart:%d", userId)
}