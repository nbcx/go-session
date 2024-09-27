package session

import "net/http"

// ManagerConfig define the session config
type Config struct {
	EnableSetCookie         bool          `json:"enableSetCookie,omitempty"`
	DisableHTTPOnly         bool          `json:"disableHTTPOnly"`
	Secure                  bool          `json:"secure"`
	EnableSidInHTTPHeader   bool          `json:"EnableSidInHTTPHeader"`
	EnableSidInURLQuery     bool          `json:"EnableSidInURLQuery"`
	CookieName              string        `json:"cookieName"`
	Gclifetime              int64         `json:"gclifetime"`
	Maxlifetime             int64         `json:"maxLifetime"`
	CookieLifeTime          int           `json:"cookieLifeTime"`
	ProviderConfig          string        `json:"providerConfig"`
	Domain                  string        `json:"domain"`
	SessionIDLength         int64         `json:"sessionIDLength"`
	SessionNameInHTTPHeader string        `json:"SessionNameInHTTPHeader"`
	SessionIDPrefix         string        `json:"sessionIDPrefix"`
	CookieSameSite          http.SameSite `json:"cookieSameSite"`
}

func (c *Config) Opts(opts ...ManagerConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

type ManagerConfigOpt func(config *Config)

func NewManagerConfig(opts ...ManagerConfigOpt) *Config {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

// CfgCookieName set key of session id
func CfgCookieName(cookieName string) ManagerConfigOpt {
	return func(config *Config) {
		config.CookieName = cookieName
	}
}

// CfgSessionIdLength set len of session id
func CfgSessionIdLength(length int64) ManagerConfigOpt {
	return func(config *Config) {
		config.SessionIDLength = length
	}
}

// CfgSessionIdPrefix set prefix of session id
func CfgSessionIdPrefix(prefix string) ManagerConfigOpt {
	return func(config *Config) {
		config.SessionIDPrefix = prefix
	}
}

// CfgSetCookie whether set `Set-Cookie` header in HTTP response
func CfgSetCookie(enable bool) ManagerConfigOpt {
	return func(config *Config) {
		config.EnableSetCookie = enable
	}
}

// CfgGcLifeTime set session gc lift time
func CfgGcLifeTime(lifeTime int64) ManagerConfigOpt {
	return func(config *Config) {
		config.Gclifetime = lifeTime
	}
}

// CfgMaxLifeTime set session lift time
func CfgMaxLifeTime(lifeTime int64) ManagerConfigOpt {
	return func(config *Config) {
		config.Maxlifetime = lifeTime
	}
}

// CfgCookieLifeTime set cookie lift time
func CfgCookieLifeTime(lifeTime int) ManagerConfigOpt {
	return func(config *Config) {
		config.CookieLifeTime = lifeTime
	}
}

// CfgProviderConfig configure session provider
func CfgProviderConfig(providerConfig string) ManagerConfigOpt {
	return func(config *Config) {
		config.ProviderConfig = providerConfig
	}
}

// CfgDomain set cookie domain
func CfgDomain(domain string) ManagerConfigOpt {
	return func(config *Config) {
		config.Domain = domain
	}
}

// CfgSessionIdInHTTPHeader enable session id in http header
func CfgSessionIdInHTTPHeader(enable bool) ManagerConfigOpt {
	return func(config *Config) {
		config.EnableSidInHTTPHeader = enable
	}
}

// CfgSetSessionNameInHTTPHeader set key of session id in http header
func CfgSetSessionNameInHTTPHeader(name string) ManagerConfigOpt {
	return func(config *Config) {
		config.SessionNameInHTTPHeader = name
	}
}

// EnableSidInURLQuery enable session id in query string
func CfgEnableSidInURLQuery(enable bool) ManagerConfigOpt {
	return func(config *Config) {
		config.EnableSidInURLQuery = enable
	}
}

// DisableHTTPOnly set HTTPOnly for http.Cookie
func CfgHTTPOnly(HTTPOnly bool) ManagerConfigOpt {
	return func(config *Config) {
		config.DisableHTTPOnly = !HTTPOnly
	}
}

// CfgSecure set Secure for http.Cookie
func CfgSecure(Enable bool) ManagerConfigOpt {
	return func(config *Config) {
		config.Secure = Enable
	}
}

// CfgSameSite set http.SameSite
func CfgSameSite(sameSite http.SameSite) ManagerConfigOpt {
	return func(config *Config) {
		config.CookieSameSite = sameSite
	}
}
