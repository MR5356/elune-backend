package response

const (
	CodeSuccess      = 0
	CodeParamError   = -1
	CodeNotFound     = -404
	CodeUnknownError = -999
)

const (
	MsgSuccess    = "success"
	MsgParamError = "param error"
	MsgNotFound   = "not found"

	MsgUnknownError = "unknown error"
)

var (
	responseMap = map[int]string{
		CodeSuccess:    MsgSuccess,
		CodeParamError: MsgParamError,
		CodeNotFound:   MsgNotFound,
	}
)
