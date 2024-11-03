package session

import (
	"context"

	"github.com/alexedwards/scs/v2"
	"github.com/pablor21/goms/pkg/auth"
	"github.com/pablor21/goms/pkg/errors"
)

type Session struct {
	ctx     context.Context
	manager *scs.SessionManager
}

func NewSession(ctx context.Context, manager *scs.SessionManager) *Session {
	return &Session{
		ctx:     ctx,
		manager: manager,
	}
}

func (s *Session) Put(key string, value interface{}) {
	s.manager.Put(s.ctx, key, value)
}

func (s *Session) GetString(key string) string {
	return s.manager.GetString(s.ctx, key)
}

func (s *Session) GetInt(key string) int {
	return s.manager.GetInt(s.ctx, key)
}

func (s *Session) GetInt64(key string) int64 {
	return s.manager.GetInt64(s.ctx, key)
}

func (s *Session) GetFloat(key string) float64 {
	return s.manager.GetFloat(s.ctx, key)
}

func (s *Session) GetBool(key string) bool {
	return s.manager.GetBool(s.ctx, key)
}

func (s *Session) Get(key string) interface{} {
	return s.manager.Get(s.ctx, key)
}

func (s *Session) Remove(key string) {
	s.manager.Remove(s.ctx, key)
}

func (s *Session) Clear() {
	s.manager.Clear(s.ctx)
}

func (s *Session) Destroy() {
	s.manager.Destroy(s.ctx)
}

func (s *Session) RenewToken() error {
	return s.manager.RenewToken(s.ctx)
}

func (s *Session) SetPrincipal(principal auth.Principal) {
	s.Put(string(auth.PrincipalKeyContextKey), principal)
}

func (s *Session) GetPrincipal() (res auth.Principal, err error) {
	rawPrincipal := s.Get(string(auth.PrincipalKeyContextKey))
	if principal, ok := rawPrincipal.(auth.DflPrincipal); ok {
		res = &principal
	} else {
		err = errors.ErrUnauthorized
	}
	return
}
