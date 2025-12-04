package plugins

import (
	"net/http"

	"github.com/google/uuid"
)

func (f *Plugins) HandleFiles() {
	f.Frontend.HandleFunc("/plugin/", func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.URL.Query().Get("jwt")
		pluginId := r.URL.Query().Get("plugin")

		if jwtToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		pluginUuid, err := uuid.Parse(pluginId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p, ok := f.Plugins[pluginUuid]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		data := p.Tokens.AccessToken
		data.Token = jwtToken
		_, err = f.Plugins[pluginUuid].Tokens.ValidateToken(data)
		if err != nil {
			if err == ErrTokenExpired {
				p := f.Plugins[pluginUuid]
				_, err = f.Plugins[pluginUuid].Tokens.ValidateToken(p.Tokens.RefreshToken)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if err := p.Tokens.GetToken(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write([]byte(p.Tokens.AccessToken.Token))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		http.StripPrefix("/plugin/"+p.Plugnin.Id+"/", http.FileServer(http.Dir("./plugins/"+p.Plugnin.Id))).ServeHTTP(w, r)
	})
}

func (f *Plugins) ListenAndServe(addr string) error {
	f.HandleFiles()
	return http.ListenAndServe(addr, f.Frontend)
}
