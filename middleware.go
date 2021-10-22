package gosnel

import "net/http"

// SessionLoad loads and saves session on every request
func (g *Gosnel) SessionLoad(next http.Handler) http.Handler {
	g.InfoLog.Println("SessionLoad called")
	return g.Session.LoadAndSave(next)
}
