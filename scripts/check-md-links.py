#!/usr/bin/env python3
#
# check-md-links.py — report broken relative markdown links across
# every tracked .md file.
#
# Usage:
#   scripts/check-md-links.py
#
# Scans each tracked `.md` file for `[text](target)` links, resolves
# every relative target against the linking file's directory, and
# prints the ones that do not exist on disk. Skips http(s), mailto,
# tel, `#anchor`, and `/`-rooted (host-absolute) targets — the last
# class is tracked separately in docs/backlog.md `repo-rooted-doc-links`.
#
# Exits non-zero if any broken link is found. Not wired into CI; run
# on demand as a guard after doc moves/renames.

import os
import re
import subprocess
import sys

LINK_RE = re.compile(r"\[(?:[^\]]*)\]\(([^)]+)\)")
SKIP_PREFIXES = ("http://", "https://", "mailto:", "tel:", "#")


def main() -> int:
    root = subprocess.check_output(
        ["git", "rev-parse", "--show-toplevel"]
    ).decode().strip()
    files = subprocess.check_output(
        ["git", "ls-files", "*.md"], cwd=root
    ).decode().splitlines()

    broken = []
    for rel in files:
        path = os.path.join(root, rel)
        try:
            text = open(path, encoding="utf-8").read()
        except (OSError, UnicodeDecodeError):
            continue
        for m in LINK_RE.finditer(text):
            target = m.group(1).strip()
            if " " in target:  # strip optional (path "title")
                target = target.split(" ", 1)[0]
            if target.startswith(SKIP_PREFIXES):
                continue
            if target.startswith("/"):  # host-absolute; see backlog
                continue
            tgt = target.split("#", 1)[0]
            if not tgt:
                continue
            resolved = os.path.normpath(
                os.path.join(os.path.dirname(path), tgt)
            )
            if not os.path.exists(resolved):
                broken.append((rel, target))

    for rel, target in broken:
        print(f"{rel}  ->  {target}")
    if broken:
        print(f"\n{len(broken)} broken relative link(s).", file=sys.stderr)
        return 1
    print("No broken relative links.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
