package request

// TODO: change Username everywhere for Username (including db)

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
