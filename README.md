# scribe-check

A local-first fabrication-review CLI for Markdown articles, powered by Gemma 4 8B via Ollama. You hand it a draft and a folder of sources; it tells you which concrete claims in the draft aren't corroborated by the sources you handed it.

It runs on a laptop. It doesn't call any cloud. It doesn't use RAG. The whole article and all sources go into the model in one prompt — Gemma 4 8B's 128K context makes this practical.

## Why this exists

Every technical writer fabricates by accident. A quote drifts one word. A coauthor's name creeps in from a different paper. A year gets transposed. An off-by-one ratio survives copy-edit. `scribe-check` is the kind of pre-publish pass an editor would do, run by a small local model, fast enough to run on every draft.

## Requirements

- Go 1.22+
- [Ollama](https://ollama.com/) running locally with `gemma4:latest` pulled (`ollama pull gemma4:latest`)
- ~10 GB free RAM for the 8B Q4_K_M model
- [`just`](https://github.com/casey/just) (optional, for the recipes)

## Install

```bash
git clone <repo> && cd scribe-check
just build       # or: go build ./...
```

## Quickstart

```bash
just demo-fab    # planted-fabrication article — should flag Petrova, *elementary cells*, 240M rods
just demo-clean  # clean article — should be quiet except for borderline over-flags
just demo        # both, back-to-back
```

Or call the binary directly:

```bash
./scribe-check examples/article-with-fabrications.md examples/sources/
```

Other recipes: `just --list`.

## Usage

```
scribe-check [flags] <article.md> <sources-path>
```

`<sources-path>` is either a single source file or a directory of `.md`/`.txt`/`.html` sources.

Flags:

| Flag | Default | Purpose |
|---|---|---|
| `--model` | `gemma4:latest` | Ollama model tag |
| `--host` | `http://localhost:11434` | Ollama host |
| `--num-ctx` | `0` (auto) | Ollama context window in tokens. `0` auto-sizes from estimated content (4K → 8K → 16K → 32K → 64K → 128K buckets, capped at Gemma 4's native 131072). |
| `--max-tokens` | `120000` | Reject if article + sources exceed this token estimate (leaves headroom under the 131072 model max) |
| `--json` | `false` | Print raw JSON only (no terminal table) |
| `--out` | `""` | Also write raw JSON to this file |
| `--quiet` | `false` | Suppress progress prints to stderr |

Exit codes: `0` clean, `1` findings present, `2` usage error.

### About the 128K context window

`gemma4:latest` exposes a 131072-token context (Ollama reports `gemma4.context_length = 131072`). By default `scribe-check` allocates *only as much KV cache as the content needs*, in power-of-two buckets up to the full 128K. A 3,000-word article + a curated citations file fits comfortably in the 16K bucket. A long-form essay plus a folder of raw source HTML can grow into the 64K or 128K buckets without any code change — set `--num-ctx 131072` explicitly to pre-allocate the full window. Latency scales roughly linearly with the chosen bucket; auto-sizing keeps the writer's loop tight on the common case and only spends the full window when the input actually demands it.

## What it flags

Five categories. The system prompt instructs the model to scan each separately:

1. **quotes_flagged** — italicized or quoted passages whose verbatim text doesn't appear in any source. Especially: terminology drift (article italicizes `*X*`, sources italicize `*Y*`).
2. **names_flagged** — people, products, papers, specifications mentioned in the article but absent from the sources. Catches fabricated coauthors, invented company names.
3. **specifics_flagged** — numbers, dates, years, percentages, versions, ratios — anything with a digit. Catches off-by-N errors, transposed years, fabricated quantities.
4. **style_drift** — British orthography (`colour`, `behaviour`, `-ised`) in a US-English article, or vice-versa. Verbatim quotes preserve source spelling and are not flagged.
5. **temporal_drift** — `today`, `this morning`, `yesterday`, weekday names in articles that are otherwise evergreen — these are leak risks.

## Honest characterization of what it does and doesn't catch

`scribe-check` is high-recall, modest-precision by design. The model is told to over-flag rather than under-flag. On real articles:

- It reliably catches **gross numeric drift** (120 → 240, off-by-2× errors, wrong years).
- It catches **fabricated names** when they appear in lists alongside corroborated names (a fake coauthor next to real ones).
- It catches **terminology drift** in italicized terms (`*elementary cells*` vs `*simple cells*`).
- It will **over-flag borderline cases**: derived numbers ("a twenty-to-one ratio" from 120M/6M), round-number paraphrases ("a thirty-year-old standard" from 1992), slightly-rephrased corroborated claims. Expect ~3–7 such false positives on a clean 3,000-word article. A human dismisses them in seconds; the cost of false-positive review is much cheaper than the cost of a missed fabrication.
- It does **not** catch logical contradictions (a sentence that contradicts another sentence in the same article).
- It does **not** catch subtle paraphrase drift below a word or two.

Treat the output as a checklist of items to glance at, not as a final judgment.

## How to use it well

- Hand it your `citations.md`, your fact-list, or the actual source HTML/markdown that backs the article. Sources can be a single curated file or a folder of raw materials.
- Run it as the last pass before publish. Skim the output. Dismiss the over-flags; investigate the rest.
- For longest articles, raise `--num-ctx` toward 131072 (Gemma 4's max). Latency scales roughly linearly with context.

## Example output

A planted-fabrication run on a real article:

```
⚑ scribe-check: 3+ finding(s)

QUOTES FLAGGED  (1)
  1. *elementary cells*
     at: They discovered that individual neurons in the primary visual cortex...
     concern: The article italicizes *elementary cells*, but the source uses
              *simple cells* as the term Hubel and Wiesel introduced.
     closest: structures they later called *simple cells*

NAMES FLAGGED  (1)
  1. Petrova
     concern: 'Petrova' is a fabricated coauthor. Source lists Ahmed, Natarajan,
              and Rao only.

SPECIFICS FLAGGED  (1+)
  1. The human eye contains roughly 240 million rod cells...
     concern: The article claims 240 million rod cells. The source corroborates
              roughly 120M rods. This is an off-by-two-fold error.
     closest: roughly 120M rods, ~6-7M cones per human retina
```

The full JSON for this run is in `examples/output-fabrications.json`. The clean-article run is in `examples/output-clean.json`.

## Why Gemma 4 8B

The job needs three things simultaneously and only the 8B variant has all three:

- **Structural reasoning** that 2B/4B don't reliably produce.
- **128K context** so a full article + 5–10 source documents fit in one prompt without RAG.
- **Local execution** because per-draft scans should cost nothing.

The 26B MoE and 31B dense variants do this better, but at 5–10× the latency. `scribe-check` is meant to run between drafts. On the 8B, that's seconds; on the 31B, it's long enough that you'd skip it.

## License

Apache 2.0.
