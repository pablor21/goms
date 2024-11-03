package auth

const (
	ClaimEmail     = "email"
	ClaimPhone     = "phone"
	ClaimFirstName = "first_name"
	ClaimLastName  = "last_name"
	ClaimLang      = "lang"
	ClaimRole      = "role"
	ClaimAvatar    = "avatar"
	ClaimStatus    = "status"
)

type Principal interface {
	GetID() interface{}
	SetId(id interface{}) Principal
	GetRoles() []string
	SetRoles(roles []string) Principal
	SetClaims(claims map[string]interface{}) Principal
	GetClaims() map[string]interface{}
	GetClaim(key string) (interface{}, bool)
	SetClaim(key string, value interface{}) Principal
}
