package domain

type UpsertUserRequest struct {
	GoogleID  string
	FullName  string
	Email     string
	AvatarURL string
}

type PatchUserRequest struct {
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone"`
}
