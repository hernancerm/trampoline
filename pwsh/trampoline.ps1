# PowerShell layer of zsh-trampoline. Dot-source it from $PROFILE:
#   . ~/dev/repos/zsh-trampoline/pwsh/trampoline.ps1
# Set $ZT_SET_KEYBINDS = 0 before dot-sourcing to bind Invoke-ZtJump yourself, and $ZT_PREFIX to
# change the key (PSReadLine chord syntax, see `Get-Help Set-PSReadLineKeyHandler`).

$script:ZtCore = Join-Path $PSScriptRoot '../core/trampoline.clj'

# The listing, fzf and quoting live in the core, shared with the zsh layer.
function Invoke-ZtJump {
  $cmd = bb $script:ZtCore pick --shell pwsh
  if ($LASTEXITCODE -ne 0 -or -not $cmd) {
    # On error the core printed it to stderr. Redraw the prompt below it.
    [Microsoft.PowerShell.PSConsoleReadLine]::InvokePrompt()
    return
  }
  [Microsoft.PowerShell.PSConsoleReadLine]::RevertLine()
  [Microsoft.PowerShell.PSConsoleReadLine]::Insert($cmd)
  [Microsoft.PowerShell.PSConsoleReadLine]::AcceptLine()
}

if ($null -eq $ZT_PREFIX) {
  $ZT_PREFIX = 'Ctrl+t'
}
if ($ZT_SET_KEYBINDS -ne 0) {
  Set-PSReadLineKeyHandler -Chord $ZT_PREFIX -ScriptBlock { Invoke-ZtJump }
}
