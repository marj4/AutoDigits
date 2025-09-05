package request

type RequestUserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RequestFromTransaction struct {
	TxType string `json:"tx_type"`
	Amount string `json:"amount"`
}
