package gosnel

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
)

func (g *Gosnel) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gosnel) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gosnel) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
	fp := path.Join(pathToFile, fileName)
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; file=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)
	return nil
}

func (g *Gosnel) Error404(w http.ResponseWriter, r *http.Request) {
	g.ErrorStatus(w, http.StatusNotFound)
}

func (g *Gosnel) Error500(w http.ResponseWriter, r *http.Request) {
	g.ErrorStatus(w, http.StatusInternalServerError)
}
func (g *Gosnel) ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	g.ErrorStatus(w, http.StatusUnauthorized)
}

func (g *Gosnel) ErrorForbidden(w http.ResponseWriter, r *http.Request) {
	g.ErrorStatus(w, http.StatusForbidden)
}

func (g *Gosnel) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
