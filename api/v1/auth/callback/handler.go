package callback

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/matthiasharzer/go-stats-tracker/api/v1/auth"
	"golang.org/x/oauth2"
)

//go:embed index.html
var templateHTML string

type TemplateData struct {
	DeveloperKey  string
	AppID         string
	AccessToken   string
	UserAccessKey string
	State         string
}

func Handler(sharedState *auth.SharedState, oauth oauth2.Config, filePickerDeveloperKey string, filePickerAppID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateData2 := TemplateData{
			DeveloperKey:  filePickerDeveloperKey,
			AppID:         filePickerAppID,
			AccessToken:   "token.AccessToken",
			State:         "state",
			UserAccessKey: "userAccessKey",
		}
		tmpl2 := template.Must(template.New("index").Parse(templateHTML))
		err2 := tmpl2.Execute(w, templateData2)
		if err2 != nil {
			http.Error(w, "Failed to execute template: "+err2.Error(), http.StatusInternalServerError)
			return
		}
		return
		state := r.FormValue("state")
		if state == "" {
			http.Error(w, "State is empty", http.StatusBadRequest)
			return
		}

		userAccessKey := sharedState.GetUserAccessKey(state)
		if userAccessKey == "" {
			http.Error(w, "Unknown callback state", http.StatusBadRequest)
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

		sharedState.SetRefreshToken(state, token.RefreshToken)

		templateData := TemplateData{
			DeveloperKey:  filePickerDeveloperKey,
			AppID:         filePickerAppID,
			AccessToken:   token.AccessToken,
			State:         state,
			UserAccessKey: userAccessKey,
		}

		tmpl := template.Must(template.New("index").Parse(templateHTML))
		err = tmpl.Execute(w, templateData)
		if err != nil {
			http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
