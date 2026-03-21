package domain

type UpsertUserRequest struct {
	GoogleID  string
	FullName  string
	Email     string
	AvatarURL string
}
