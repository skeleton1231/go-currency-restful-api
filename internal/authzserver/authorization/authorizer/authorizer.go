package authorizer

import (
	"github.com/ory/ladon"
	"github.com/skeleton1231/go-gin-restful-api-boilerplate/internal/authzserver/authorization"
)

// PolicyGetter defines function to get policy for a given user.
type PolicyGetter interface {
	GetPolicy(key string) ([]*ladon.DefaultPolicy, error)
}

// Authorization implements authorization.AuthorizationInterface interface.
type Authorization struct {
	getter PolicyGetter
}

// NewAuthorization create a new Authorization instance.
func NewAuthorization(getter PolicyGetter) authorization.AuthorizationInterface {
	return &Authorization{getter}
}

// Create create a policy.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) Create(policy *ladon.DefaultPolicy) error {
	return nil
}

// Update update a policy.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) Update(policy *ladon.DefaultPolicy) error {
	return nil
}

// Delete delete a policy by the given identifier.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) Delete(id string) error {
	return nil
}

// DeleteCollection batch delete policies by the given identifiers.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) DeleteCollection(idList []string) error {
	return nil
}

// Get returns the policy detail by the given identifier.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) Get(id string) (*ladon.DefaultPolicy, error) {
	return &ladon.DefaultPolicy{}, nil
}

// List returns all the policies under the username.
func (auth *Authorization) List(username string) ([]*ladon.DefaultPolicy, error) {
	return auth.getter.GetPolicy(username)
}

// LogRejectedAccessRequest write rejected subject access to redis.
func (auth *Authorization) LogRejectedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {

}

// LogGrantedAccessRequest write granted subject access to redis.
func (auth *Authorization) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {

}
