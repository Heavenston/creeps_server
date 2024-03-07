package apimodel

type User struct {
	Id         int    `json:"id"`
	DiscordId  string `json:"discord_id"`
	DiscordTag string `json:"discord_tag"`
	AvatarUrl  *string `json:"avatar_url"`
	Username   string `json:"username"`
}
