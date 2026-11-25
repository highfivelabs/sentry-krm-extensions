// Copyright 2026 High Five Labs
// SPDX-License-Identifier: Apache-2.0

package extensions

import (
	"net/http"

	"github.com/highfivelabs/krm-rest-runtime/rest"
)

// Vendor is how the manager reaches Sentry. Sentry hosts per data-storage
// region (us, de) and can be self-hosted, so the Connection's baseURL
// overrides the default.
var Vendor = &rest.Vendor{
	Name:           "sentry",
	DefaultBaseURL: "https://us.sentry.io/api/0",
	CredentialKeys: []string{"token"},
	Authorize:      func(r *http.Request, c map[string]string) { r.Header.Set("Authorization", "Bearer "+c["token"]) },
}
