#!/usr/bin/env python3

import argparse
import itertools
import os
import re
import subprocess
import sys
import time

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
    parser.add_argument("--go_module_root", required=True, help="Path to the root of the go module to fuzz")
    args = parser.parse_args()

    # use float for remaining_seconds so we can represent infinity
    if args.seconds:
        total_time = float(args.seconds)
    else:
        total_time = float("inf")
    
    start_time = time.time()    
    remaining_seconds = total_time

    fuzzers = discover_fuzzers(args.go_module_root)
    num_fuzzers = len(fuzzers)
    print(f"🐝 Discovered {num_fuzzers} fuzzers:", file=sys.stderr)
    for fuzzfn, path in fuzzers.items():
        print(f"{fuzzfn} in {path}", file=sys.stderr)

    if num_fuzzers == 0:
        print(f"No fuzzers found, this is likely an error. Exiting.")
        exit(1)

    # run forever or until --seconds, with increasingly longer durations per fuzz run
    durations_seconds = itertools.chain([5, 10, 30, 90, 270], itertools.repeat(600))
    if args.ci:
        # In CI - default to 60s fuzzes for scheduled runs, and 45 seconds for everything else
        durations_seconds = [60] if os.getenv('GITHUB_EVENT_NAME') == 'scheduled' else [45]
        if args.seconds:
            # However, if seconds was specified, evenly divide total time among all fuzzers
            # leaving a 10 second buffer for processing/building time between fuzz runs
            actual_fuzz_time = total_time - (num_fuzzers * 10)
            if actual_fuzz_time <= 5 * num_fuzzers:
                print(f"Seconds (--seconds {args.seconds}) is too low to properly run fuzzers for 5sec each. Exiting.")
                exit(1)
            durations_seconds = [ actual_fuzz_time / num_fuzzers ]

    for duration_seconds in durations_seconds:
        print(f"🐝 Running each fuzzer for {duration_seconds}s before switching to next fuzzer", file=sys.stderr)
        for fuzzfn, path in fuzzers.items():
            elapsed_time = time.time() - start_time
            remaining_seconds = total_time - elapsed_time 
            
            if remaining_seconds <= 0:
                print(f"🐝 Time budget of {args.seconds}s is exhausted. Exiting.", file=sys.stderr)
                return

            next_duration_seconds = min(remaining_seconds, duration_seconds)
            print(f"🐝 Running {fuzzfn} in {path} for {next_duration_seconds}s (Elapsed: {elapsed_time:.2f}s, Remaining: {remaining_seconds:.2f}s)", file=sys.stderr)
            run_fuzzer(fuzzfn, path, next_duration_seconds, args.go_module_root)
            print(f"🐝 Completed running {fuzzfn} in {path} for {next_duration_seconds}s.", file=sys.stderr)

def discover_fuzzers(go_module_root):
    fuzzers = {}
    for root, dirs, files in os.walk(go_module_root):
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
                fuzzers[fuzzfn] = os.path.relpath(root, go_module_root)
    return fuzzers

def run_fuzzer(fuzzfn, dir, duration_seconds, go_module_root):
    subprocess.check_call(["go", "test", "-run=^$", f"-fuzz=^{fuzzfn}$", f"-fuzztime={duration_seconds}s", f"./{dir}"], cwd=go_module_root)

if __name__ == "__main__":
    main()