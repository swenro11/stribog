package entity

type Topic struct {
	ID            uint
	Title         string
	ProjectID     uint
	Status        string `gorm:"default:New"`
	Keywords      []Keyword
	ParentTopicID *uint
	ParentTopic   *Topic
}
