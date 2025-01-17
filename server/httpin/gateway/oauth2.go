package gateway

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-apps/utils"
	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
)

func (g *gateway) remoteOAuth2Connect(w http.ResponseWriter, req *http.Request, sessionID, actingUserID string) {
	appID := appIDVar(req)

	if appID == "" {
		httputils.WriteError(w, utils.NewInvalidError("app_id not specified"))
		return
	}

	connectURL, err := g.proxy.GetRemoteOAuth2ConnectURL(sessionID, actingUserID, appID)
	if err != nil {
		g.log.WithError(err).Warnw("Failed to get remote OuAuth2 connect URL",
			"app_id", appID,
			"acting_user_id", actingUserID)
		httputils.WriteError(w, err)
		return
	}

	http.Redirect(w, req, connectURL, http.StatusTemporaryRedirect)
}

func (g *gateway) remoteOAuth2Complete(w http.ResponseWriter, req *http.Request, sessionID, actingUserID string) {
	appID := appIDVar(req)

	if appID == "" {
		httputils.WriteError(w, utils.NewInvalidError("app_id not specified"))
		return
	}

	q := req.URL.Query()
	urlValues := map[string]interface{}{}
	for key := range q {
		urlValues[key] = q.Get(key)
	}

	err := g.proxy.CompleteRemoteOAuth2(sessionID, actingUserID, appID, urlValues)
	if err != nil {
		g.log.WithError(err).Warnw("Failed to complete remote OuAuth2",
			"app_id", appID,
			"acting_user_id", actingUserID)
		httputils.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(`
	<!DOCTYPE html>
	<html>
		<head>
			<script>
				window.close();
			</script>
		</head>
		<body>
			<p>Completed connecting your account. Please close this window.</p>
		</body>
	</html>
	`))
}
