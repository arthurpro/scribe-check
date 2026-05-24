package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	defaultModel     = "gemma4:latest"
	gemma4MaxContext = 131072 // Gemma 4's full 128K window
)

func main() {
	var (
		model      = flag.String("model", defaultModel, "Ollama model tag")
		host       = flag.String("host", "http://localhost:11434", "Ollama host")
		jsonOnly   = flag.Bool("json", false, "print raw JSON output only (no terminal table)")
		outFile    = flag.String("out", "", "write JSON output to this file in addition to stdout")
		maxTokens  = flag.Int("max-tokens", 120000, "reject if article + sources exceed this token estimate")
		numCtx     = flag.Int("num-ctx", 0, "Ollama context window in tokens (0 = auto-size from content, up to 131072)")
		quiet      = flag.Bool("quiet", false, "suppress progress prints to stderr")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `scribe-check — fabrication-review for Markdown articles via Gemma 4 on Ollama

Usage:
  scribe-check [flags] <article.md> <sources-path>

  <sources-path> may be a single file or a directory of .md/.txt/.html sources.

Flags:
`)
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}
	articlePath, sourcesPath := flag.Arg(0), flag.Arg(1)

	logf := func(format string, args ...any) {
		if !*quiet {
			fmt.Fprintf(os.Stderr, format, args...)
		}
	}

	logf("scribe-check: loading article %s\n", articlePath)
	article, err := loadArticle(articlePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load article: %v\n", err)
		os.Exit(1)
	}

	logf("scribe-check: loading sources from %s\n", sourcesPath)
	sources, err := loadSources(sourcesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load sources: %v\n", err)
		os.Exit(1)
	}
	if len(sources) == 0 {
		fmt.Fprintf(os.Stderr, "no source files found under %s\n", sourcesPath)
		os.Exit(1)
	}
	logf("scribe-check: %d source file(s) loaded\n", len(sources))

	estTokens := estimateTokens(article, sources)
	if estTokens > *maxTokens {
		fmt.Fprintf(os.Stderr, "estimated %d tokens exceeds --max-tokens=%d\n", estTokens, *maxTokens)
		os.Exit(1)
	}
	chosenCtx := chooseNumCtx(*numCtx, estTokens)
	logf("scribe-check: ~%d tokens (limit %d), num_ctx=%d, calling %s\n", estTokens, *maxTokens, chosenCtx, *model)

	report, raw, err := runReview(*host, *model, chosenCtx, article, sources)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		os.Exit(1)
	}

	if *outFile != "" {
		if err := os.WriteFile(*outFile, []byte(raw), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write --out: %v\n", err)
			os.Exit(1)
		}
		logf("scribe-check: wrote %s\n", *outFile)
	}

	if *jsonOnly {
		fmt.Println(strings.TrimSpace(raw))
		return
	}

	renderReport(os.Stdout, report)
	if hasFindings(report) {
		os.Exit(1)
	}
}

// chooseNumCtx picks an Ollama num_ctx that fits the estimated content with
// ~25% headroom for the system prompt and the JSON completion, rounded up to
// a power-of-two bucket so KV-cache allocations stay tidy. Capped at Gemma 4's
// native 131072. If the user passed --num-ctx explicitly, that wins.
func chooseNumCtx(explicit, estimated int) int {
	if explicit > 0 {
		if explicit > gemma4MaxContext {
			return gemma4MaxContext
		}
		return explicit
	}
	target := estimated + estimated/4 + 2048
	for _, b := range []int{4096, 8192, 16384, 32768, 65536, gemma4MaxContext} {
		if b >= target {
			return b
		}
	}
	return gemma4MaxContext
}
