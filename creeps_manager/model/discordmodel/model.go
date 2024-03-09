package discordmodel

import (
	"fmt"
	"strconv"
)

type User struct {
	// the user's id
	Id string `json:"id"`
	// the user's username, not unique across the platform
	Username string `json:"username"`
	// the user's Discord-tag
	Discriminator string `json:"discriminator"`
	// the user's display name, if it is set. For bots, this is the application name
	GlobalName *string `json:"global_name"`
	// the user's avatar hash
	Avatar *string `json:"avatar"`
	// whether the user belongs to an OAuth2 application
	Bot *bool `json:"bot"`
	// whether the user is an Official Discord System user (part of the urgent message system)
	System *bool `json:"system"`
	// whether the user has two factor enabled on their account
	MfaEnabled *bool `json:"mfa_enabled"`
	// the user's banner hash
	Banner *string `json:"banner"`
	// the user's banner color encoded as an integer representation of hexadecimal color code
	AccentColor *int `json:"accent_color"`
	// the user's chosen language option
	Locale *string `json:"locale"`
	// whether the email on this account has been verified
	Verified *bool `json:"verified"`
	// the user's email
	Email *string `json:"email"`
	// the flags on a user's account
	Flags *int `json:"flags"`
	// the type of Nitro subscription on a user's account
	PremiumType *int `json:"premium_type"`
	// the public flags on a user's account
	PublicFlags *int `json:"public_flags"`
	// the user's avatar decoration hash
	AvatarDecoration *string `json:"avatar_decoration"`
}

func (user *User) AvatarUrl() string {
	if user.Avatar != nil {
		return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", user.Id, *user.Avatar)
	} else {
		id, _ := strconv.ParseInt(user.Id, 10, 64)
		return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", (id >> 22) % 6)
	}
}

type AccessTokenRequest struct {
	// Must be set to 'authorization_code'
	GrantType string `form:"grant_type"`
	// The code from the querystring
	Code string `form:"code"`
	// The 'redirect_uri' associated with this authorization, usually from your authorization URL
	RedirectUri string `form:"redirect_uri"`
}

type RefreshTokenRequest struct {
	// Must be set to 'refresh_token'
	GrantType string `form:"grant_type"`
	// The user's refresh token
	RefreshToken string `form:"refresh_token"`
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (a *AccessTokenResponse) AuthHeader() string {
	return fmt.Sprintf("Bearer %s", a.AccessToken)
}
