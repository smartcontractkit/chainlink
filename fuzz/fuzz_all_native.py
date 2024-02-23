#!/usr/bin/env python3

import argparse
import itertools
import os
import re
import subprocess
import sys

LIBROOT = "../"

def main():
    parser = argparse.ArgumentParser(
        formatter_class=argparse.RawDescriptionHelpFormatter,
        description="\n".join([
            "Fuzz helper to run all native go fuzzers in chainlink",
            "",
        ]),
    )
    parser.add_argument("--ci", required=False, help="In CI mode we run each parser only briefly once", action="store_true")
    parser.add_argument("--seconds", required=False, help="Run for this many seconds of total fuzz time before exiting")
    args = parser.parse_args()

    # use float for remaining_seconds so we can represent infinity
    if args.seconds:
        remaining_seconds = float(args.seconds)
    else:
        remaining_seconds = float("inf")

    fuzzers = discover_fuzzers()
    print(f"🐝 Discovered fuzzers:", file=sys.stderr)
    for fuzzfn, path in fuzzers.items():
        print(f"{fuzzfn} in {path}", file=sys.stderr)

    if args.ci:
        # only run each fuzzer once for 60 seconds in CI
        durations_seconds = [60]
    else:
        # run forever or until --seconds, with increasingly longer durations per fuzz run
        durations_seconds = itertools.chain([5, 10, 30, 90, 270], itertools.repeat(600))

    for duration_seconds in durations_seconds:
        print(f"🐝 Running each fuzzer for {duration_seconds}s before switching to next fuzzer", file=sys.stderr)
        for fuzzfn, path in fuzzers.items():
            if remaining_seconds <= 0:
                print(f"🐝 Time budget of {args.seconds}s is exhausted. Exiting.", file=sys.stderr)
                return

            next_duration_seconds = min(remaining_seconds, duration_seconds)
            remaining_seconds -= next_duration_seconds

            print(f"🐝 Running {fuzzfn} in {path} for {next_duration_seconds}s before switching to next fuzzer", file=sys.stderr)
            run_fuzzer(fuzzfn, path, next_duration_seconds)
            print(f"🐝 Completed running {fuzzfn} in {path} for {next_duration_seconds}s. Total remaining time is {remaining_seconds}s", file=sys.stderr)

def discover_fuzzers():
    fuzzers = {}
    for root, dirs, files in os.walk(LIBROOT):
        for file in files:
            if not file.endswith("test.go"): continue
            with open(os.path.join(root, file), "r") as f:
                text = f.read()
            # ignore multiline comments
            text = re.sub(r"(?s)/[*].*?[*]/", "", text)
            # ignore single line comments *except* build tags
            text = re.sub(r"//.*", "", text)
            # Find every function with a name like FuzzXXX
            for fuzzfn in re.findall(r"func\s+(Fuzz\w+)", text):
                if fuzzfn in fuzzers:
                    raise Exception(f"Duplicate fuzz function: {fuzzfn}")
                fuzzers[fuzzfn] = os.path.relpath(root, LIBROOT)
    return fuzzers

def run_fuzzer(fuzzfn, dir, duration_seconds):
    subprocess.check_call(["go", "test", "-run=^$", f"-fuzz=^{fuzzfn}$", f"-fuzztime={duration_seconds}s", f"./{dir}"], cwd=LIBROOT)

if __name__ == "__main__":
    main()