#!/usr/bin/env bash

# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

repos=$(az acr repository list -o tsv --name capzci)

for repo in $repos
do
  filters+="--filter $repo:.* "
done

PURGE_CMD="acr purge $filters --ago 1d --untagged"

az acr task create --name midnight_capz_purge --cmd "${PURGE_CMD}" \
  --schedule "0 0 * * *" --registry capzci --context /dev/null
