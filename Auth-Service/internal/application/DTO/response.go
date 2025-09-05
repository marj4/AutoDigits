package dto

type Response struct {
	Status   string
	Message  string
	Uuid     string
	Password string
	Role     string
}

const (
	Error    = "Error"
	Success  = "Success"
	NotFound = "Not found"
)
