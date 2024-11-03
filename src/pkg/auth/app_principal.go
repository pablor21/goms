package auth

type DflPrincipal struct {
	ID     interface{}            `json:"id"`
	Roles  []string               `json:"roles"`
	Claims map[string]interface{} `json:"claims"`
}

func NewPrincipal() Principal {
	return &DflPrincipal{}
}

func (p *DflPrincipal) GetID() interface{} {
	return p.ID
}

func (p *DflPrincipal) SetId(id interface{}) Principal {
	p.ID = id
	return p
}

func (p *DflPrincipal) GetRoles() []string {
	return p.Roles
}

func (p *DflPrincipal) SetRoles(roles []string) Principal {
	p.Roles = roles
	return p
}

func (p *DflPrincipal) SetClaims(claims map[string]interface{}) Principal {
	p.Claims = claims
	return p
}

func (p *DflPrincipal) GetClaims() map[string]interface{} {
	return p.Claims
}

func (p *DflPrincipal) GetClaim(key string) (interface{}, bool) {
	value, ok := p.Claims[key]
	return value, ok
}

func (p *DflPrincipal) SetClaim(key string, value interface{}) Principal {
	p.Claims[key] = value
	return p
}
