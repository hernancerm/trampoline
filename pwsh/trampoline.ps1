# PowerShell layer of trampoline. Dot-source it from $PROFILE:
#   . ~/dev/repos/trampoline/pwsh/trampoline.ps1
# Set $TR_SET_KEYBINDS = 0 before dot-sourcing to bind Invoke-TrJump yourself, and $TR_PREFIX to
# change the key (PSReadLine chord syntax, see `Get-Help Set-PSReadLineKeyHandler`).

$script:TrCore = Join-Path $PSScriptRoot '../core/trampoline.clj'

# The listing, fzf and quoting live in the core, shared with the zsh layer.
function Invoke-TrJump {
  $cmd = bb $script:TrCore pick --shell pwsh
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
