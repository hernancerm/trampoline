<a href="https://github.com/hernancerm/trampoline/actions/workflows/ci.yml" target="_blank">
  <img src="https://github.com/hernancerm/trampoline/actions/workflows/ci.yml/badge.svg" />
</a>

# trampoline

zsh (or PowerShell) + [fzf](https://github.com/junegunn/fzf) for filesystem navigation.

## What is this?

This is a command for zsh and PowerShell to facilitate two use cases:

1. `cd` to commonly visited directories.
2. Open in `$EDITOR` commonly visited files.

After the [shell setup](#shell-setup) this keybind is set:

- <kbd>Ctrl-t</kbd> List [configured](#configuration-file) directories and files in fzf.

## Installation

Install [fzf](https://github.com/junegunn/fzf) and the `trampoline` binary. With
[mise](https://mise.jdx.dev/):

```bash
mise use -g github:hernancerm/trampoline
```

OR: Download it from [releases](https://github.com/hernancerm/trampoline/releases) and add
to `$PATH`.

## Shell setup

### Zsh

Add to `~/.zshrc`:

```bash
eval "$(trampoline init zsh)"
```

OR: If you want to set a custom keybind then use:

```bash
TR_SET_KEYBINDS=0
eval "$(trampoline init zsh)"
bindkey "^t" tr-jump
```

OR: If you use [jeffreytse/zsh-vi-mode](https://github.com/jeffreytse/zsh-vi-mode):

```bash
TR_SET_KEYBINDS=0
eval "$(trampoline init zsh)"
function zvm_after_init {
  zvm_define_widget tr-jump
  zvm_bindkey viins "^t" tr-jump
}
```

### PowerShell

Add to `$PROFILE`:

```powershell
trampoline init pwsh | Out-String | Invoke-Expression
```

OR: If you want to set a custom keybind then use:

```powershell
$TR_SET_KEYBINDS = 0
trampoline init pwsh | Out-String | Invoke-Expression
Set-PSReadLineKeyHandler -Chord 'Ctrl+t' -ScriptBlock { Invoke-TrJump }
```

## Configuration file

trampoline supports two config files:

1. `~/.trampoline`
2. `~/.trampoline.local` (for what you don't want to commit to source control.)

The configs are merged in order, no overrides. I recommend using the first config file for
common dirs/files in your setup, while the second config file for sensitive or
machine-specific dirs/files. Both files share the same syntax.

Example contents of the config file:

```text
~/dev/work/*
~/dev/temp
~/dev/Some Dir With Spaces
~/dev/Another Dir With Spaces/*
${HISTFILE}
$HOME/foo/*
```

Explanation of the config file's syntax:

- Each line has the absolute path of one directory or one file, listed as is.
- A directory ending in `/*` lists its level-1 sub-dirs instead, hidden ones included.
  On Windows `\*` works too.
- The tilde (`~`) character is expanded to the user home directory.
- Exported env variables are supported, as `$VAR` or `${VAR}`.

## Development

Run the tests with `go test ./...`. To try a change, build the binary onto `$PATH`, ahead
of any installed release:

```bash
go build -o ~/.local/bin/trampoline .
```

## Similar projects

- <https://github.com/ajeetdsouza/zoxide>
- <https://github.com/wting/autojump>
- <https://github.com/agkozak/zsh-z>
