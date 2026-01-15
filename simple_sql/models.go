package simple_sql

type Personal struct {
	Id       int    `json:"id"`
	FullName string `json:"full_name"`
	Position string `json:"position"`
}

func NewPersonal(id int, fullName string, position string) *Personal {
	return &Personal{
		Id:       id,
		FullName: fullName,
		Position: position,
	}
}
