package dto

type Response struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	UUID     string // Need to Auth-Service
	Password string // Need to Auth-Service
	Role     string // Need to Auth-Service
}

const (
	Error    = "Internal error"
	NotExist = "Not exist"
	Success  = "Success"
	Exist    = "Exist"
)
