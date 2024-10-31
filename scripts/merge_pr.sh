#!/usr/bin/env bash

# This script publishes a draft release

# standard bash error handling
set -o nounset  # treat unset variables as an error and exit immediately.
set -o errexit  # exit immediately when a command fails.
set -E          # needs to be set if we want the ERR trap
set -o pipefail # prevents errors in a pipeline from being masked

GITHUB_URL=https://api.github.com/repos/KsaweryZietara/kyma-environment-broker
GITHUB_AUTH_HEADER="Authorization: Bearer ${GH_TOKEN}"

CURL_RESPONSE=$(curl -L \
  -X PUT \
  -H "Accept: application/vnd.github+json" \
  -H "${GITHUB_AUTH_HEADER}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  ${GITHUB_URL}/pulls/${PR_NUMBER}/merge \
  -d '{"commit_title":"Bump","commit_message":"Bump sec-scanners-config.yaml, KEB images and Chart"}')
