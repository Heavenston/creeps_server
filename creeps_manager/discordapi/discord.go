package discordapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/Heavenston/creeps_server/creeps_manager/model/discordmodel"
	"github.com/ajg/form"
	"github.com/rs/zerolog/log"
)

type IDiscordAuth interface {
	AuthHeader() string
}

type DiscordAppAuth struct {
	ClientId     string
	ClientSecret string
}

func (a *DiscordAppAuth) AuthHeader() string {
	data := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.ClientId, a.ClientSecret)))
	return fmt.Sprintf("Basic %s", data)
}

type DiscordBearerAuth struct {
	AccessToken string
	DiscordId *string
}

func (a *DiscordBearerAuth) AuthHeader() string {
	return fmt.Sprintf("Bearer %s", a.AccessToken)
}

type ErrNonOk struct {
	status int
	body string
}

func (e *ErrNonOk) Error() string {
	return fmt.Sprintf("%d: %s", e.status, e.body)
}

const API_BASE_URI string = "https://discord.com/api/v10"

func applyAuth(auth IDiscordAuth, req *http.Request) {
	req.Header.Add("Authorization", auth.AuthHeader())
}

func get(auth IDiscordAuth, uri string, out any) error {
	req, err := http.NewRequest("GET", API_BASE_URI + uri, nil)
	if err != nil {
		return err
	}

	applyAuth(auth, req)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return &ErrNonOk{
			status: resp.StatusCode,
			body: string(body),
		}
	}

	log.Trace().
		Str("response_body", string(body)).
		Str("uri", uri).
		Type("auth_type", auth).
		Msg("Made a discord get request")

	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}

	return nil
}

func MakeAccessTokenRequest(
	auth *DiscordAppAuth,
	code string,
	redirectUri string,
) (discordmodel.AccessTokenResponse, error) {
	values, err := form.EncodeToValues(&discordmodel.AccessTokenRequest{
		GrantType:   "authorization_code",
		Code:        code,
		RedirectUri: redirectUri,
	})
	values.Add("client_id", auth.ClientId)
	values.Add("client_secret", auth.ClientSecret)
	if err != nil {
		return discordmodel.AccessTokenResponse{}, err
	}

	resp, err := http.PostForm("https://discord.com/api/oauth2/token", values)
	if err != nil {
		log.Debug().Any("values", values).Err(err).Msg("access token error")
		return discordmodel.AccessTokenResponse{}, err
	}

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Debug().Any("values", values).Err(err).Msg("access token error")
		return discordmodel.AccessTokenResponse{}, fmt.Errorf("%s", string(body))
	}

	var atr discordmodel.AccessTokenResponse
	err = json.Unmarshal(body, &atr)

	return atr, err
}

var userCache sync.Map

func GetCurrentUser(auth IDiscordAuth) (user discordmodel.User, err error) {
	if ba, ok := auth.(*DiscordBearerAuth); ok && ba.DiscordId != nil {
		val, ok := userCache.Load(ba.DiscordId)
		if ok {
			user = *val.(*discordmodel.User)
			return
		}
	}

	err = get(auth, "/users/@me", &user)

	if err == nil {
		userCache.Store(user.Id, &user)
	}

	return
}
