package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "SUCCESS"
	message[ServerCommonError] = "internal Server Error"
	message[UnauthorizedError] = "not authorized"
	message[ReuqestParamError] = "bad request param"
	message[CaptchaExpireError] = "verification code has expired"
	message[WrongCaptchaError] = "wrong verification code"
	message[TokenExpireError] = "token失效，请重新登陆"
	message[TokenGenerateError] = "生成token失败"
	message[DbError] = "数据库繁忙,请稍后再试"
	message[DbUpdateAffectedZeroError] = "更新数据影响行数为0"
	message[DataNoExistError] = "Data not exist"
	// user service
	message[UserAlreadyExistError] = "user already signed up"
	message[MaxCaptchaSendTimeError] = "maxmum captcha send time reached"
	message[NotAuthorizedError] = "not authorized"
	message[ErrorInvalidEmail] = "invalid email"
	message[ErrorSignedupInGoogle] = "user has signed up by google, please sign in with google"
	message[ErrorWrongPassword] = "incrrect password"
	message[UserNotExistedError] = "user not exist"
	// product service
	message[OutOfStockError] = "product is out of stock"
}

func MapErrMsg(errcode uint32) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "internal Server Error"
	}
}

func IsCodeErr(errcode uint32) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
