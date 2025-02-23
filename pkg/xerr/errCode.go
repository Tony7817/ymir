package xerr

// 成功返回
const OK uint32 = 200

/**(前3位代表业务,后三位代表具体功能)**/

/**全局错误码*/
//服务器开小差
const ServerCommonError uint32 = 100001

// 请求参数错误
const ErrorReuqestParam uint32 = 100002

// token过期
const TokenExpireError uint32 = 100003

// 生成token失败
const TokenGenerateError uint32 = 100004

// 数据库繁忙,请稍后再试
const DbError uint32 = 100005

// 更新数据影响行数为0
const DbUpdateAffectedZeroError uint32 = 100006

// 数据不存在
const DataNoExistError uint32 = 100007

const UnauthorizedError uint32 = 100008

// user service 200000-200999
const UserAlreadyExistError uint32 = 201001
const CaptchaExpireError uint32 = 201002
const WrongCaptchaError uint32 = 201003
const MaxCaptchaSendTimeError uint32 = 201004
const ErrorNotAuthorized uint32 = 201005
const ErrorInvalidEmail uint32 = 201006
const ErrorSignedupInGoogle uint32 = 201007
const ErrorWrongPassword uint32 = 201008
const UserNotExistedError uint32 = 201009

// 订单服务 300000 - 300999
const ErrorRequestOrderMaximunReach uint32 = 300001
const ErrorOrderNotExist uint32 = 300002
const ErrorCreateOrder uint32 = 300003
const ErrorPayOrder uint32 = 300004
const ErrorInvalidOrderStatus uint32 = 300005
const ErrorDeleteOrder uint32 = 300006
const ErrorIdempotence uint32 = 300007

// 商品服务
const ErrorOutOfStockError uint32 = 400001
const ErrorStockNotEnough uint32 = 400002

//支付服务
