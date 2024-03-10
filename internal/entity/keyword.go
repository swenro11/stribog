package entity

type Keyword struct {
	ID      uint
	Title   string
	Slug    *string
	Topic   string
	Status  string    `gorm:"default:New"`
	Article []Article `gorm:"many2many:article_keywords;"`
}
