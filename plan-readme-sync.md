# Auto-sync README.md with `examples/00-docs/` outputs

## Context

The `justfile` currently warns: *"example0 is the example shown in the README. Be sure to change the README if the output ever changes."* That's a manual step and will drift. We want a command that regenerates example0 and then updates README.md in-place so the code blocks always match the committed example outputs. Target environment is Linux only; GNU awk/bash is assumed.

## What the README actually mirrors

Three blocks in `README.md` come from `examples/00-docs/`:

| # | README lines | Source file | Notes |
|---|---|---|---|
| 1 | 11–24 (```go) | `config.go` | README version prepends a `//go:generate` directive that is purely illustrative — real generation goes through `generate.go`. Do **not** auto-sync this block. |
| 2 | 28–153 (```go) | `config_gen.go` | Skips line 1 (`// Code generated …` header) |
| 3 | 173–179 (plain fence) | `.env` | No language fence; `.env` is gitignored but regenerated every run |

A fourth block (lines 157–169, illustrative `main()`) is hand-written and must NOT be auto-synced.

## Strategy: HTML-comment anchors + awk sync script

### Anchor scheme

Wrap each synced block with paired markers:

```markdown
<!-- sync:begin file=examples/00-docs/config_gen.go lang=go skip=1 -->
```go
…current contents…
```
<!-- sync:end -->
```

Attributes on the begin marker:
- `file=` — path relative to repo root (source of truth)
- `lang=` — fence language (empty string → no language tag, for the `.env` block)
- `skip=N` — drop the first N lines of the source file (handles the `// Code generated` header)

Everything between the two markers (inclusive of the fenced block) is replaced by the tool. Prose and any block outside markers is untouched. Block 1 stays unmarked and therefore unchanged.

### The tool: `tools/readme-sync.sh` (bash + awk)

Single awk script, driven by a thin bash wrapper. Approximate shape:

```bash
#!/usr/bin/env bash
set -euo pipefail
readme="${1:-README.md}"
check="${CHECK:-0}"
tmp=$(mktemp)
awk '
  /<!-- sync:begin / {
    match($0, /file=[^ ]+/);    file   = substr($0, RSTART+5, RLENGTH-5)
    match($0, /lang=[^ ]+/);    lang   = RLENGTH>0 ? substr($0, RSTART+5, RLENGTH-5) : ""
    match($0, /skip=[0-9]+/);   skip   = RLENGTH>0 ? substr($0, RSTART+5, RLENGTH-5)+0 : 0
    print
    print "```" lang
    n = 0
    while ((getline line < file) > 0) { n++; if (n > skip) print line }
    close(file)
    print "```"
    skipping = 1; next
  }
  /<!-- sync:end -->/ { skipping = 0; print; next }
  !skipping { print }
' "$readme" > "$tmp"

if [ "$check" = "1" ]; then
  diff -u "$readme" "$tmp" && rm "$tmp" || { rm "$tmp"; exit 1; }
else
  mv "$tmp" "$readme"
fi
```

- Default run: rewrites README in place.
- `CHECK=1 tools/readme-sync.sh` → diff mode, exits non-zero if out of sync. Good for CI / pre-commit.
- POSIX-ish awk; GNU awk on Linux handles this without issue.
- Idempotent: running twice produces no diff.

### Why awk over a Go tool or embedmd

- **awk**: ~30 lines, no build step, no binary, matches the repo's existing shell style (`test/run.sh`, justfile recipes are all shell-driven). Best fit for this scope.
- **Go tool**: overkill for 3 blocks of line-oriented replacement. Would need its own module, tools directory, and build wiring.
- **embedmd**: third-party dependency with regex-range syntax; more capable than we need.

## Files to create / modify

- **Create** `tools/readme-sync.sh` — the sync script described above. `chmod +x` it.
- **Modify** `README.md`:
  - Wrap block 2 with `<!-- sync:begin file=examples/00-docs/config_gen.go lang=go skip=1 -->` … `<!-- sync:end -->`.
  - Wrap block 3 with `<!-- sync:begin file=examples/00-docs/.env lang= -->` … `<!-- sync:end -->`.
  - Leave block 1 (struct with `//go:generate` directive) untouched — it's illustrative.
- **Modify** `justfile`:
  - New recipe `readme-sync: example0` → runs `./tools/readme-sync.sh`.
  - New recipe `readme-check` → runs `CHECK=1 ./tools/readme-sync.sh` for verification.
  - Add `readme-sync` to the `update-all` dependency chain so pre-commit workflow keeps README in sync.
- **Modify** `agents.md` — note that README sync happens via `just update-all`; do not hand-edit content between `<!-- sync:begin/end -->` markers.

## Verification

1. `just example0` regenerates `config_gen.go` and `.env`.
2. `just readme-sync` updates `README.md` in place.
3. `git diff README.md` shows no change on the first run (source files unchanged).
4. Temporarily edit `examples/00-docs/config.go` (e.g., rename a field), `just update-all`, confirm blocks 2 and 3 update to match the new generated output; block 1 stays as-is.
5. `just readme-check` exits 0 when in sync, non-zero with a diff when out of sync.
6. Running `just readme-sync` twice in a row produces no diff (idempotent).
