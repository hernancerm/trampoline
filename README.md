# trampoline

zsh (or PowerShell) + [fzf](https://github.com/junegunn/fzf) for filesystem navigation.

## What is this?

This is a zsh and PowerShell plugin to facilitate two use cases:

1. `cd` to commonly visited directories.
2. Open in `$EDITOR` commonly visited files.

After sourcing the plugin this keybind is set:

- <kbd>Ctrl-t</kbd> List [configured](#configuration-file) directories and files in fzf.

## Installation

### Without a plugin manager

1. Install [fzf](https://github.com/junegunn/fzf) and [babashka](https://babashka.org/)
   (`bb`).

2. Clone the trampoline repository by executing the below command:

    ```bash
    git clone 'https://github.com/hernancerm/trampoline.git' \
      "${HOME}/.local/share/trampoline"
    ```

#### Zsh

3. Source the plugin (also add command to your `~/.zshrc` to enable on future sessions):

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

#### PowerShell

3. Source the PowerShell layer from `$PROFILE`:

    ```powershell
    . "${HOME}/.local/share/trampoline/pwsh/trampoline.ps1"
    ```
    
    OR: If you want to set a custom keybind then use:
    
    ```powershell
    $TR_SET_KEYBINDS = 0
    . "${HOME}/.local/share/trampoline/pwsh/trampoline.ps1"
    Set-PSReadLineKeyHandler -Chord 'Ctrl+t' -ScriptBlock { Invoke-TrJump }
    ```

### With a plugin manager

#### Zsh

With [Sheldon](https://github.com/rossmacarthur/sheldon):

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
- Environment variables are supported, as `$VAR` or `${VAR}`. Only exported ones are seen.

## Similar projects

- <https://github.com/ajeetdsouza/zoxide>
- <https://github.com/wting/autojump>
- <https://github.com/agkozak/zsh-z>
