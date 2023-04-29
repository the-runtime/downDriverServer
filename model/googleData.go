package model

type GoogleUserData struct {
	Id            string `json:"id""`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"'`
}
