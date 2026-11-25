#!/usr/bin/env bash
# Copyright 2026 High Five Labs
# SPDX-License-Identifier: Apache-2.0
# Download the vendor's OpenAPI specs at the commits vendor.yaml pins.
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p specs
[ -s specs/sentry.json ] || curl -fsSL -o specs/sentry.json "https://raw.githubusercontent.com/getsentry/sentry-api-schema/d8e95f7f3c95b8aa5f60983f6f2239b71ac529d6/openapi-derefed.json"
