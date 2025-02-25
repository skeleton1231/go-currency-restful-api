// Copyright 2020 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authzserver

import (
	"github.com/pkg/errors"
	"github.com/skeleton1231/go-iam-ecommerce-microservice/internal/authzserver/load/cache"
	"github.com/skeleton1231/go-iam-ecommerce-microservice/internal/pkg/middleware"
	"github.com/skeleton1231/go-iam-ecommerce-microservice/internal/pkg/middleware/auth"
)

func newCacheAuth() middleware.AuthStrategy {
	return auth.NewCacheStrategy(getSecretFunc())
}

func getSecretFunc() func(string) (auth.Secret, error) {
	return func(kid string) (auth.Secret, error) {
		cli, err := cache.GetCacheInsOr(nil)
		if err != nil || cli == nil {
			return auth.Secret{}, errors.Wrap(err, "get cache instance failed")
		}

		secret, err := cli.GetSecret(kid)
		if err != nil {
			return auth.Secret{}, err
		}

		return auth.Secret{
			Username: secret.Username,
			ID:       secret.SecretId,
			Key:      secret.SecretKey,
			Expires:  secret.Expires,
		}, nil
	}
}
