# Code Review: gico

**Reviewer:** Claude (automated)  
**Date:** 2026-02-17  
**Commit:** b165dc2  

---

## Executive Summary

gico is a well-structured Go project with a clean separation between core library (`commits/`, `types/`, `graph/`) and executables (`cmd/`). The architecture is sound — channel-based parallel repo parsing, memoization cache, pluggable rendering (SVG, terminal). However, there are several bugs (including a crash-on-panic HTTP server, incorrect leap year logic, and race conditions in the cache), zero test coverage, stale dead code, and some naming/idiom issues.

**Severity breakdown:** 🔴 3 critical | 🟡 8 moderate | 🟢 6 minor

---

## 🔴 Critical Issues

### 1. HTTP handlers panic instead of returning errors (`cmd/gico-server/main.go:33-34, 53-54, 79-80`)

Every handler calls `panic(err)` on failure. This crashes the entire server on any bad request or transient git error.

```go
// BAD: crashes server
if err != nil {
    panic(err)
}

// GOOD: return HTTP error
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}
```

Also, line 59 ignores the error from `json.Marshal`:
```go
b, _ := json.Marshal(freq)  // if this fails, writes empty body with 200
```

### 2. Incorrect leap year calculation (multiple files)

The leap year check `year%4 == 0` is wrong. Years divisible by 100 but not 400 are NOT leap years (e.g., 1900, 2100).

**Affected locations:**
- `commits/chancommits.go:58, 113-114, 156-157, 209-210`
- `commits/commits.go:40-41, 119-120`
- `ui/ui.go:70-73`
- `ui/graph.go:22-23`

**Fix:** Extract a helper and use the correct algorithm:

```go
func isLeapYear(year int) bool {
    return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func yearLength(year int) int {
    if isLeapYear(year) {
        return 366
    }
    return 365
}
```

This is duplicated 8+ times — a single `yearLength()` function in `types` or `commits` would DRY it up.

### 3. Cache goroutine leak / race condition (`commits/cache.go:84-89, 98-102, 154-160, 175-180`)

Every `CacheRepo`/`CacheGraph`/etc. spawns a goroutine that sleeps for 1 hour then deletes the entry. Problems:

1. **Goroutine leak:** If the same key is cached repeatedly (e.g., TUI refreshing), goroutines accumulate — each sleeping for an hour. In a server context with frequent requests, this is unbounded.
2. **Stale deletion:** The goroutine captures the key at spawn time. If the entry is refreshed (re-cached), the old goroutine still fires and deletes the *new* entry. The comment on line 158 ("optimization, check if the creation time has changed") acknowledges this but doesn't implement the fix.

**Fix:** Use a `time.Time` comparison before deleting, or better, use a single background reaper:

```go
func CacheRepo(path string, commits []types.Commit) {
    mapTex.Lock()
    defer mapTex.Unlock()
    repoCache[path] = types.ExpRepo{Commits: commits, Created: time.Now()}
}

// Single reaper goroutine started in init()
func startReaper() {
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        for range ticker.C {
            mapTex.Lock()
            now := time.Now()
            for k, v := range repoCache {
                if now.Sub(v.Created) > time.Hour {
                    delete(repoCache, k)
                }
            }
            mapTex.Unlock()
        }
    }()
}
```

---

## 🟡 Moderate Issues

### 4. `hashSlice` appends a newline (`commits/cache.go:42`)

```go
return fmt.Sprintf("%x\n", b)  // trailing \n in cache key
```

This works consistently (always there) but is surprising and would break if any key comparison strips whitespace. Remove the `\n`.

### 5. `OpenRepo` dereferences potentially nil pointer (`commits/common.go:36`)

```go
r, err := git.PlainOpenWithOptions(directory, &(git.PlainOpenOptions{DetectDotGit: true}))
return Repo{Repo: *r, Path: directory}, err
```

If `err != nil`, `r` may be nil, causing a nil pointer dereference *before* the caller can check `err`. Fix:

```go
r, err := git.PlainOpenWithOptions(directory, &git.PlainOpenOptions{DetectDotGit: true})
if err != nil {
    return Repo{}, err
}
return Repo{Repo: *r, Path: directory}, nil
```

### 6. `SetColorScheme` appends instead of replacing (`graph/common/common.go:30-34`)

```go
func SetColorScheme(c []color.Color) {
    for _, c := range c {
        colorScheme = append(colorScheme, sc.FromRGBA(c.RGBA()))
    }
}
```

This *appends* to the existing scheme. Callers likely expect it to *replace*. Also, the parameter `c` shadows the function parameter — use a different name:

```go
func SetColorScheme(colors []color.Color) {
    colorScheme = colorScheme[:0]
    for _, clr := range colors {
        colorScheme = append(colorScheme, sc.FromRGBA(clr.RGBA()))
    }
}
```

### 7. `colorsLoaded` sync.Once is unused everywhere

`graph/common/common.go:12`, `graph/svg/svg.go:22-23`, `graph/term/term.go:15-16` all declare `colorsLoaded sync.Once` and `colorScheme` variables that are never used. staticcheck flags these. The color init is done in `init()` instead. Remove the dead declarations.

### 8. `svgToPng` is dead code with hardcoded filenames (`graph/svg/svg.go:92-116`)

This function reads `in.svg` and writes `out.png` — clearly leftover development code. It's unused (staticcheck U1000). Remove it or move it to a `cmd/` tool if needed.

### 9. Author filter uses regex without escaping (`commits/chancommits.go:239, commits/commits.go:132`)

Author names/emails are passed directly to `regexp.Compile`. An email like `user+tag@example.com` contains `+` which is a regex quantifier. A name with `(` would cause a compile error.

```go
// Fix: escape the input or use strings.Contains
r, err := regexp.Compile(regexp.QuoteMeta(a))
```

Or if regex matching is intentional, document it clearly in the API.

### 10. `GetWeekFreq` doesn't handle year boundaries correctly (`commits/commits.go:22-33`)

When `today < 6` (first week of year), it fetches last year's freq and prepends it. But `today` is adjusted by adding 365 (or 366), then it slices `freq[today-6 : today+1]`. The logic assumes `freq` is a simple concatenation, which works, but:
- If `curYear` is a leap year, the combined slice has 366 + N entries, but indexing by `today + 365/366` could go wrong at boundaries.
- No test coverage makes this fragile.

### 11. `Frequency` (sync version) silently continues on error (`commits/commits.go:52-53`)

```go
commits, err := repo.GetCommitSet()
if err != nil {
    log.Printf("skipping repo %s\n", repo.Path)
}
// continues to use `commits` even when err != nil
```

Should `continue` after the log:

```go
if err != nil {
    log.Printf("skipping repo %s\n", repo.Path)
    continue
}
```

---

## 🟢 Minor Issues

### 12. Naming conventions

Per Go conventions and Tai's standards:

- `mapTex` → `mapMu` (it's a mutex, not a texture) — `commits/cache.go:15`
- `yst` → `yearStr` — `cmd/gico-server/main.go:44, 65`
- `y` → `parsedYear` — `cmd/gico-server/main.go:46, 68`
- `gfreq` → `globalFreq` or `mergedFreq` — `commits/commits.go:44`
- `aName`/`aEmail` → `authorName`/`authorEmail` — `cmd/mgfetch/main.go:17-18`
- `cs` → `commitSet` — `commits/commits.go:75`
- `cc` → `commitChan` — `commits/chancommits.go:19, 70`
- `a` (in cache functions) → `authorHash` — `commits/cache.go:46`
- `r` (in cache functions) → `repoHash` — `commits/cache.go:47`
- `sb` → `buf` — `graph/svg/svg.go:37`
- Single-letter `s`, `c`, `x` used as non-loop variables throughout `graph/` — expand these

### 13. `delegateKeyMap` and `newDelegateKeyMap` unused (`ui/listdelegate.go:52-62`)

Dead code. Remove or wire it up.

### 14. `highlightedEntry` field unused (`ui/settings.go:25`)

```go
highlightedEntry int  // never read or written
```

### 15. `colorsLoaded` sync.Once pattern abandoned

The `init()` approach in `graph/common/common.go:21-28` overwrites `colors` immediately (line 23 assigns, line 24 reassigns). The first palette is dead code:

```go
colors := []string{"#767960", "#a7297f", "#e8ca89", "#f5efd6", "#158266"}
colors = []string{"#0e4429", "#006d32", "#26a641", "#39d353"}  // overwrites above
```

### 16. `Content-Type` for SVG endpoints should be `image/svg+xml`

`cmd/gico-server/main.go:30, 84` set `Content-Type: text/html` for SVG responses. Should be `image/svg+xml`.

### 17. Gorilla mux/handlers are archived

`github.com/gorilla/mux` and `github.com/gorilla/handlers` are archived (Dec 2023). Consider migrating to `net/http` (Go 1.22+ has pattern matching) or `chi`.

---

## Test Coverage

**There are zero test files in the entire project.** Priority areas for testing:

1. **`types/helpers.go` — `Freq.Merge()`**: Easy to unit test, handles edge cases (different lengths)
2. **`commits/cache.go`**: Cache hit/miss/expiry logic
3. **`graph/common/common.go` — `ColorForFrequency`, `MinMax`**: Pure functions, easy to test
4. **`commits/chancommits.go` — `FilterCChanByYear`, `FilterCChanByAuthor`**: Channel pipeline correctness
5. **Leap year helper** (once extracted): Table-driven tests for 1900, 2000, 2004, 2100

---

## Dependency Hygiene

- **`go 1.19`**: Very old. The project uses generics-era dependencies (bubbletea v0.25). Consider bumping to at least `go 1.21`.
- **Gorilla packages archived**: `gorilla/mux` v1.8.0 and `gorilla/handlers` v1.5.1 — see item 17.
- **`srwiley/oksvg` and `srwiley/rasterx`**: Only used in dead `svgToPng()` function. Removing that function drops two dependencies.
- **`golang.org/x/image`**: Also only pulled in by the dead SVG-to-PNG code.
- **`imdario/mergo` + `dario.cat/mergo`**: Both in go.sum — transitive dep duplication from go-git. Not actionable but worth noting.

Run `go mod tidy` after removing dead code to clean up.

---

## Security Issues

1. **Regex injection via author filter** (item 9): User-supplied author strings from query params are compiled as regex. A malicious `?author=(.*)` matches everything; pathological patterns could cause ReDoS.
2. **No input validation on HTTP endpoints**: Year parameter is safely handled (defaults on parse failure), but author is passed through unchecked.
3. **Server binds to all interfaces** (`cmd/gico-server/main.go:88`): `":8822"` listens on 0.0.0.0. For a local tool, consider `"127.0.0.1:8822"` or make it configurable.
4. **MD5 used for cache keys** (`commits/cache.go:39`): Not a security issue (not used cryptographically), but `hash/fnv` would be faster and more idiomatic for non-crypto hashing.

---

## Architecture Notes

**Good patterns:**
- Clean separation of concerns: `types/` for data, `commits/` for git ops, `graph/` for rendering, `ui/` for TUI
- Channel-based parallel repo parsing is well-structured
- Cache layer is a reasonable optimization for the TUI use case
- BubbleTea integration is solid with separate models per view

**Suggestions:**
- The `Repo` type wrapping `git.Repository` by value (`commits/common.go:16-19`) copies the entire struct. Use a pointer: `Repo *git.Repository`
- `RepoSet` as `[]string` is fine but consider a dedicated type with methods for validation
- The `Freq` type (`[]int`) having `String()` and `StringSelected()` methods that depend on `graph/term` creates a circular concern — the data type knows about its terminal rendering. Consider making these standalone functions in `graph/term` instead.

---

## README Accuracy

- **"svg-server"** is mentioned in the gico-server section but the binary is `gico-server` (renamed in commit b3e3047). Update the description.
- **Example GIF** placeholder exists but no image is linked.
- The backtick formatting for mgfetch has a stray backtick at the end of "terminal.`"
- No mention of the `--author` or any CLI flags (there aren't any — consider adding flag support).

---

## Summary of Recommended Actions

| Priority | Action |
|----------|--------|
| P0 | Fix panics in HTTP handlers → return proper HTTP errors |
| P0 | Fix nil pointer deref in `OpenRepo` |
| P0 | Fix leap year calculation (extract `yearLength` helper) |
| P1 | Replace goroutine-per-cache-entry with single reaper |
| P1 | Escape regex in author filter (or use `strings.Contains`) |
| P1 | Add tests for core functions |
| P1 | Remove dead code (svgToPng, unused vars, delegateKeyMap) |
| P2 | Fix Content-Type for SVG endpoints |
| P2 | Fix silent error continuation in `Frequency()` |
| P2 | Rename variables per Go conventions |
| P2 | Bump Go version, consider replacing archived gorilla deps |
| P3 | Bind server to localhost by default |
| P3 | Fix README inaccuracies |
| P3 | Remove `hashSlice` newline, switch MD5 → FNV |
