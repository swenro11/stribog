package entity

type Keyword struct {
	ID      uint
	Title   string
	Slug    *string
	Article []Article `gorm:"many2many:article_keywords;"`
}
