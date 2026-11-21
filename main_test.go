package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func eq[T comparable](t *testing.T, want, got T) {
	t.Helper()
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestExpand(t *testing.T) {
	h := home()
	eq(t, h, expand("~"))
	eq(t, h+"/a", expand("~/a"))
	eq(t, "~a", expand("~a"))
	eq(t, "a/~", expand("a/~"))
	eq(t, h+"/a", expand("$HOME/a"))
	eq(t, h+"/a", expand("${HOME}/a"))
	eq(t, "/a", expand("$TRAMPOLINE_TEST_UNSET/a"))
}

func TestShorten(t *testing.T) {
	eq(t, "~/a", shorten(home()+"/a"))
	eq(t, "/elsewhere", shorten("/elsewhere"))
}

func TestQuoteZsh(t *testing.T) {
	eq(t, "/a/b-c", quoteZsh("/a/b-c"))
	eq(t, "'/a b'", quoteZsh("/a b"))
	eq(t, `'/it'\''s'`, quoteZsh("/it's"))
}

func TestQuotePwsh(t *testing.T) {
	eq(t, `C:\a\b`, quotePwsh(`C:\a\b`))
	eq(t, `'C:\a b'`, quotePwsh(`C:\a b`))
	eq(t, "'it''s'", quotePwsh("it's"))
}

// tree makes a temp dir holding the dirs `a`, `.b`, `C` and `with space`, and the file `f.txt`.
// It returns the dir and its sub-dirs in the expected order.
func tree(t *testing.T) (string, []string) {
	// EvalSymlinks as macOS temp dirs sit behind the /var symlink.
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	var subs []string
	for _, d := range []string{"a", ".b", "C", "with space"} {
		p := filepath.Join(dir, d)
		if err := os.Mkdir(p, 0o755); err != nil {
			t.Fatal(err)
		}
		subs = append(subs, p)
	}
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	return dir, subs
}

func paths(its []item) []string {
	var out []string
	for _, it := range its {
		out = append(out, it.path)
	}
	return out
}

func TestSubDirs(t *testing.T) {
	dir, subs := tree(t)
	if got := subDirs(dir); !slices.Equal(subs, got) {
		t.Errorf("want %v, got %v", subs, got)
	}
}

func TestItems(t *testing.T) {
	dir, subs := tree(t)
	file := filepath.Join(dir, "f.txt")
	t.Run("a dir with /* lists its sub-dirs", func(t *testing.T) {
		if got := paths(items([]string{dir + "/*"})); !slices.Equal(subs, got) {
			t.Errorf("want %v, got %v", subs, got)
		}
	})
	t.Run(`a dir with \* lists its sub-dirs`, func(t *testing.T) {
		if got := paths(items([]string{dir + `\*`})); !slices.Equal(subs, got) {
			t.Errorf("want %v, got %v", subs, got)
		}
	})
	t.Run("a dir lists itself, displayed as written", func(t *testing.T) {
		got := items([]string{dir})
		if !slices.Equal([]item{{dir, dir}}, got) {
			t.Errorf("got %v", got)
		}
	})
	t.Run("a file lists itself", func(t *testing.T) {
		got := items([]string{file})
		if !slices.Equal([]item{{file, file}}, got) {
			t.Errorf("got %v", got)
		}
	})
	t.Run("a missing path lists nothing", func(t *testing.T) {
		missing := filepath.Join(dir, "missing")
		if got := items([]string{missing, missing + "/*"}); len(got) != 0 {
			t.Errorf("got %v", got)
		}
	})
	t.Run("lines are listed in order", func(t *testing.T) {
		got := paths(items([]string{file, dir}))
		if !slices.Equal([]string{file, dir}, got) {
			t.Errorf("got %v", got)
		}
	})
}
