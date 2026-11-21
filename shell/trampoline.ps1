# PowerShell layer of trampoline, printed by `trampoline init pwsh`. Evaluate it from $PROFILE:
#   trampoline init pwsh | Out-String | Invoke-Expression
# Set $TR_SET_KEYBINDS = 0 before that to bind Invoke-TrJump yourself, and $TR_PREFIX to
# change the key (PSReadLine chord syntax, see `Get-Help Set-PSReadLineKeyHandler`).

# The listing, fzf and quoting live in the core, shared with the zsh layer.
function Invoke-TrJump {
  $cmd = trampoline pick --shell pwsh
  if ($LASTEXITCODE -ne 0 -or -not $cmd) {
    # On error the core printed it to stderr. Redraw the prompt below it.
    [Microsoft.PowerShell.PSConsoleReadLine]::InvokePrompt()
    return
  }
  [Microsoft.PowerShell.PSConsoleReadLine]::RevertLine()
  [Microsoft.PowerShell.PSConsoleReadLine]::Insert($cmd)
  [Microsoft.PowerShell.PSConsoleReadLine]::AcceptLine()
}

if ($null -eq $TR_PREFIX) {
  $TR_PREFIX = 'Ctrl+t'
}
if ($TR_SET_KEYBINDS -ne 0) {
  Set-PSReadLineKeyHandler -Chord $TR_PREFIX -ScriptBlock { Invoke-TrJump }
}
