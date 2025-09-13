package jsonrep

type UserAuth struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}
