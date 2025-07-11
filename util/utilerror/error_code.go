package utilerror

const (
	ErrorCodeSystemError              = -600
	ErrorCodeRequestParamCanNotBeNull = -1000
	ErrorCodeFile                     = -1001
)

var (
	ErrorRequestParamCanNotBeNull = NewError("Request param can not be null")
)
