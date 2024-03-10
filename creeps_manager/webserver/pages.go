package webserver

import (
	"net/http"

	"github.com/Heavenston/creeps_server/creeps_manager/model"
	"github.com/Heavenston/creeps_server/creeps_manager/templates"
	"github.com/a-h/templ"
)

func (self *WebServer) getIndex(w http.ResponseWriter, r *http.Request) {
    var games []model.Game
    self.Db.Find(&games)

    ctx := templ.WithChildren(r.Context(), templates.Index(games))
    templates.Layout(templates.IndexHeader()).
        Render(ctx, w)
}

func (self *WebServer) getCreateGame(w http.ResponseWriter, r *http.Request) {
    ctx := templ.WithChildren(r.Context(), templates.CreateGame())
    templates.Layout(templates.CreateGameHeader()).
        Render(ctx, w)
}

