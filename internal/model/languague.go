package model

type Language struct {
	Id        string `db:"id" json:"id"`
	Code      string `db:"code" json:"code"`
	Name      string `db:"name" json:"name"`
	FlagURL   string `db:"flag_url" json:"flag_url"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

// CreateLanguageRequest represents the request body for creating a language
type CreateLanguageRequest struct {
	Code    string `json:"code" validate:"required"`
	Name    string `json:"name" validate:"required"`
	FlagURL string `json:"flag_url"`
}

// UpdateLanguageRequest represents the request body for updating a language
type UpdateLanguageRequest struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	FlagURL string `json:"flag_url"`
}
