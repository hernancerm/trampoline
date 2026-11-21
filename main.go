// Shell-independent core of trampoline. The shell layers (shell/), printed by `init`, only bind
// a key, run `trampoline pick` and execute the command it prints.
//
// Usage:
//
//	trampoline init <zsh|pwsh>       Print the shell layer, to eval from the shell's rc file.
//	trampoline --version             Print the version.
//	trampoline items [--type <d|f>]  Print the configured items, one per line, optionally only
//	                                 the directories (d) or the files (f).
//	trampoline pick --shell <shell>  Pick an item in fzf and print the command to run it, quoted
//	                                 for <shell>: zsh or pwsh. Prints nothing when the pick is
//	                                 cancelled.
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Set by the release build with `-ldflags -X main.version=<version>`.
var version = "dev"

var (
	//go:embed shell/trampoline.zsh
	zshLayer string
	//go:embed shell/trampoline.ps1
	pwshLayer string
)

// item is one line of the picker. display keeps the line as written in the config (e.g.
// `${HISTFILE}`), while path is the expanded filepath.
type item struct {
	display string
	path    string
}

// home prefers `$HOME`, as the shells use it for `~`.
func home() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	h, _ := os.UserHomeDir()
	return h
}

func configFiles() []string {
	return []string{filepath.Join(home(), ".trampoline"), filepath.Join(home(), ".trampoline.local")}
}

var envVar = regexp.MustCompile(`\$\{(\w+)\}|\$(\w+)`)

// expand expands a leading `~` and the environment variables `$VAR` and `${VAR}`. Not
// os.ExpandEnv, which also expands specials like `$1` and `$@`.
func expand(s string) string {
	if s == "~" || strings.HasPrefix(s, "~/") || strings.HasPrefix(s, `~\`) {
		s = home() + s[1:]
	}
	return envVar.ReplaceAllStringFunc(s, func(m string) string {
		return os.Getenv(strings.Trim(m, "${}"))
	})
}

// shorten replaces the home directory prefix with `~`, for display.
func shorten(s string) string {
	if h := home(); strings.HasPrefix(s, h) {
		return "~" + s[len(h):]
	}
	return s
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func configLines() []string {
	var lines []string
	for _, f := range configFiles() {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
			if strings.TrimSpace(line) != "" {
				lines = append(lines, line)
			}
		}
	}
	return lines
}

var nonAlnum = regexp.MustCompile(`[^A-Za-z0-9]`)

// subDirs lists the level-1 sub-dirs, hidden ones included. Sorted ignoring case and
// punctuation, like `sort` in an English locale: `a`, `.b`, `C`.
func subDirs(dir string) []string {
	entries, _ := os.ReadDir(dir)
	var dirs []string
	for _, e := range entries {
		// Stat, not e.IsDir(), so symlinks to dirs count.
		if p := filepath.Join(dir, e.Name()); isDir(p) {
			dirs = append(dirs, p)
		}
	}
	key := func(s string) string { return strings.ToLower(nonAlnum.ReplaceAllString(s, "")) }
	sort.Slice(dirs, func(i, j int) bool {
		ki, kj := key(dirs[i]), key(dirs[j])
		if ki != kj {
			return ki < kj
		}
		return dirs[i] < dirs[j]
	})
	return dirs
}

// items lists the items of the config lines.
func items(lines []string) []item {
	var out []item
	for _, line := range lines {
		// Only a trailing `:0` is the suffix, so a Windows drive (`C:\`) is not taken for one.
		noExpand := strings.HasSuffix(line, ":0")
		raw := strings.TrimSuffix(line, ":0")
		path := filepath.Clean(expand(raw))
		switch {
		case isFile(path), isDir(path) && noExpand:
			out = append(out, item{shorten(raw), path})
		case isDir(path):
			for _, d := range subDirs(path) {
				out = append(out, item{shorten(d), d})
			}
		}
	}
	return out
}

var (
	zshSafe  = regexp.MustCompile(`^[\w./~@%+=:,-]+$`)
	pwshSafe = regexp.MustCompile(`^[\w./\\:-]+$`)
)

func quoteZsh(s string) string {
	if zshSafe.MatchString(s) {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func quotePwsh(s string) string {
	if pwshSafe.MatchString(s) {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func pick(shell string) error {
	all := items(configLines())
	displayToPath := map[string]string{}
	lines := make([]string, len(all))
	for i, it := range all {
		// The first item wins when two display the same.
		if _, ok := displayToPath[it.display]; !ok {
			displayToPath[it.display] = it.path
		}
		lines[i] = it.display
	}
	cmd := exec.Command("fzf")
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n"))
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	// fzf exits non-zero when cancelled, which is not an error here.
	var exitErr *exec.ExitError
	if err != nil && !errors.As(err, &exitErr) {
		return err
	}
	path, ok := displayToPath[strings.TrimRight(string(out), "\r\n")]
	if !ok {
		return nil
	}
	quote := quoteZsh
	if shell == "pwsh" {
		quote = quotePwsh
	}
	run := "cd"
	if !isDir(path) {
		run = os.Getenv("EDITOR")
		if run == "" {
			run = "vim"
		}
	}
	fmt.Println(run + " " + quote(path))
	return nil
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "ERROR: "+msg)
	os.Exit(1)
}

func main() {
	args := append(os.Args[1:], "", "", "")
	cmd, flag, value := args[0], args[1], args[2]
	// Only listing needs a config, so `init` works in an rc file before the config exists.
	if (cmd == "items" || cmd == "pick") && !isFile(configFiles()[0]) && !isFile(configFiles()[1]) {
		fail("No config file found for trampoline, no ~/.trampoline or ~/.trampoline.local")
	}
	switch cmd {
	case "--version":
		fmt.Println(version)
	case "init":
		switch flag {
		case "zsh":
			fmt.Print(zshLayer)
		case "pwsh":
			fmt.Print(pwshLayer)
		default:
			fail("Usage: trampoline init <zsh|pwsh>")
		}
	case "items":
		if flag != "" && !(flag == "--type" && (value == "d" || value == "f")) {
			fail("Usage: trampoline items [--type <d|f>]")
		}
		for _, it := range items(configLines()) {
			if flag == "" || (value == "d") == isDir(it.path) {
				fmt.Println(it.display)
			}
		}
	case "pick":
		if flag != "--shell" || (value != "zsh" && value != "pwsh") {
			fail("Usage: trampoline pick --shell <zsh|pwsh>")
		}
		if err := pick(value); err != nil {
			fail(err.Error())
		}
	default:
		fail("Usage: trampoline <init|items|pick|--version>")
	}
}
