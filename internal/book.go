package internal

type Book struct {
	ISBN          string  `json:"isbn"`
	Title         string  `json:"title"`
	SubTitle      string  `json:"subTitle"`
	Authors       string  `json:"authors"`
	Categories    string  `json:"categories"`
	Thumbnail     string  `json:"thumbnail"`
	Description   string  `json:"description"`
	PublishedYear int     `json:"publishedYear"`
	AverageRating float32 `json:"averageRating"`
	NumPages      int     `json:"numPages"`
	RatingCount   int     `json:"ratingCount"`
}

func NewBook(isbn, title, subtitle, thumbnail, description string, publishYear, numPages, ratingCount int, avgRating float32, authors, categories string) *Book {
	return &Book{
		ISBN:          isbn,
		Title:         title,
		SubTitle:      subtitle,
		Authors:       authors,
		Categories:    categories,
		Thumbnail:     thumbnail,
		Description:   description,
		PublishedYear: publishYear,
		AverageRating: avgRating,
		NumPages:      numPages,
		RatingCount:   ratingCount,
	}
}
