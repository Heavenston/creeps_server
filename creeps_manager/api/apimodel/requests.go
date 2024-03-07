package apimodel

type CreateGameRequest struct {
	Config *GameConfig `json:"config",omitempty`
	Name   string      `json:"name"`
}
