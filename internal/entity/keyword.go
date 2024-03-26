package entity

type Keyword struct {
	ID       uint
	Title    string
	Slug     *string
	TopicID  uint
	Status   string    `gorm:"default:New"`
	Source   string    `gorm:"default:Unknown"`
	Articles []Article `gorm:"many2many:article_keywords;"`
}
