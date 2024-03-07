package discordmodel

type User struct {
	// the user's id
	id string
	// the user's username, not unique across the platform
	username string
	// the user's Discord-tag
	discriminator string
	// the user's display name, if it is set. For bots, this is the application name
	global_name *string
	// the user's avatar hash
	avatar *string
	// whether the user belongs to an OAuth2 application
	bot *bool
	// whether the user is an Official Discord System user (part of the urgent message system)
	system *bool
	// whether the user has two factor enabled on their account
	mfa_enabled *bool
	// the user's banner hash
	banner *string
	// the user's banner color encoded as an integer representation of hexadecimal color code
	accent_color *int
	// the user's chosen language option
	locale *string
	// whether the email on this account has been verified
	verified *bool
	// the user's email
	email *string
	// the flags on a user's account
	flags *int
	// the type of Nitro subscription on a user's account
	premium_type *int
	// the public flags on a user's account
	public_flags *int
	// the user's avatar decoration hash
	avatar_decoration *string
}

type AccessTokenRequest struct {
	// Must be set to 'authorization_code'
	GrantType string `json:"grant_type"`
	// The code from the querystring
	Code string `json:"code"`
	// The 'redirect_uri' associated with this authorization, usually from your authorization URL
	RedirectUri string `json:"redirect_uri"`
}

type RefreshTokenRequest struct {
	// Must be set to 'refresh_token'
	GrantType string `json:"grant_type"`
	// The user's refresh token
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
