package entity

type Task struct {
	ID               uint
	Uuid             string
	Status           string
	ErrorDescription *string
	Prompt           *string
}
