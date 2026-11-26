#!/usr/bin/env python3
"""
Inspect problem corrector queue distribution in Valkey/Redis.

Prints each corrector along with:
  - how many problems they own (1st in `prob:<id>` queue)
  - how many problems they are 2nd in the queue for
  - per-difficulty (A/B/C) breakdown and uniformity summary stats
"""

from __future__ import annotations

import argparse
import math
import socket
import sqlite3
import sys
from collections import defaultdict
from typing import Any


class ValkeyClient:
    """Minimal zero-dependency RESP2 client over TCP."""

    def __init__(self, host: str = "localhost", port: int = 6379, db: int = 0, password: str | None = None):
        self._sock = socket.create_connection((host, port), timeout=10)
        self._rfile = self._sock.makefile("rb")
        if password:
            self.cmd("AUTH", password)
        if db != 0:
            self.cmd("SELECT", str(db))

    def close(self) -> None:
        try:
            self._rfile.close()
        finally:
            self._sock.close()

    def __enter__(self) -> "ValkeyClient":
        return self

    def __exit__(self, *_: Any) -> None:
        self.close()

    def _read_line(self) -> bytes:
        line = self._rfile.readline()
        if not line:
            raise ConnectionError("Valkey connection closed unexpectedly")
        return line.rstrip(b"\r\n")

    def _read_resp(self) -> Any:
        prefix = self._rfile.read(1)
        if not prefix:
            raise ConnectionError("Valkey connection closed unexpectedly")
        if prefix == b"+":
            return self._read_line().decode("utf-8", errors="replace")
        if prefix == b"-":
            raise RuntimeError(f"Valkey error: {self._read_line().decode('utf-8', errors='replace')}")
        if prefix == b":":
            return int(self._read_line())
        if prefix == b"$":
            length = int(self._read_line())
            if length == -1:
                return None
            data = self._rfile.read(length)
            self._rfile.read(2)  # trailing CRLF
            return data.decode("utf-8", errors="replace")
        if prefix == b"*":
            count = int(self._read_line())
            if count == -1:
                return None
            return [self._read_resp() for _ in range(count)]
        raise RuntimeError(f"Unexpected RESP prefix: {prefix!r}")

    def cmd(self, *args: str) -> Any:
        buf = [f"*{len(args)}\r\n".encode()]
        for arg in args:
            encoded = arg.encode("utf-8")
            buf.append(f"${len(encoded)}\r\n".encode())
            buf.append(encoded)
            buf.append(b"\r\n")
        self._sock.sendall(b"".join(buf))
        return self._read_resp()

    def scan_iter(self, pattern: str, count: int = 500) -> list[str]:
        keys: list[str] = []
        cursor = "0"
        while True:
            cursor, batch = self.cmd("SCAN", cursor, "MATCH", pattern, "COUNT", str(count))
            keys.extend(batch)
            if cursor == "0":
                break
        return keys


def load_usernames(pb_db_path: str) -> dict[str, str]:
    names: dict[str, str] = {}
    conn = sqlite3.connect(f"file:{pb_db_path}?mode=ro", uri=True)
    try:
        for cid, username in conn.execute("SELECT id, username FROM correctors"):
            if username:
                names[cid] = username
    finally:
        conn.close()
    return names


def stats_summary(values: list[int]) -> str:
    if not values:
        return "n/a"
    n = len(values)
    mean = sum(values) / n
    variance = sum((v - mean) ** 2 for v in values) / n
    stdev = math.sqrt(variance)
    return f"min={min(values)}  max={max(values)}  mean={mean:.2f}  stdev={stdev:.2f}"


def fmt_diffs(counts: dict[str, int]) -> str:
    return f"({counts['A']}/{counts['B']}/{counts['C']})"


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Show how many problems each corrector owns (1st) and is 2nd in queue for in Valkey."
    )
    parser.add_argument("--host", default="localhost", help="Valkey host (default: localhost)")
    parser.add_argument("--port", type=int, default=6379, help="Valkey port (default: 6379)")
    parser.add_argument("--db", type=int, default=0, help="Valkey DB index (default: 0)")
    parser.add_argument("--password", default=None, help="Valkey password (optional)")
    parser.add_argument(
        "--include-instances",
        action="store_true",
        help="Include generated auto-problem instances (keys containing '/')",
    )
    parser.add_argument(
        "--sort",
        choices=("owned", "second", "total", "id"),
        default="owned",
        help="Sort output by column (default: owned)",
    )
    parser.add_argument(
        "--pb-db",
        default=None,
        help="Optional path to PocketBase data.db to resolve corrector usernames",
    )
    args = parser.parse_args()

    usernames: dict[str, str] = {}
    if args.pb_db:
        usernames = load_usernames(args.pb_db)

    with ValkeyClient(args.host, args.port, args.db, args.password) as rdb:
        correctors: set[str] = set()

        for key in rdb.scan_iter("corradmin:*"):
            cid = key[len("corradmin:") :]
            if cid:
                correctors.add(cid)

        for key in rdb.scan_iter("corrtoken:*"):
            cid = rdb.cmd("GET", key)
            if cid:
                correctors.add(cid)

        owned: dict[str, int] = defaultdict(int)
        second: dict[str, int] = defaultdict(int)
        owned_diff: dict[str, dict[str, int]] = defaultdict(lambda: defaultdict(int))
        second_diff: dict[str, dict[str, int]] = defaultdict(lambda: defaultdict(int))

        prob_keys = rdb.scan_iter("prob:*")
        total_probs = 0
        skipped_instances = 0
        empty_queue_probs = 0
        short_queue_probs = 0

        for key in prob_keys:
            prob_id = key[len("prob:") :]
            if not args.include_instances and "/" in prob_id:
                skipped_instances += 1
                continue

            total_probs += 1
            raw_queue, diff = rdb.cmd("HMGET", key, "queue", "diff")
            diff = diff or "?"
            queue = [c for c in (raw_queue or "").split(":") if c]

            for c in queue:
                correctors.add(c)

            if len(queue) == 0:
                empty_queue_probs += 1
                continue

            first_corr = queue[0]
            owned[first_corr] += 1
            owned_diff[first_corr][diff] += 1

            if len(queue) >= 2:
                second_corr = queue[1]
                second[second_corr] += 1
                second_diff[second_corr][diff] += 1
            else:
                short_queue_probs += 1

    rows = []
    for cid in correctors:
        label = f"{usernames[cid]} ({cid})" if cid in usernames else cid
        o = owned[cid]
        s = second[cid]
        rows.append((cid, label, o, s, o + s))

    if args.sort == "owned":
        rows.sort(key=lambda r: (-r[2], -r[3], r[1]))
    elif args.sort == "second":
        rows.sort(key=lambda r: (-r[3], -r[2], r[1]))
    elif args.sort == "total":
        rows.sort(key=lambda r: (-r[4], -r[2], r[1]))
    else:
        rows.sort(key=lambda r: r[1])

    name_width = max([len("Corrector")] + [len(r[1]) for r in rows])

    header = (
        f"{'Corrector':<{name_width}}  "
        f"{'Owned':>5} {'(A/B/C)':<11}  "
        f"{'2nd':>5} {'(A/B/C)':<11}  "
        f"{'Top-2':>5}"
    )
    print(header)
    print("-" * len(header))
    for cid, label, o, s, tot in rows:
        print(
            f"{label:<{name_width}}  "
            f"{o:>5} {fmt_diffs(owned_diff[cid]):<11}  "
            f"{s:>5} {fmt_diffs(second_diff[cid]):<11}  "
            f"{tot:>5}"
        )

    print("-" * len(header))
    print(f"Correctors: {len(rows)} | Problems: {total_probs}", end="")
    if skipped_instances:
        print(f" (skipped {skipped_instances} auto-generated instances)", end="")
    print()
    if empty_queue_probs:
        print(f"Warning: {empty_queue_probs} problem(s) have an empty queue")
    if short_queue_probs:
        print(f"Warning: {short_queue_probs} problem(s) have fewer than 2 correctors in queue")

    if rows:
        print(f"Owned  distribution: {stats_summary([r[2] for r in rows])}")
        print(f"2nd    distribution: {stats_summary([r[3] for r in rows])}")
        print(f"Top-2  distribution: {stats_summary([r[4] for r in rows])}")

    return 0


if __name__ == "__main__":
    sys.exit(main())
