package schema

type RequestModel struct {
	Year int `validate:"required,number,min=2020,max=2025"`
}
