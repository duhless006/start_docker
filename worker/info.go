package worker

type Worker struct {
	Id       int    `json:"id"`
	FullName string `json:"full_name"`
	Position string `json:"position"`
}

func NewWorker(id int, fullName string, position string) *Worker {
	return &Worker{
		Id:       id,
		FullName: fullName,
		Position: position,
	}
}
