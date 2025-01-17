// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package proxy

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/mmclient"
)

func (p *Proxy) EnableApp(client mmclient.Client, sessionID string, cc *apps.Context, appID apps.AppID) (string, error) {
	app, err := p.GetInstalledApp(appID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get app. appID: %s", appID)
	}

	if !app.Disabled {
		return fmt.Sprintf("%s is already enabled", app.DisplayName), nil
	}

	_, err = client.EnableBot(app.BotUserID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to enable bot account for %s", app.AppID)
	}

	// Enable the app in the store first to allow calls to it
	app.Disabled = false
	err = p.store.App.Save(app)
	if err != nil {
		return "", errors.Wrapf(err, "failed to save app. appID: %s", appID)
	}

	var message string
	if app.OnEnable != nil {
		resp := p.Call(sessionID, cc.ActingUserID, &apps.CallRequest{
			Call:    *app.OnEnable,
			Context: cc,
		})
		if resp.Type == apps.CallResponseTypeError {
			p.log.WithError(err).Warnw("OnEnable failed, enabling app anyway",
				"app_id", app.AppID)
		} else {
			message = resp.Markdown
		}
	}

	if message == "" {
		message = fmt.Sprintf("Enabled %s", app.DisplayName)
	}

	p.log.Infow("Enabled app", "app_id", app.AppID)

	p.dispatchRefreshBindingsEvent(cc.ActingUserID)

	return message, nil
}

func (p *Proxy) DisableApp(client mmclient.Client, sessionID string, cc *apps.Context, appID apps.AppID) (string, error) {
	app, err := p.GetInstalledApp(appID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get app. appID: %s", appID)
	}

	if app.Disabled {
		return fmt.Sprintf("%s is already disabled", app.DisplayName), nil
	}

	// Call the app first as later it's disabled
	var message string
	if app.OnDisable != nil {
		resp := p.Call(sessionID, cc.ActingUserID, &apps.CallRequest{
			Call:    *app.OnDisable,
			Context: cc,
		})
		if resp.Type == apps.CallResponseTypeError {
			p.log.WithError(err).Warnw("OnDisable failed, disabling app anyway",
				"app_id", app.AppID)
		} else {
			message = resp.Markdown
		}
	}

	if message == "" {
		message = fmt.Sprintf("Disabled %s", app.DisplayName)
	}

	// disable app, not removing the data
	_, err = client.DisableBot(app.BotUserID)
	if err != nil {
		return "", errors.Wrapf(err, "failed to disable bot account for %s", app.AppID)
	}

	app.Disabled = true
	err = p.store.App.Save(app)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get app. appID: %s", appID)
	}

	p.log.Infow("Disabled app",
		"app_id", app.AppID)

	p.dispatchRefreshBindingsEvent(cc.ActingUserID)

	return message, nil
}

func (p *Proxy) AppIsEnabled(app *apps.App) bool {
	if app.AppType == apps.AppTypeBuiltin {
		return true
	}
	if app.Disabled {
		return false
	}
	if m, _ := p.store.Manifest.Get(app.AppID); m == nil {
		return false
	}
	return true
}
