package callback

import (
	_ "embed"
	"html/template"
	"net/http"
	"os"

	"github.com/matthiasharzer/go-stats-tracker/api"
	"golang.org/x/oauth2"
)

//go:embed index.html
var templateHTML string

type TemplateData struct {
	DeveloperKey  string
	AccessToken   string
	UserAccessKey string
	State         string
}

func getTemplateData() TemplateData {
	developerKey := os.Getenv("DEVELOPER_KEY")
	accessToken := os.Getenv("ACCESS_TOKEN")
	userID := os.Getenv("USER_ID")

	if developerKey == "" || accessToken == "" || userID == "" {
		panic("Environment variables DEVELOPER_KEY, ACCESS_TOKEN, and USER_ID must be set")
	}

	return TemplateData{
		DeveloperKey:  developerKey,
		AccessToken:   accessToken,
		UserAccessKey: userID,
	}
}

func Handler(sharedState *api.SharedState, oauth oauth2.Config) http.HandlerFunc {
	developerKEy := os.Getenv("DEVELOPER_KEY")
	if developerKEy == "" {
		panic("Environment variable DEVELOPER_KEY must be set")
	}
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
			DeveloperKey:  developerKEy,
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
