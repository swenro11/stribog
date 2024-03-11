package entity

type Article struct {
	ID               uint
	Title            string
	Slug             *string
	ShortDescription *string
	Body             *string
	Status           string
	RewriteNotes     *string
	Prompt           *string
	Images           []Image
	Keyword          []Keyword `gorm:"many2many:article_keywords;"`
}
