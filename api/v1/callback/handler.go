package callback

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/matthiasharzer/go-stats-tracker/api"
	"golang.org/x/oauth2"
)

//go:embed index.html
var templateHTML string

type TemplateData struct {
	APIKey        string
	AccessToken   string
	UserAccessKey string
	State         string
}

func Handler(sharedState *api.SharedState, oauth oauth2.Config, filePickerAPIKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state == "" {
			http.Error(w, "State is empty", http.StatusBadRequest)
			return
		}

		authFlowState := sharedState.PeekState(state)
		if authFlowState == nil {
			http.Error(w, "Unknown state", http.StatusBadRequest)
			return
		}

		code := r.FormValue("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		token, err := oauth.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		sharedState.UpdateState(state, api.AuthFlowState{
			UserAccessKey: authFlowState.UserAccessKey,
			RefreshToken:  token.RefreshToken,
		})

		templateData := TemplateData{
			APIKey:        filePickerAPIKey,
			AccessToken:   token.AccessToken,
			UserAccessKey: authFlowState.UserAccessKey,
			State:         state,
		}

		tmpl := template.Must(template.New("index").Parse(templateHTML))
		err = tmpl.Execute(w, templateData)
		if err != nil {
			http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
