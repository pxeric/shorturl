package models

//定义返回code
type ReturnCode int

const (
	Fail         ReturnCode = -1
	MissParams              = -2
	Success                 = 1
	NotExist                = 21
	Sended                  = 22
	SmsCodeError            = 23
	NoLogin                 = 24
	Exist                   = 25
	NoPermission            = 57
)
