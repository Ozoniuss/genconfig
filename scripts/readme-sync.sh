#!/usr/bin/env bash
set -euo pipefail
tmp=$(mktemp)
awk '
  /<!-- sync:begin / {
    match($0, /file=[^ ]+/); file = substr($0, RSTART+5, RLENGTH-5)
    match($0, /lang=[^ ]+/); lang = RLENGTH>0 ? substr($0, RSTART+5, RLENGTH-5) : ""
    match($0, /skip=[0-9]+/); skip = RLENGTH>0 ? substr($0, RSTART+5, RLENGTH-5)+0 : 0
    print
    print "```" lang
    n = 0
    while ((getline line < file) > 0) if (++n > skip) print line
    close(file)
    print "```"
    skipping = 1; next
  }
  /<!-- sync:end -->/ { skipping = 0 }
  !skipping { print }
' README.md > "$tmp"
mv "$tmp" README.md
