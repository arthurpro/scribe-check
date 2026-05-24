package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	reset  = "\x1b[0m"
	bold   = "\x1b[1m"
	dim    = "\x1b[2m"
	red    = "\x1b[31m"
	yellow = "\x1b[33m"
	cyan   = "\x1b[36m"
	green  = "\x1b[32m"
)

func renderReport(w io.Writer, r Report) {
	color := isTerminal(w)
	cw := func(c, s string) string {
		if !color {
			return s
		}
		return c + s + reset
	}

	total := len(r.QuotesFlagged) + len(r.NamesFlagged) + len(r.SpecificsFlagged) +
		len(r.StyleDrift) + len(r.TemporalDrift)

	fmt.Fprintln(w)
	if total == 0 {
		fmt.Fprintln(w, cw(green+bold, "✓ scribe-check: no findings"))
	} else {
		fmt.Fprintf(w, "%s %d finding(s)\n", cw(red+bold, "⚑ scribe-check:"), total)
	}
	fmt.Fprintln(w)

	renderCategory(w, color, cw(red, "QUOTES FLAGGED"), r.QuotesFlagged, func(f Finding) (string, string, string) {
		return f.Passage, f.LocationHint, f.ClosestSourceMatch
	})
	renderCategory(w, color, cw(red, "NAMES FLAGGED"), r.NamesFlagged, func(f Finding) (string, string, string) {
		return f.Entity, "", ""
	})
	renderCategory(w, color, cw(yellow, "SPECIFICS FLAGGED"), r.SpecificsFlagged, func(f Finding) (string, string, string) {
		return f.Claim, "", ""
	})
	renderCategory(w, color, cw(cyan, "STYLE DRIFT"), r.StyleDrift, func(f Finding) (string, string, string) {
		return f.Passage, "", ""
	})
	renderCategory(w, color, cw(cyan, "TEMPORAL DRIFT"), r.TemporalDrift, func(f Finding) (string, string, string) {
		return f.Passage, "", ""
	})

	if strings.TrimSpace(r.Summary) != "" {
		fmt.Fprintln(w, bold+"SUMMARY"+reset)
		fmt.Fprintln(w, wrap(r.Summary, 88))
		fmt.Fprintln(w)
	}
}

func renderCategory(w io.Writer, color bool, header string, findings []Finding, extract func(Finding) (primary, hint, match string)) {
	if len(findings) == 0 {
		return
	}
	fmt.Fprintf(w, "%s  (%d)\n", header, len(findings))
	for i, f := range findings {
		primary, hint, match := extract(f)
		fmt.Fprintf(w, "  %d. %s\n", i+1, truncate(primary, 200))
		if hint != "" {
			fmt.Fprintf(w, "     %sat:%s %s\n", dimOnly(color), resetOnly(color), truncate(hint, 120))
		}
		fmt.Fprintf(w, "     %sconcern:%s %s\n", dimOnly(color), resetOnly(color), truncate(f.Concern, 200))
		if match != "" && match != "(none)" {
			fmt.Fprintf(w, "     %sclosest:%s %s\n", dimOnly(color), resetOnly(color), truncate(match, 200))
		}
	}
	fmt.Fprintln(w)
}

func dimOnly(color bool) string {
	if !color {
		return ""
	}
	return dim
}

func resetOnly(color bool) string {
	if !color {
		return ""
	}
	return reset
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func wrap(s string, width int) string {
	words := strings.Fields(s)
	var lines []string
	var cur strings.Builder
	for _, w := range words {
		if cur.Len() > 0 && cur.Len()+1+len(w) > width {
			lines = append(lines, cur.String())
			cur.Reset()
		}
		if cur.Len() > 0 {
			cur.WriteByte(' ')
		}
		cur.WriteString(w)
	}
	if cur.Len() > 0 {
		lines = append(lines, cur.String())
	}
	return strings.Join(lines, "\n")
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
