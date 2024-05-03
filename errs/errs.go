package errs

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Event   string `json:"event"`
}

func NewErr(e error, event string, msgCode string) *Err {
	msg := e.Error()
	var code int
	switch msgCode {
	case "not found":
		code = 404
	case "bad data":
		code = 400
	default:
		code = 500
	}
	return &Err{
		Code:    code,
		Message: msg,
		Event:   event,
	}
}

func WrapErr(e *Err, msg string) *Err {
	e.Message = msg + e.Message
	return e
}
