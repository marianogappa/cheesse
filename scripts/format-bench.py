#!/usr/bin/env python3
"""Parses Go benchmark output and produces a human-readable Markdown table."""
import sys, re

LABELS = {
    "BenchmarkCalculateAllActions_Start": "Generate legal moves (opening)",
    "BenchmarkCalculateAllActions_Kiwipete": "Generate legal moves (complex position)",
    "BenchmarkDoAction_Start": "Apply a move (opening)",
    "BenchmarkDoAction_Kiwipete": "Apply a move (complex position)",
    "BenchmarkNewGameFromFEN": "Parse a FEN string",
    "BenchmarkPerft3_Start": "Perft(3) from starting position",
    "BenchmarkAIDepth0_Start": "AI move, Easy (depth 0, opening)",
    "BenchmarkAIDepth1_Start": "AI move, Medium (depth 1, opening)",
    "BenchmarkAIDepth2_Endgame": "AI move, Hard (depth 2, endgame)",
    "BenchmarkParseNotation_50MoveGame": "Parse 50-move game (auto-detect)",
    "BenchmarkConvertNotation_ToICCF": "Convert 13-move game to ICCF",
}

def fmt_time(ns_str):
    ns = float(ns_str)
    if ns >= 1_000_000_000:
        return f"{ns / 1_000_000_000:.2f}s"
    if ns >= 1_000_000:
        return f"{ns / 1_000_000:.1f}ms"
    if ns >= 1_000:
        return f"{ns / 1_000:.0f}us"
    return f"{ns:.0f}ns"

def fmt_bytes(b_str):
    b = int(b_str)
    if b >= 1_000_000:
        return f"{b / 1_000_000:.1f} MB"
    if b >= 1_000:
        return f"{b / 1_000:.0f} KB"
    return f"{b} B"

def fmt_allocs(a_str):
    a = int(a_str)
    if a >= 1_000_000:
        return f"{a / 1_000_000:.1f}M"
    if a >= 1_000:
        return f"{a / 1_000:.1f}K"
    return str(a)

print("| What | Time | Memory | Allocations |")
print("|---|---|---|---|")

for line in sys.stdin:
    line = line.strip()
    if not line.startswith("Benchmark"):
        continue
    parts = line.split()
    # Go bench format: Name-N iters time ns/op bytes B/op allocs allocs/op
    name = re.sub(r'-\d+$', '', parts[0])
    time_ns = parts[2]
    bytes_ = parts[4]
    allocs = parts[6]
    label = LABELS.get(name, name)
    print(f"| {label} | {fmt_time(time_ns)} | {fmt_bytes(bytes_)} | {fmt_allocs(allocs)} |")
