// Copyright (c) 2020-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package store

import (
	"crypto/sha1" // nolint:gosec
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/config"
	"github.com/mattermost/mattermost-plugin-apps/server/httpout"
	"github.com/mattermost/mattermost-plugin-apps/upstream/upaws"
	"github.com/mattermost/mattermost-plugin-apps/utils"
)

type ManifestStore interface {
	config.Configurable

	AsMap() map[apps.AppID]*apps.Manifest
	DeleteLocal(apps.AppID) error
	Get(apps.AppID) (*apps.Manifest, error)
	GetFromS3(apps.AppID, apps.AppVersion) (*apps.Manifest, error)
	InitGlobal(httpout.Service) error
	StoreLocal(*apps.Manifest) error
}

// manifestStore combines global (aka marketplace) manifests, and locally
// installed ones. The global list is loaded on startup. The local manifests are
// stored in KV store, and the list of their keys is stored in the config, as a
// map of AppID->sha1(manifest).
type manifestStore struct {
	*Service

	// mutex guards local, the pointer to the map of locally-installed
	// manifests.
	mutex sync.RWMutex

	global map[apps.AppID]*apps.Manifest
	local  map[apps.AppID]*apps.Manifest
}

var _ ManifestStore = (*manifestStore)(nil)

// InitGlobal reads in the list of known (i.e. marketplace listed) app
// manifests.
func (s *manifestStore) InitGlobal(httpOut httpout.Service) error {
	bundlePath, err := s.mm.System.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "can't get bundle path")
	}
	assetPath := filepath.Join(bundlePath, "assets")
	f, err := os.Open(filepath.Join(assetPath, config.ManifestsFile))
	if err != nil {
		return errors.Wrap(err, "failed to load global list of available apps")
	}
	defer f.Close()

	global := map[apps.AppID]*apps.Manifest{}
	manifestLocations := map[apps.AppID]string{}
	err = json.NewDecoder(f).Decode(&manifestLocations)
	if err != nil {
		return err
	}

	conf := s.conf.GetConfig()
	var data []byte
	for appID, loc := range manifestLocations {
		parts := strings.SplitN(loc, ":", 2)
		switch {
		case len(parts) == 1:
			data, err = s.getDataFromS3(appID, apps.AppVersion(parts[0]))
		case len(parts) == 2 && parts[0] == "s3":
			data, err = s.getDataFromS3(appID, apps.AppVersion(parts[1]))
		case len(parts) == 2 && parts[0] == "file":
			data, err = os.ReadFile(filepath.Join(assetPath, parts[1]))
		case len(parts) == 2 && (parts[0] == "http" || parts[0] == "https"):
			data, err = httpOut.GetFromURL(loc, conf.DeveloperMode)
		default:
			s.log.WithError(err).Errorw("Failed to load global manifest",
				"app_id", appID)
			continue
		}
		if err != nil {
			s.log.WithError(err).Errorw("Failed to load global manifest",
				"app_id", appID,
				"loc", loc)
			continue
		}

		var m *apps.Manifest
		m, err = apps.ManifestFromJSON(data)
		if err != nil {
			s.log.WithError(err).Errorw("Failed to load global manifest",
				"app_id", appID,
				"loc", loc)
			continue
		}
		if m.AppID != appID {
			err = errors.Errorf("mismatched app ids while getting manifest %s != %s", m.AppID, appID)
			s.log.WithError(err).Errorw("Failed to load global manifest",
				"app_id", appID,
				"loc", loc)
			continue
		}
		global[appID] = m
	}

	s.mutex.Lock()
	s.global = global
	s.mutex.Unlock()

	return nil
}

func DecodeManifest(data []byte) (*apps.Manifest, error) {
	var m apps.Manifest
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	err = m.IsValid()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *manifestStore) Configure(conf config.Config) {
	updatedLocal := map[apps.AppID]*apps.Manifest{}

	for id, key := range conf.LocalManifests {
		var m *apps.Manifest
		err := s.mm.KV.Get(config.KVLocalManifestPrefix+key, &m)
		switch {
		case err != nil:
			s.log.WithError(err).Errorw("Failed to load local manifest",
				"app_id", id)

		case m == nil:
			s.log.WithError(utils.ErrNotFound).Errorw("Failed to load local manifest",
				"app_id", id)

		default:
			updatedLocal[apps.AppID(id)] = m
		}
	}

	s.mutex.Lock()
	s.local = updatedLocal
	s.mutex.Unlock()
}

func (s *manifestStore) Get(appID apps.AppID) (*apps.Manifest, error) {
	s.mutex.RLock()
	local := s.local
	global := s.global
	s.mutex.RUnlock()

	m, ok := local[appID]
	if ok {
		return m, nil
	}
	m, ok = global[appID]
	if ok {
		return m, nil
	}
	return nil, utils.ErrNotFound
}

func (s *manifestStore) AsMap() map[apps.AppID]*apps.Manifest {
	s.mutex.RLock()
	local := s.local
	global := s.global
	s.mutex.RUnlock()

	out := map[apps.AppID]*apps.Manifest{}
	for id, m := range global {
		out[id] = m
	}
	for id, m := range local {
		out[id] = m
	}
	return out
}

func (s *manifestStore) StoreLocal(m *apps.Manifest) error {
	conf := s.conf.GetConfig()
	prevSHA := conf.LocalManifests[string(m.AppID)]

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	sha := fmt.Sprintf("%x", sha1.Sum(data)) // nolint:gosec
	if sha == prevSHA {
		return nil
	}

	_, err = s.mm.KV.Set(config.KVLocalManifestPrefix+sha, m)
	if err != nil {
		return err
	}

	s.mutex.RLock()
	local := s.local
	s.mutex.RUnlock()
	updatedLocal := map[apps.AppID]*apps.Manifest{}
	for k, v := range local {
		if k != m.AppID {
			updatedLocal[k] = v
		}
	}
	updatedLocal[m.AppID] = m
	s.mutex.Lock()
	s.local = updatedLocal
	s.mutex.Unlock()

	updated := map[string]string{}
	for k, v := range conf.LocalManifests {
		updated[k] = v
	}
	updated[string(m.AppID)] = sha
	sc := conf.StoredConfig
	sc.LocalManifests = updated
	err = s.conf.StoreConfig(sc)
	if err != nil {
		return err
	}

	err = s.mm.KV.Delete(config.KVLocalManifestPrefix + prevSHA)
	if err != nil {
		s.log.WithError(err).Warnf("Failed to delete previous Manifest KV value")
	}
	return nil
}

func (s *manifestStore) DeleteLocal(appID apps.AppID) error {
	conf := s.conf.GetConfig()
	sha := conf.LocalManifests[string(appID)]

	err := s.mm.KV.Delete(config.KVLocalManifestPrefix + sha)
	if err != nil {
		return err
	}

	s.mutex.RLock()
	local := s.local
	s.mutex.RUnlock()
	updatedLocal := map[apps.AppID]*apps.Manifest{}
	for k, v := range local {
		if k != appID {
			updatedLocal[k] = v
		}
	}
	s.mutex.Lock()
	s.local = updatedLocal
	s.mutex.Unlock()

	updated := map[string]string{}
	for k, v := range conf.LocalManifests {
		updated[k] = v
	}
	delete(updated, string(appID))
	sc := conf.StoredConfig
	sc.LocalManifests = updated

	return s.conf.StoreConfig(sc)
}

// getFromS3 returns manifest data for an app from the S3
func (s *manifestStore) getDataFromS3(appID apps.AppID, version apps.AppVersion) ([]byte, error) {
	name := upaws.S3ManifestName(appID, version)
	data, err := s.aws.GetS3(s.s3AssetBucket, name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download manifest %s", name)
	}

	return data, nil
}

// GetFromS3 returns the manifest for an app from the S3
func (s *manifestStore) GetFromS3(appID apps.AppID, version apps.AppVersion) (*apps.Manifest, error) {
	data, err := s.getDataFromS3(appID, version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get manifest data")
	}

	m, err := apps.ManifestFromJSON(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal manifest data")
	}

	if m.AppID != appID {
		return nil, errors.New("mismatched app ID")
	}

	if m.Version != version {
		return nil, errors.New("mismatched app version")
	}

	return m, nil
}
