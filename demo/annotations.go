package demo

// AUTO GENERATED - DO NOT EDIT

// Some Person
type Person struct {
	Name  string `maxlen:"256" minlen:"2"`
	Email string `format:"email"`
	Age   uint8  `max:"40"`
	Phone string `pattern:"\d+"`
}

type Book struct {
	Title     string   `json:"title"`
	PageCount int32    `json:"page_count"`
	Authors   []Person `json:"authors,omitempty"`
	Content   string   `maxlen:"256" json:"content"`
}

