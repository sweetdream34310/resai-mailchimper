package contacts

type ContactAddReq struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"required"`
}

type ContactUpdateReq struct {
	ID    string `json:"-"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}
