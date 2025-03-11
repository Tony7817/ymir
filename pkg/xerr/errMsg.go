package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "SUCCESS"
	message[ServerCommonError] = "internal Server Error"
	message[UnauthorizedError] = "not authorized"
	message[ErrorReuqestParam] = "bad request param"
	message[CaptchaExpireError] = "verification code has expired"
	message[WrongCaptchaError] = "wrong verification code"
	message[TokenExpireError] = "token失效，请重新登陆"
	message[TokenGenerateError] = "生成token失败"
	message[DbError] = "数据库繁忙,请稍后再试"
	message[DbUpdateAffectedZeroError] = "update data affected zero"
	message[DataNoExistError] = "Data not exist"
	message[ErrorResourceForbiden] = "resource forbidden"
	// user service
	message[UserAlreadyExistError] = "user already signed up"
	message[MaxCaptchaSendTimeError] = "maxmum captcha send time reached"
	message[ErrorNotAuthorized] = "not authorized"
	message[ErrorInvalidEmail] = "invalid email"
	message[ErrorSignedupInGoogle] = "user has signed up by google, please sign in with google"
	message[ErrorWrongPassword] = "incrrect password"
	message[UserNotExistedError] = "user is not signed up"
	// product service
	message[ErrorOutOfStockError] = "product is out of stock"
	// order service
	message[ErrorRequestOrderMaximunReach] = "request too many orders"
	message[ErrorOrderNotExist] = "order does not exist"
	message[ErrorCreateOrder] = "create order failed, please try again"
	message[ErrorPayOrder] = "pay order failed, please try again"
	message[ErrorInvalidOrderStatus] = "invalid order"
	message[ErrorDeleteOrder] = "delete order failed, please try again"
	message[ErrorIdempotence] = "order has been created"
	message[ErrorStockNotEnough] = "stock is not enough"
	message[ErrorFetchResource] = "fetch resource failed"
	message[ErrorOrderCreated] = "order has been created"
	message[ErrorOrderPaied] = "order has been paid"
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
