package entity

type Image struct {
	ID           uint
	AltTitle     *string
	Slug         *string
	ArticleID    uint
	Sort         *uint
	Status       string `gorm:"default:New"`
	RewriteNotes *string
	Prompt       *string
	Link         *string
	Path         *string
	Base64       *string
}
