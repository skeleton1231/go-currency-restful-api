// Copyright 2020 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authorization

import (
	"github.com/marmotedu/log"
	"github.com/ory/ladon"
)

// AuditLogger outputs and cache information about granting or rejecting policies.
type AuditLogger struct {
	client AuthorizationInterface
}

// NewAuditLogger creates a AuditLogger with default parameters.
func NewAuditLogger(client AuthorizationInterface) *AuditLogger {
	return &AuditLogger{
		client: client,
	}
}

// LogRejectedAccessRequest write rejected subject access to log.
func (a *AuditLogger) LogRejectedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	a.client.LogRejectedAccessRequest(r, p, d)
	log.Debug("subject access review rejected", log.Any("request", r), log.Any("deciders", d))
}

// LogGrantedAccessRequest write granted subject access to log.
func (a *AuditLogger) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	a.client.LogGrantedAccessRequest(r, p, d)
	log.Debug("subject access review granted", log.Any("request", r), log.Any("deciders", d))
}
