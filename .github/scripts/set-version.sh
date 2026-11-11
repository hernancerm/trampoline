#!/usr/bin/env bash

# Sets the project version where it is written down: `tr_version` in trampoline.zsh.
# Run from the repo root.

set -euo pipefail

version="${1:?usage: set-version.sh <version>}"
zsh_file='trampoline.zsh'

sed -E "/^function tr_version/,/^}/ s/'[^']*'/'${version}'/" "${zsh_file}" > "${zsh_file}.tmp"
mv "${zsh_file}.tmp" "${zsh_file}"

# The sed is silent when the function is renamed or reformatted, which would
# release with the previous version. Fail instead.
if ! grep -q "echo '${version}'" "${zsh_file}"
then
  echo "error: could not set tr_version in ${zsh_file}" >&2
  exit 1
fi
