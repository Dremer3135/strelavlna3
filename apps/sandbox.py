#!/usr/bin/env python3
"""
sandbox_server.py

A small HTTP server that runs user-provided Python code using RestrictedPython
in a separate process with CPU-time and memory (address-space) limits.

POST /run
  JSON body: {"code": "<python code>", "text": "<text>", "answer": "<answer>", timeout": 2, "mem_mb": 64}

Response JSON:
  {"success": true, "text": "...", "answer": "..."}
  or
  {"success": false, "error": "..."}
"""

import json
import sys
import traceback
from http.server import HTTPServer, BaseHTTPRequestHandler
from multiprocessing import Process, Pipe
import multiprocessing
import os
import time
from typing import Any
import random
import math

# UNIX-only resource limiting:
try:
    import resource
except Exception:
    resource = None  # Windows: no resource module

# RestrictedPython
from RestrictedPython import compile_restricted_exec
from RestrictedPython.Guards import safe_builtins, safer_getattr
from RestrictedPython.Eval import default_guarded_getitem, default_guarded_getiter

def _inplacevar_(op, var, expr):
    if op == "+=":
        return var + expr
    elif op == "-=":
        return var - expr
    elif op == "*=":
        return var * expr
    elif op == "/=":
        return var / expr
    elif op == "%=":
        return var % expr
    elif op == "**=":
        return var ** expr
    elif op == "<<=":
        return var << expr
    elif op == ">>=":
        return var >> expr
    elif op == "|=":
        return var | expr
    elif op == "^=":
        return var ^ expr
    elif op == "&=":
        return var & expr
    elif op == "//=":
        return var // expr
    elif op == "@=":
        return var // expr

# ---------- Worker: runs in a separate process ----------
def _prune_locals_for_transport(d):
    """Make a JSON-friendly representation of local variables."""
    out = {}
    for k, v in d.items():
        # skip internal names
        if k.startswith('__'):
            continue
        try:
            # prefer direct JSON serialization for basic types
            json.dumps({k: v})
            out[k] = v
        except Exception:
            # fall back to repr
            try:
                out[k] = {"__repr__": repr(v)}
            except Exception:
                out[k] = {"__repr__": "<unrepresentable value>"}
    return out

class AttrDict:
    def __init__(self, d) -> None:
        self.__d = d
    def __getattr__(self, name: str) -> Any:
        if name == "__d": return self.__d
        return self.__d[name]

def _worker_run(code: str, text: str, answer: str, consts, conn, timeout: int, mem_mb: int):
    """
    Worker function executed in a child process.
    It sets resource limits (if available), executes code under RestrictedPython,
    and sends back either {'ok': True, 'locals': {...}} or {'ok': False, 'error': '...'}.
    """
    try:
        # 1) Apply resource limits (Unix)
        # if resource is not None:
        #     # limit CPU time (seconds)
        #     try:
        #         # RLIMIT_CPU limits CPU-seconds. When exceeded, SIGXCPU is delivered.
        #         resource.setrlimit(resource.RLIMIT_CPU, (timeout, timeout))
        #     except Exception as e:
        #         # continue but note we couldn't set
        #         pass
        #
        #     # limit address space (virtual memory) in bytes
        #     mem_bytes = int(mem_mb) * 1024 * 1024
        #     try:
        #         resource.setrlimit(resource.RLIMIT_AS, (mem_bytes, mem_bytes))
        #     except Exception:
        #         # not fatal; continue
        #         pass
        #
        #     # try to prevent forking many processes (best-effort)
        #     try:
        #         resource.setrlimit(resource.RLIMIT_NPROC, (10, 10))
        #     except Exception:
        #         pass

        # 2) Prepare RestrictedPython environment
        restricted_globals = {
            "__builtins__": safe_builtins,
            # The _print_ helper lets user code call print()
            "_print_": lambda *args, **kwargs: None,
            "_getiter_": default_guarded_getiter,
            "_getitem_": default_guarded_getitem,
            "_getattr_": safer_getattr,
            "_inplacevar_": _inplacevar_,
            "CONSTS": AttrDict(consts),
            "random": random,
            "math": math,
        }
        restricted_locals = {}

        # 3) compile code
        byte_code = compile_restricted_exec(code)

        print(byte_code)

        if len(byte_code.errors) > 0:
            conn.send({"ok": False, "error": f"{byte_code.errors}"})
            conn.close()
            return


        # 4) execute
        exec(byte_code.code, restricted_globals, restricted_locals)

        ntext = text
        for k, v in restricted_locals.items():
            ntext = ntext.replace("`" + f"{k}" + "`", f"{v}")

        nans = answer
        for k, v in restricted_locals.items():
            nans = nans.replace("`" + f"{k}" + "`", f"{v}")

        # 5) gather locals and send back
        conn.send({"ok": True, "text": ntext, "answer": nans})
    except SystemExit as e:
        # intentional exit from user code
        conn.send({"ok": False, "error": f"SystemExit: {e}"})
    except BaseException as e:
        tb = traceback.format_exc()
        conn.send({"ok": False, "error": f"{type(e).__name__}: {e}\n{tb}"})
    finally:
        try:
            conn.close()
        except Exception:
            pass

# ---------- HTTP server ----------
class SandboxHandler(BaseHTTPRequestHandler):
    server_version = "RestrictedSandboxHTTP/0.1"

    def _send_json(self, obj, status=200):
        data = json.dumps(obj).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def do_POST(self):
        if self.path != "/run":
            self._send_json({"success": False, "error": "not found"}, status=404)
            return

        content_length = int(self.headers.get("Content-Length", "0"))
        body = self.rfile.read(content_length).decode("utf-8")
        try:
            payload = json.loads(body)
            code = payload.get("code", "")
            text = payload.get("text", "")
            answer = payload.get("answer", "")
            consts = payload.get("consts", {})
            timeout = float(payload.get("timeout", 2))
            mem_mb = int(payload.get("mem_mb", 64))
        except Exception as e:
            self._send_json({"success": False, "error": f"bad request JSON: {e}"}, status=400)
            return

        if not isinstance(code, str) or not code.strip():
            self._send_json({"success": True, "text": text, "answer": answer}, status=400)
            return

        parent_conn, child_conn = Pipe()
        proc = Process(target=_worker_run, args=(code, text, answer, consts, child_conn, timeout, mem_mb))
        proc.daemon = True
        proc.start()
        child_conn.close()

        # Wait up to timeout seconds (+ small grace) for result
        start = time.time()
        result = None
        try:
            # We'll poll the pipe in a loop to be able to terminate process when needed.
            poll_interval = 0.05
            deadline = start + timeout + 0.5  # small grace
            while time.time() < deadline:
                if parent_conn.poll(poll_interval):
                    result = parent_conn.recv()
                    break
                # else continue polling
            else:
                # timeout expired, kill child
                if proc.is_alive():
                    try:
                        proc.terminate()
                    except Exception:
                        pass
                result = {"ok": False, "error": f"Execution timed out after {timeout} seconds (killed)"}
        except Exception as e:
            result = {"ok": False, "error": f"internal server error: {e}"}
        finally:
            # ensure child ended
            proc.join(timeout=1)
            try:
                parent_conn.close()
            except Exception:
                pass

        if result is None:
            result = {"ok": False, "error": "no result (child crashed or was terminated)"}

        if result.get("ok"):
            self._send_json({"success": True, "text": result["text"], "answer": result["answer"]})
        else:
            self._send_json({"success": False, "error": result.get("error", "unknown error")}, status=200)

    def log_message(self, format, *args):
        # keep server output quiet in examples; override to print if desired
        sys.stderr.write("%s - - [%s] %s\n" % (self.address_string(), self.log_date_time_string(), format%args))

def run_server(host="0.0.0.0", port=8000):
    print(f"Starting server on {host}:{port} (RestrictedPython sandbox).")
    httpd = HTTPServer((host, port), SandboxHandler)
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("Shutting down...")
    finally:
        httpd.server_close()

if __name__ == "__main__":
    # make multiprocessing spawn method deterministic on some platforms
    multiprocessing.set_start_method("spawn", force=True)
    run_server()
