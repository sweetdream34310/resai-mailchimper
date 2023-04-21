package label

type CreateLabelReq struct {
	Name string `json:"name" binding:"required"`
}

type PatchLabelReq struct {
	Name string `json:"name" binding:"required"`
}
