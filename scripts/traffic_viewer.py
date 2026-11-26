#!/usr/bin/env python3
"""
Live local traffic viewer for Střela/Vlna (reads all events directly from Valkey).

Because `enigine` filters WebSocket events per assigned corrector on the backend,
a single corrector WS connection cannot see other correctors' tickets. Every game
event is persisted to Valkey, so this tool streams 100% of competition traffic
(all teams & all correctors) in real time without modifying or restarting enigine.

Usage:
  # Terminal stream (on the server, or locally via `ssh -L 6379:localhost:6379 ...`)
  python3 scripts/traffic_viewer.py

  # Local Web UI on http://127.0.0.1:8088 (plus terminal output)
  python3 scripts/traffic_viewer.py --web
"""

import argparse
import collections
import datetime
import http.server
import json
import queue
import re
import socket
import socketserver
import sqlite3
import sys
import threading
import time


class ValkeyClient:
    """Minimal zero-dependency RESP2 client over TCP (Python 3.5+ compatible)."""

    def __init__(self, host="localhost", port=6379, db=0, password=None, timeout=10):
        self._sock = socket.create_connection((host, port), timeout=timeout)
        self._rfile = self._sock.makefile("rb")
        if password:
            self.cmd("AUTH", password)
        if db != 0:
            self.cmd("SELECT", str(db))

    def close(self):
        try:
            self._rfile.close()
        finally:
            self._sock.close()

    def _read_line(self):
        line = self._rfile.readline()
        if not line:
            raise ConnectionError("Valkey connection closed unexpectedly")
        return line.rstrip(b"\r\n")

    def _read_resp(self):
        prefix = self._rfile.read(1)
        if not prefix:
            raise ConnectionError("Valkey connection closed unexpectedly")
        if prefix == b"+":
            return self._read_line().decode("utf-8", errors="replace")
        if prefix == b"-":
            raise RuntimeError("Valkey error: {}".format(self._read_line().decode("utf-8", errors="replace")))
        if prefix == b":":
            return int(self._read_line())
        if prefix == b"$":
            length = int(self._read_line())
            if length == -1:
                return None
            data = self._rfile.read(length)
            self._rfile.read(2)
            return data.decode("utf-8", errors="replace")
        if prefix == b"*":
            count = int(self._read_line())
            if count == -1:
                return None
            return [self._read_resp() for _ in range(count)]
        raise RuntimeError("Unexpected RESP prefix: {!r}".format(prefix))

    def cmd(self, *args):
        buf = ["*{}\r\n".format(len(args)).encode()]
        for arg in args:
            encoded = arg.encode("utf-8")
            buf.append("${}\r\n".format(len(encoded)).encode())
            buf.append(encoded)
            buf.append(b"\r\n")
        self._sock.sendall(b"".join(buf))
        return self._read_resp()

    def scan_iter(self, pattern, count=500):
        keys = []
        cursor = "0"
        while True:
            cursor, batch = self.cmd("SCAN", cursor, "MATCH", pattern, "COUNT", str(count))
            keys.extend(batch)
            if cursor == "0":
                break
        return keys

    def start_monitor(self):
        self._sock.settimeout(None)
        resp = self.cmd("MONITOR")
        if resp != "OK":
            raise RuntimeError("MONITOR failed: {!r}".format(resp))

    def read_monitor_line(self):
        prefix = self._rfile.read(1)
        if not prefix:
            raise ConnectionError("MONITOR stream closed")
        line = self._read_line()
        if prefix == b"-":
            raise RuntimeError("Valkey error: {}".format(line.decode("utf-8", errors="replace")))
        return line


MONITOR_PREFIX_RE = re.compile(rb"^(\d+\.\d+)\s+\[(\d+)\s+[^\]]+\]\s+(.*)$")


def parse_monitor_args(raw_bytes):
    """Parse Valkey sdscatrepr quoted arguments into UTF-8 strings."""
    args = []
    n = len(raw_bytes)
    i = 0
    while i < n:
        while i < n and raw_bytes[i] in (32, 9):  # space / tab
            i += 1
        if i >= n:
            break
        if raw_bytes[i] != 34:  # '"'
            # Unquoted token fallback
            start = i
            while i < n and raw_bytes[i] not in (32, 9):
                i += 1
            args.append(raw_bytes[start:i].decode("utf-8", errors="replace"))
            continue
        i += 1  # skip opening quote
        out = bytearray()
        while i < n:
            ch = raw_bytes[i]
            if ch == 34:  # closing '"'
                i += 1
                break
            if ch == 92 and i + 1 < n:  # '\\'
                nxt = raw_bytes[i + 1]
                if nxt == 120 and i + 3 < n:  # \xHH
                    try:
                        out.append(int(raw_bytes[i + 2 : i + 4], 16))
                        i += 4
                        continue
                    except ValueError:
                        pass
                esc_map = {
                    92: 92,   # \\
                    34: 34,   # \"
                    110: 10,  # \n
                    114: 13,  # \r
                    116: 9,   # \t
                    97: 7,    # \a
                    98: 8,    # \b
                }
                if nxt in esc_map:
                    out.append(esc_map[nxt])
                    i += 2
                    continue
            out.append(ch)
            i += 1
        args.append(out.decode("utf-8", errors="replace"))
    return args


def load_pb_usernames(pb_db_path):
    names = {}
    conn = sqlite3.connect("file:{}?mode=ro".format(pb_db_path), uri=True)
    try:
        for cid, username in conn.execute("SELECT id, username FROM correctors"):
            if username:
                names[cid] = username
    finally:
        conn.close()
    return names


class MetadataCache:
    """Caches team names, problem names/diffs, and corrector usernames."""

    def __init__(self, rdb, corr_names=None):
        self._rdb = rdb
        self._lock = threading.Lock()
        self.teams = {}
        self.probs = {}
        self.corrs = dict(corr_names or {})
        self.preload()

    def preload(self):
        with self._lock:
            for key in self._rdb.scan_iter("teamname:*"):
                tid = key[len("teamname:"):]
                val = self._rdb.cmd("GET", key)
                if val:
                    self.teams[tid] = val
            for key in self._rdb.scan_iter("prob:*"):
                pid = key[len("prob:"):]
                name, diff = self._rdb.cmd("HMGET", key, "name", "diff")
                if name or diff:
                    self.probs[pid] = (name or pid, diff or "?")

    def team_name(self, tid):
        with self._lock:
            if tid in self.teams:
                return self.teams[tid]
            val = self._rdb.cmd("GET", "teamname:" + tid)
            if val:
                self.teams[tid] = val
                return val
            return tid

    def prob_label(self, pid):
        with self._lock:
            if pid in self.probs:
                name, diff = self.probs[pid]
                return "[{}] {}".format(diff, name)
            name, diff = self._rdb.cmd("HMGET", "prob:" + pid, "name", "diff")
            if not name and "/" in pid:
                base = pid.split("/", 1)[0]
                if base in self.probs:
                    name, diff = self.probs[base]
                else:
                    name, diff = self._rdb.cmd("HMGET", "prob:" + base, "name", "diff")
            if name or diff:
                self.probs[pid] = (name or pid, diff or "?")
                return "[{}] {}".format(diff or "?", name or pid)
            return pid

    def corr_name(self, cid):
        if cid in self.corrs:
            return self.corrs[cid]
        return cid


def decode_event(ts_float, args, meta):
    """Turn a parsed Valkey MONITOR command into a structured game event (or None)."""
    if not args:
        return None
    cmd = args[0].upper()
    ts_str = datetime.datetime.fromtimestamp(ts_float).strftime("%H:%M:%S.%f")[:-3]

    if cmd == "RPUSH" and len(args) >= 3 and args[1].startswith("tline:"):
        key_parts = args[1].split(":", 2)
        if len(key_parts) != 3:
            return None
        tid, pid = key_parts[1], key_parts[2]
        atom = args[2]
        sparts = atom.split(":")
        if len(sparts) < 4:
            return None
        mside = sparts[0]
        mtype = sparts[1]
        msg = ":".join(sparts[2:-1])

        if mtype == "answer":
            kind = "ANSWER"
        elif mtype == "grade":
            kind = "GRADE"
        elif mtype in ("copy", "paste", "window-focus"):
            kind = "ANTICHEAT"
        else:
            kind = "CHAT"

        return {
            "ts": ts_str,
            "kind": kind,
            "subtype": mtype,
            "side": mside,
            "team_id": tid,
            "team": meta.team_name(tid),
            "prob_id": pid,
            "prob": meta.prob_label(pid),
            "detail": msg,
        }

    if cmd == "SET" and len(args) >= 3:
        key, val = args[1], args[2]

        if key.startswith("tstate:"):
            parts = key.split(":", 2)
            if len(parts) != 3:
                return None
            tid, pid = parts[1], parts[2]
            kind_map = {"bought": "BUY", "sold": "SELL", "solved": "SOLVE"}
            kind = kind_map.get(val, "STATE")
            return {
                "ts": ts_str,
                "kind": kind,
                "subtype": val,
                "side": "system",
                "team_id": tid,
                "team": meta.team_name(tid),
                "prob_id": pid,
                "prob": meta.prob_label(pid),
                "detail": val,
            }

        if key.startswith("money:"):
            tid = key[len("money:"):]
            return {
                "ts": ts_str,
                "kind": "MONEY",
                "subtype": "money",
                "side": "system",
                "team_id": tid,
                "team": meta.team_name(tid),
                "prob_id": "",
                "prob": "",
                "detail": "money -> {}".format(val),
            }

        if key.startswith("tcorr:"):
            parts = key.split(":", 2)
            if len(parts) != 3:
                return None
            tid, pid = parts[1], parts[2]
            return {
                "ts": ts_str,
                "kind": "ASSIGN",
                "subtype": "tcorr",
                "side": "system",
                "team_id": tid,
                "team": meta.team_name(tid),
                "prob_id": pid,
                "prob": meta.prob_label(pid),
                "detail": "assigned -> {}".format(meta.corr_name(val)),
            }

        if key == "state":
            return {
                "ts": ts_str,
                "kind": "STATE",
                "subtype": "state",
                "side": "admin",
                "team_id": "",
                "team": "*",
                "prob_id": "",
                "prob": "",
                "detail": "contest state -> {}".format(val),
            }

        if key.startswith("rank:"):
            tid = key[len("rank:"):]
            return {
                "ts": ts_str,
                "kind": "STATE",
                "subtype": "rank",
                "side": "system",
                "team_id": tid,
                "team": meta.team_name(tid),
                "prob_id": "",
                "prob": "",
                "detail": "rank -> #{}".format(val),
            }

    return None


ANSI_COLORS = {
    "BUY": "\033[36m",        # cyan
    "SELL": "\033[33m",       # yellow
    "SOLVE": "\033[32;1m",    # bold green
    "ANSWER": "\033[35m",     # magenta
    "GRADE": "\033[34;1m",    # bold blue
    "CHAT": "\033[37m",       # white
    "MONEY": "\033[32m",      # green
    "ASSIGN": "\033[90m",     # gray
    "ANTICHEAT": "\033[31m",  # red
    "STATE": "\033[33;1m",    # bold yellow
}
ANSI_RESET = "\033[0m"


def format_terminal(ev, use_color=True):
    tag = "{:<9}".format(ev["kind"])
    if use_color:
        color = ANSI_COLORS.get(ev["kind"], "")
        tag = "{}{}{}".format(color, tag, ANSI_RESET)

    who = ev["team"]
    if ev["prob"]:
        who = "{} | {}".format(who, ev["prob"])

    if ev["kind"] in ("CHAT", "ANSWER", "GRADE", "ANTICHEAT"):
        arrow = "corr->team" if ev["side"] == "admin" else "team->corr"
        extra = "({}) {}".format(ev["subtype"], ev["detail"]) if ev["kind"] == "ANTICHEAT" else ev["detail"]
        return "{} {} [{}] {} :: {}".format(ev["ts"], tag, arrow, who, extra)

    return "{} {} {} :: {}".format(ev["ts"], tag, who, ev["detail"])


class EventHub:
    """Thread-safe ring buffer + fan-out queues for SSE web clients."""

    def __init__(self, maxlen=500):
        self._lock = threading.Lock()
        self._history = collections.deque(maxlen=maxlen)
        self._clients = []

    def publish(self, ev):
        payload = json.dumps(ev, ensure_ascii=False)
        with self._lock:
            self._history.append(payload)
            dead = []
            for q in self._clients:
                try:
                    q.put_nowait(payload)
                except queue.Full:
                    dead.append(q)
            for q in dead:
                self._clients.remove(q)

    def subscribe(self):
        q = queue.Queue(maxsize=1000)
        with self._lock:
            snapshot = list(self._history)
            self._clients.append(q)
        return snapshot, q

    def unsubscribe(self, q):
        with self._lock:
            if q in self._clients:
                self._clients.remove(q)


HTML_PAGE = """<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8" />
<title>Střela/Vlna — Live Traffic Viewer</title>
<style>
  :root {
    --bg: #0f1117; --panel: #171a22; --border: #282d3a; --text: #e6e9ef; --muted: #8b93a6;
  }
  * { box-sizing: border-box; }
  body { margin: 0; font: 13px/1.45 ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; background: var(--bg); color: var(--text); display: flex; flex-direction: column; height: 100vh; }
  header { padding: 10px 14px; background: var(--panel); border-bottom: 1px solid var(--border); display: flex; flex-wrap: wrap; gap: 10px; align-items: center; }
  header h1 { margin: 0; font-size: 15px; font-weight: 700; margin-right: 6px; }
  input[type=text] { background: #0b0d12; color: var(--text); border: 1px solid var(--border); border-radius: 4px; padding: 5px 8px; min-width: 220px; font: inherit; }
  .pills { display: flex; flex-wrap: wrap; gap: 4px; }
  .pill { border: 1px solid var(--border); background: #0b0d12; color: var(--muted); padding: 3px 8px; border-radius: 4px; cursor: pointer; font: inherit; font-size: 11px; }
  .pill.on { color: #fff; border-color: #5a6782; background: #232a3b; }
  button.ctrl { background: #232a3b; color: var(--text); border: 1px solid var(--border); border-radius: 4px; padding: 4px 10px; cursor: pointer; font: inherit; }
  #status { margin-left: auto; font-size: 12px; color: var(--muted); }
  #log { flex: 1; overflow-y: auto; padding: 8px 12px; }
  .row { display: grid; grid-template-columns: 92px 86px 260px 1fr; gap: 10px; padding: 3px 6px; border-bottom: 1px solid #1a1e28; align-items: baseline; }
  .row:hover { background: #161a24; }
  .ts { color: var(--muted); }
  .kind { font-weight: 700; }
  .k-BUY { color: #4dd0e1; }
  .k-SELL { color: #ffd54f; }
  .k-SOLVE { color: #66bb6a; }
  .k-ANSWER { color: #ce93d8; }
  .k-GRADE { color: #64b5f6; }
  .k-CHAT { color: #e6e9ef; }
  .k-MONEY { color: #81c784; }
  .k-ASSIGN { color: #788296; }
  .k-ANTICHEAT { color: #ef5350; }
  .k-STATE { color: #ffb74d; }
  .who { color: #c5cde0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .detail { white-space: pre-wrap; word-break: break-word; }
  .side { font-size: 11px; padding: 1px 4px; border-radius: 3px; margin-right: 6px; }
  .s-player { background: #2d1f3a; color: #e1bee7; }
  .s-admin { background: #1c2e45; color: #90caf9; }
</style>
</head>
<body>
<header>
  <h1>Live Traffic</h1>
  <input id="q" type="text" placeholder="Filter team, prob, text..." />
  <div class="pills" id="pills"></div>
  <button class="ctrl" id="pauseBtn">Pause</button>
  <button class="ctrl" id="clearBtn">Clear</button>
  <span id="status">Connecting...</span>
</header>
<div id="log"></div>
<script>
  const KINDS = ["BUY","SELL","SOLVE","ANSWER","GRADE","CHAT","MONEY","ASSIGN","ANTICHEAT","STATE"];
  const enabled = new Set(KINDS);
  const events = [];
  let paused = false;
  let autoScroll = true;

  const logEl = document.getElementById("log");
  const qEl = document.getElementById("q");
  const statusEl = document.getElementById("status");
  const pillsEl = document.getElementById("pills");
  const pauseBtn = document.getElementById("pauseBtn");
  const clearBtn = document.getElementById("clearBtn");

  KINDS.forEach(k => {
    const b = document.createElement("button");
    b.className = "pill on k-" + k;
    b.textContent = k;
    b.onclick = () => {
      if (enabled.has(k)) { enabled.delete(k); b.classList.remove("on"); }
      else { enabled.add(k); b.classList.add("on"); }
      renderAll();
    };
    pillsEl.appendChild(b);
  });

  logEl.addEventListener("scroll", () => {
    autoScroll = (logEl.scrollTop + logEl.clientHeight >= logEl.scrollHeight - 40);
  });
  qEl.addEventListener("input", renderAll);
  pauseBtn.onclick = () => { paused = !paused; pauseBtn.textContent = paused ? "Resume" : "Pause"; if (!paused) renderAll(); };
  clearBtn.onclick = () => { events.length = 0; logEl.innerHTML = ""; };

  function matches(ev) {
    if (!enabled.has(ev.kind)) return false;
    const q = qEl.value.trim().toLowerCase();
    if (!q) return true;
    return (ev.team + " " + ev.team_id + " " + ev.prob + " " + ev.prob_id + " " + ev.detail + " " + ev.subtype).toLowerCase().includes(q);
  }

  function makeRow(ev) {
    const div = document.createElement("div");
    div.className = "row";
    const ts = document.createElement("span"); ts.className = "ts"; ts.textContent = ev.ts;
    const kind = document.createElement("span"); kind.className = "kind k-" + ev.kind; kind.textContent = ev.kind;
    const who = document.createElement("span"); who.className = "who";
    who.textContent = ev.prob ? (ev.team + " | " + ev.prob) : ev.team;
    who.title = who.textContent;
    const det = document.createElement("span"); det.className = "detail";
    if (ev.side === "player" || ev.side === "admin") {
      const s = document.createElement("span");
      s.className = "side s-" + ev.side;
      s.textContent = ev.side === "admin" ? "corr" : "team";
      det.appendChild(s);
    }
    const txt = ev.kind === "ANTICHEAT" ? ("[" + ev.subtype + "] " + ev.detail) : ev.detail;
    det.appendChild(document.createTextNode(txt));
    div.append(ts, kind, who, det);
    return div;
  }

  function renderAll() {
    logEl.innerHTML = "";
    const frag = document.createDocumentFragment();
    for (const ev of events) {
      if (matches(ev)) frag.appendChild(makeRow(ev));
    }
    logEl.appendChild(frag);
    if (autoScroll) logEl.scrollTop = logEl.scrollHeight;
  }

  const es = new EventSource("/events");
  es.onopen = () => { statusEl.textContent = "Connected (" + events.length + " events)"; };
  es.onerror = () => { statusEl.textContent = "Reconnecting..."; };
  es.onmessage = (e) => {
    const ev = JSON.parse(e.data);
    events.push(ev);
    if (events.length > 2000) events.shift();
    statusEl.textContent = "Connected (" + events.length + " events)";
    if (!paused && matches(ev)) {
      logEl.appendChild(makeRow(ev));
      if (logEl.children.length > 2000) logEl.removeChild(logEl.firstChild);
      if (autoScroll) logEl.scrollTop = logEl.scrollHeight;
    }
  };
</script>
</body>
</html>
"""


class ThreadingHTTPServer(socketserver.ThreadingMixIn, http.server.HTTPServer):
    daemon_threads = True
    allow_reuse_address = True


def start_web_server(hub, bind_host, bind_port):
    class Handler(http.server.BaseHTTPRequestHandler):
        def log_message(self, fmt, *args):
            pass  # keep terminal output clean

        def do_GET(self):
            if self.path == "/" or self.path.startswith("/?"):
                body = HTML_PAGE.encode("utf-8")
                self.send_response(200)
                self.send_header("Content-Type", "text/html; charset=utf-8")
                self.send_header("Content-Length", str(len(body)))
                self.end_headers()
                self.wfile.write(body)
                return

            if self.path == "/events":
                self.send_response(200)
                self.send_header("Content-Type", "text/event-stream; charset=utf-8")
                self.send_header("Cache-Control", "no-cache")
                self.send_header("Connection", "keep-alive")
                self.end_headers()
                snapshot, q = hub.subscribe()
                try:
                    for payload in snapshot:
                        self.wfile.write(("data: " + payload + "\n\n").encode("utf-8"))
                    self.wfile.flush()
                    while True:
                        try:
                            payload = q.get(timeout=15)
                            self.wfile.write(("data: " + payload + "\n\n").encode("utf-8"))
                            self.wfile.flush()
                        except queue.Empty:
                            self.wfile.write(b": keepalive\n\n")
                            self.wfile.flush()
                except (BrokenPipeError, ConnectionResetError, OSError):
                    pass
                finally:
                    hub.unsubscribe(q)
                return

            self.send_error(404)

    server = ThreadingHTTPServer((bind_host, bind_port), Handler)
    t = threading.Thread(target=server.serve_forever, daemon=True)
    t.start()
    return server


def main():
    parser = argparse.ArgumentParser(
        description="Stream all Střela/Vlna competition traffic live from Valkey."
    )
    parser.add_argument("--host", default="localhost", help="Valkey host (default: localhost)")
    parser.add_argument("--port", type=int, default=6379, help="Valkey port (default: 6379)")
    parser.add_argument("--db", type=int, default=0, help="Valkey DB index to watch (default: 0)")
    parser.add_argument("--password", default=None, help="Valkey password (optional)")
    parser.add_argument("--pb-db", default=None, help="Optional PocketBase data.db path for corrector usernames")
    parser.add_argument("--team", default=None, help="Filter terminal output by team name or ID (substring)")
    parser.add_argument("--prob", default=None, help="Filter terminal output by prob name or ID (substring)")
    parser.add_argument("--hide-anticheat", action="store_true", help="Hide copy/paste/window-focus in terminal")
    parser.add_argument("--hide-assign", action="store_true", help="Hide corrector assignment events in terminal")
    parser.add_argument("--no-color", action="store_true", help="Disable ANSI colors in terminal")
    parser.add_argument("--quiet", action="store_true", help="Suppress terminal event output (useful with --web)")
    parser.add_argument("--web", action="store_true", help="Start a local web dashboard")
    parser.add_argument("--web-host", default="127.0.0.1", help="Web UI bind host (default: 127.0.0.1)")
    parser.add_argument("--web-port", type=int, default=8088, help="Web UI port (default: 8088)")
    args = parser.parse_args()

    corr_names = load_pb_usernames(args.pb_db) if args.pb_db else {}
    lookup_rdb = ValkeyClient(args.host, args.port, args.db, args.password)
    meta = MetadataCache(lookup_rdb, corr_names)

    hub = EventHub() if args.web else None
    if args.web:
        start_web_server(hub, args.web_host, args.web_port)
        sys.stderr.write("Web UI listening on http://{}:{}\n".format(args.web_host, args.web_port))
        sys.stderr.flush()

    mon_rdb = ValkeyClient(args.host, args.port, args.db, args.password)
    mon_rdb.start_monitor()
    sys.stderr.write(
        "Streaming all Valkey traffic on {}:{} (db={}, {} teams / {} probs cached)... Press Ctrl+C to stop.\n".format(
            args.host, args.port, args.db, len(meta.teams), len(meta.probs)
        )
    )
    sys.stderr.flush()

    use_color = (not args.no_color) and sys.stdout.isatty()
    team_filter = args.team.lower() if args.team else None
    prob_filter = args.prob.lower() if args.prob else None

    try:
        while True:
            raw = mon_rdb.read_monitor_line()
            m = MONITOR_PREFIX_RE.match(raw)
            if not m:
                continue
            ts_float = float(m.group(1))
            db_idx = int(m.group(2))
            if db_idx != args.db:
                continue
            cmd_args = parse_monitor_args(m.group(3))
            ev = decode_event(ts_float, cmd_args, meta)
            if ev is None:
                continue

            if hub is not None:
                hub.publish(ev)

            if args.quiet:
                continue
            if args.hide_anticheat and ev["kind"] == "ANTICHEAT":
                continue
            if args.hide_assign and ev["kind"] == "ASSIGN":
                continue
            if team_filter and team_filter not in (ev["team"] + " " + ev["team_id"]).lower():
                continue
            if prob_filter and prob_filter not in (ev["prob"] + " " + ev["prob_id"]).lower():
                continue

            print(format_terminal(ev, use_color=use_color), flush=True)
    except KeyboardInterrupt:
        pass
    finally:
        mon_rdb.close()
        lookup_rdb.close()

    return 0


if __name__ == "__main__":
    sys.exit(main())
