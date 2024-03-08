package atomic

const (
	ERR_INTERNAL_SERVER_ERROR = "7xTpC28cYvozzFR4B40su"
	ERR_UNAUTHORIZED          = "45EEaIn6XgQ6sH8Rm3WRo"
	ERR_FORBIDDEN             = "SsKEOBm0hk1NZTXTyXdp2"
)

func (e *engine) initErrors() {
	e.Errors[ERR_INTERNAL_SERVER_ERROR] = &Error{
		UUID: ERR_INTERNAL_SERVER_ERROR,
		Eng:  "Internal Server Error (500)",
		Zht:  "內部伺服器錯誤 (500)",
		Zhs:  "内部服务器错误 (500)",
	}

	e.Errors[ERR_UNAUTHORIZED] = &Error{
		UUID: ERR_UNAUTHORIZED,
		Eng:  "Unauthorized (401)",
		Zht:  "未經授權 (401)",
		Zhs:  "未经授权 (401)",
	}

	e.Errors[ERR_FORBIDDEN] = &Error{
		UUID: ERR_FORBIDDEN,
		Eng:  "Forbidden (403)",
		Zht:  "無法訪問 (403)",
		Zhs:  "无法访问 (403)",
	}
}
