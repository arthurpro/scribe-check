package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Article struct {
	Path string
	Body string
}

type Source struct {
	Path string
	Body string
}

func loadArticle(path string) (Article, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Article{}, err
	}
	return Article{Path: path, Body: string(b)}, nil
}

func loadSources(path string) ([]Source, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return []Source{{Path: path, Body: string(b)}}, nil
	}

	var sources []Source
	err = filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(p))
		switch ext {
		case ".md", ".txt", ".html", ".htm":
			b, err := os.ReadFile(p)
			if err != nil {
				return fmt.Errorf("read %s: %w", p, err)
			}
			sources = append(sources, Source{Path: p, Body: string(b)})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(sources, func(i, j int) bool { return sources[i].Path < sources[j].Path })
	return sources, nil
}

// estimateTokens is a coarse 4-chars-per-token approximation, deliberately
// pessimistic on the high end so we err on the side of fitting.
func estimateTokens(a Article, sources []Source) int {
	total := len(a.Body)
	for _, s := range sources {
		total += len(s.Body)
	}
	return total / 4
}
