# zsh layer of trampoline, printed by `trampoline init zsh`. Eval it from ~/.zshrc:
#   eval "$(trampoline init zsh)"
# No guard against a second eval: a `return` here would also end the ~/.zshrc evaling it, and
# running it twice only rebinds the same widget.

# CONFIGURATION

TR_PREFIX="${TR_PREFIX:-^t}"
TR_SET_KEYBINDS="${TR_SET_KEYBINDS:-1}"

# WIDGET

# The listing, fzf and quoting live in the core, shared with the PowerShell layer.
function tr-jump {
  local cmd
  cmd="$(trampoline pick --shell zsh)"
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
