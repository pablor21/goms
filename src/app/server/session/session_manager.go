package session

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/pablor21/goms/app/config"
	server_config "github.com/pablor21/goms/app/server/config"
	"github.com/pablor21/goms/app/server/session/stores/filesystem"
)

type sessionKeyType string

var sessionKey sessionKeyType = "session"

type SessionManager struct {
	stores map[string]*scs.SessionManager
}

var sessionManager *SessionManager

func GetSessionManager() *SessionManager {
	if sessionManager == nil {
		sessionManager = &SessionManager{
			stores: make(map[string]*scs.SessionManager),
		}
	}
	return sessionManager
}

// GetStore returns a session store by name
// The store config is defined in the config file
func (sm *SessionManager) GetStore(name string) *scs.SessionManager {
	var cfg server_config.SessionConfig
	var ok bool
	if cfg, ok = config.GetConfig().Server.Auth[name]; !ok {
		panic("session store not found")
	}

	if _, ok := sm.stores[name]; !ok {
		m := scs.New()
		m.Lifetime = time.Duration(cfg.Lifetime) * time.Second
		m.Store = sm.initSessionStore(cfg, m)

		var sameSite http.SameSite = http.SameSiteDefaultMode

		switch cfg.Cookie.SameSite {
		case "lax":
			sameSite = http.SameSiteLaxMode
		case "strict":
			sameSite = http.SameSiteStrictMode
		case "none":
			sameSite = http.SameSiteNoneMode
		}
		m.Cookie.Name = cfg.Cookie.Name
		m.Cookie.Path = cfg.Cookie.Path
		m.Cookie.Domain = cfg.Cookie.Domain
		m.Cookie.Secure = cfg.Cookie.Secure
		m.Cookie.HttpOnly = cfg.Cookie.HttpOnly
		m.Cookie.SameSite = sameSite
		m.Cookie.Persist = true
		m.Lifetime = time.Duration(cfg.Cookie.MaxAge) * time.Second

		sm.stores[name] = m
	}
	return sm.stores[name]
}

func (sm *SessionManager) initSessionStore(cfg server_config.SessionConfig, s *scs.SessionManager) scs.Store {

	switch cfg.Store.Type {
	case "memory":
		return memstore.New()
	case "filesystem":
		return filesystem.New(cfg.Store.Config, s)
	default:
		panic(fmt.Sprintf("session store type %s not supported", cfg.Store.Type))
	}
}

func SetSession(ctx context.Context, s *scs.SessionManager) context.Context {
	return context.WithValue(ctx, sessionKey, NewSession(ctx, s))
}

func GetSession(ctx context.Context) *Session {
	if s, ok := ctx.Value(sessionKey).(*Session); ok {
		return s
	}
	return nil
}
