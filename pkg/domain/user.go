package domain

type User struct {
	ID             string `json:"id"`
	GoogleID       string `json:"google_id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
	AccessToken    string `json:"-"`
	RefreshToken   string `json:"-"`
}
