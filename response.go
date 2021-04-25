package yu

// 返回code值
const (
	CodeOK = iota
	CodeErr
)

// Resp http统一返回结构
type Resp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

// RespOK 创建一个返回OK的结构
func RespOK(data interface{}) *Resp {
	return &Resp{Code: CodeOK, Data: data}
}

// RespErr 创建一个返回error的结构
func RespErr(err error) *Resp {
	return &Resp{Code: CodeErr, Data: err.Error()}
}

// RespMsg 创建一个返回字符串信息的结构
func RespMsg(code int, msg string) *Resp {
	return &Resp{Code: code, Data: msg}
}

// NewResp 创建一个返回结构
func NewResp(code int, data interface{}) *Resp {
	return &Resp{Code: code, Data: data}
}
