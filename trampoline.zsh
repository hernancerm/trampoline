# Do not source this script multiple times.
command -v tr_version > /dev/null && return

# CONFIGURATION

## Get the version.
function tr_version {
  echo '0.1.1-dev'
}

# Set by trampoline.plugin.zsh. The fallback covers sourcing this file directly.
TR_PATH="${TR_PATH:-${${(%):-%x}:A:h}}"
TR_PREFIX="${TR_PREFIX:-^t}"
TR_SET_KEYBINDS="${TR_SET_KEYBINDS:-1}"

# WIDGET

# The listing, fzf and quoting live in the core, shared with the PowerShell layer.
function tr-jump {
  local cmd
  cmd="$(bb "${TR_PATH}/core/trampoline.clj" pick --shell zsh)"
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

zle -N tr-jump

# Set keybinds.

if [[ TR_SET_KEYBINDS -eq 1 ]]; then
  bindkey "${TR_PREFIX}" tr-jump
fi
