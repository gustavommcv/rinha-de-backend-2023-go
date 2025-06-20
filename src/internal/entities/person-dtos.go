package entities

type PersonRequestDTO struct {
	Surname   string   `json:"apelido"`
	Name      string   `json:"nome"`
	Birthdate string   `json:"nascimento"`
	Stack     []string `json:"stack"`
}

type PersonResponseDTO struct {
	Id        string   `json:"id"`
	Surname   string   `json:"apelido"`
	Name      string   `json:"nome"`
	Birthdate string   `json:"nascimento"`
	Stack     []string `json:"stack"`
}
