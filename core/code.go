package core

var (
	Success = NewError(0, "success")
	// 10000-10100 系统错误
	DBError    = NewError(10000, "db is error")
	InnerError = NewError(10001, "inner error")
	// 10100-10200 定时循环任务相关错误
	CronCycleError = NewError(10100, "cron cycle is error")

	// 10200-10300 固定定时任务执行相关错误
	FixCronError = NewError(10200, "cron is error")

	// 10400-10500 Job执行相关错误
	JobError = NewError(10500, "job flow is error")

	AlreadyRegister  = NewError(10100, "user already register")
	NameOrPwdError   = NewError(10101, "username or password error")
	TokenError       = NewError(10102, "token error")
	TokenExpired     = NewError(10103, "token expired")
	TokenInvalid     = NewError(10104, "token invalid")
	UserNotExist     = NewError(10105, "user not exist")
	UserAlreadyExist = NewError(10106, "user already exist")
)
