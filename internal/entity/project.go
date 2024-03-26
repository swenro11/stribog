package entity

type Project struct {
	ID            uint
	Title         string
	Status        string `gorm:"default:New"`
	DirectoryPath string
	Description   string
	Url           string
	Topics        []Topic
}
