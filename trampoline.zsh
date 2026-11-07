# Do not source this script multiple times.
command -v zt_version > /dev/null && return

# CONFIGURATION

## Get the version.
function zt_version {
  echo '0.1.1-dev'
}

# Set by trampoline.plugin.zsh. The fallback covers sourcing this file directly.
ZT_PATH="${ZT_PATH:-${${(%):-%x}:A:h}}"
ZT_PREFIX="${ZT_PREFIX:-^t}"
ZT_SET_KEYBINDS="${ZT_SET_KEYBINDS:-1}"

# WIDGET

# The listing, fzf and quoting live in the core, shared with the PowerShell layer.
function zt-jump {
  local cmd
  cmd="$(bb "${ZT_PATH}/core/trampoline.clj" pick --shell zsh)"
  if [[ $? -ne 0 ]]; then
    # The core printed the error to stderr. Accept the line so the error stays visible above it.
    zle .accept-line
    return 1
  fi
  if [[ -z "${cmd}" ]]; then
    zle .reset-prompt
    return
  fi
  BUFFER="${cmd}"
  zle .reset-prompt
  zle .accept-line
}

zle -N zt-jump

# Set keybinds.

if [[ ZT_SET_KEYBINDS -eq 1 ]]; then
  bindkey "${ZT_PREFIX}" zt-jump
fi
