package createControllerTask

type InputCreate struct {
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description" validate:"required,max=500"`
}
