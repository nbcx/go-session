package session

import (
	"net/http"
	"path/filepath"
)

var manager *Manager

func init() {
	conf := new(Config)
	conf.CookieName = "_nb_session_id"
	conf.EnableSetCookie = true
	conf.GcLifetime = 3600
	conf.Secure = false
	conf.CookieLifeTime = 0
	conf.ProviderConfig = filepath.ToSlash("")
	conf.DisableHTTPOnly = false
	conf.Domain = ""
	conf.EnableSidInHTTPHeader = false
	conf.SessionNameInHTTPHeader = "_nb_session_id"
	conf.EnableSidInURLQuery = false
	conf.CookieSameSite = 1
	conf.SessionIDPrefix = ""
	Register("memory", mempder)
	if err := Set("memory", conf); err != nil {
		panic(err)
	}
}

// Set default session
func Set(provideName string, config *Config) (err error) {
	if manager != nil {
		manager.Destroy()
	}
	manager, err = New(provideName, config)
	return
}

// New session manager
func New(provideName string, config *Config) (*Manager, error) {
	manager, err := NewManager(provideName, config)
	if err == nil {
		go manager.GC()
	}
	return manager, err
}

func Default() *Manager {
	return manager
}

func Start(w http.ResponseWriter, r *http.Request) (session Store, err error) {
	return manager.SessionStart(w, r)
}

func Destroy(w http.ResponseWriter, r *http.Request) error {
	return manager.SessionDestroy(w, r)
}

func RegenerateID(w http.ResponseWriter, r *http.Request) (Store, error) {
	return manager.SessionRegenerateID(w, r)
}
