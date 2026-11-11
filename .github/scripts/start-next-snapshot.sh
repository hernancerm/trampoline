#!/usr/bin/env bash

# Returns the working tree to a dev version after a release, so builds off
# main are never mistaken for the released version. Run from the repo root.

set -euo pipefail

version="${1:?usage: start-next-snapshot.sh <released version>}"

IFS=. read -r major minor patch <<< "${version}"
next="${major}.${minor}.$((patch + 1))-dev"

.github/scripts/set-version.sh "${next}"

echo "${next}"
