# trampoline

zsh (or PowerShell) + [fzf](https://github.com/junegunn/fzf) for filesystem navigation.

## What is this?

This is a zsh and PowerShell plugin to facilitate two use cases:

1. `cd` to commonly visited directories.
2. Open in `$EDITOR` commonly visited files.

After sourcing the plugin this keybind is set:

- <kbd>Ctrl-t</kbd> List [configured](#configuration-file) directories and files in fzf.

## Installation

First, install [fzf](https://github.com/junegunn/fzf) and the `trampoline` binary. With
[mise](https://mise.jdx.dev/):

```bash
mise use -g github:hernancerm/trampoline
```

OR: Download it from [releases](https://github.com/hernancerm/trampoline/releases) and add
to `$PATH`.

### Zsh

With [Sheldon](https://github.com/rossmacarthur/sheldon) (Zsh plugin manager):

```toml
[plugins.trampoline]
github = "hernancerm/trampoline"

# File: ~/.config/sheldon/plugins.toml
```

```bash
# Source plugins.
eval "$(sheldon source)"

# File: ~/.zshrc
```

Alternatively, if not using Sheldon:

1. Clone the trampoline repository by executing the below command:

    ```bash
    git clone 'https://github.com/hernancerm/trampoline.git' \
      "${HOME}/.local/share/trampoline"
    ```

2. Source the plugin (also add the command to your `~/.zshrc` to enable on future
   sessions):

    ```bash
    source "${HOME}/.local/share/trampoline/trampoline.plugin.zsh"
    ```

    OR: If you want to set a custom keybind then use:

    ```bash
    TR_SET_KEYBINDS=0
    source "${HOME}/.local/share/trampoline/trampoline.plugin.zsh"
    bindkey "^t" tr-jump
    ```

    OR: If you use [jeffreytse/zsh-vi-mode](https://github.com/jeffreytse/zsh-vi-mode):

    ```bash
    TR_SET_KEYBINDS=0
    source "${HOME}/.local/share/trampoline/trampoline.plugin.zsh"
    function zvm_after_init {
      zvm_define_widget tr-jump
      zvm_bindkey viins "^t" tr-jump
    }
    ```

### PowerShell

1. Clone the trampoline repository by executing the below command:

    ```powershell
    git clone 'https://github.com/hernancerm/trampoline.git' "${HOME}/.local/share/trampoline"
    ```

2. Source the PowerShell layer from `$PROFILE`:

    ```powershell
    . "${HOME}/.local/share/trampoline/pwsh/trampoline.ps1"
    ```

    OR: If you want to set a custom keybind then use:

    ```powershell
    $TR_SET_KEYBINDS = 0
    . "${HOME}/.local/share/trampoline/pwsh/trampoline.ps1"
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
~/dev/work
~/dev/temp:0
~/dev/Some Dir With Spaces:0
~/dev/Another Dir With Spaces
${HISTFILE}
$HOME/foo/
```

Explanation of the config file's syntax:

- Each line has the absolute path of one directory or one file.
- By default, directories are expanded to their level-1 sub-dirs.
- Directories may be suffixed with `:0` to list the dir itself instead of its level-1
  sub-dirs.
- Only a trailing `:0` is read as the suffix, so paths may contain a colon (`C:\Users`).
- The tilde (`~`) character is expanded to the user home directory.
- Exported env variables are supported, as `$VAR` or `${VAR}`.

## Development

To run [src/trampoline/core.clj](/src/trampoline/core.clj) from source with
[babashka](https://babashka.org/) instead of the binary, set this before sourcing the
plugin:

- Zsh: `TR_CMD=(bb ~/dev/repos/trampoline/src/trampoline/core.clj)`
- PowerShell: `$TR_CMD = 'bb', "$HOME/dev/repos/trampoline/src/trampoline/core.clj"`

The core must run on both bb and the JVM. `bb test` runs the tests on both. Build the
binary with `bb bin`.

## Similar projects

- <https://github.com/ajeetdsouza/zoxide>
- <https://github.com/wting/autojump>
- <https://github.com/agkozak/zsh-z>
