set shell := ["bash", "-uc"]

# List recipes.
default:
    @just --list

# Build the binary.
build:
    go build ./...

# Tidy go.mod.
tidy:
    go mod tidy

# Remove the compiled binary.
clean:
    rm -f scribe-check

# Run scribe-check against an article and a sources path.
check ARTICLE SOURCES: build
    ./scribe-check {{ARTICLE}} {{SOURCES}}

# Demo: clean article (expected: zero or borderline findings).
demo-clean: build
    ./scribe-check ./examples/good-article.md ./examples/sources/

# Demo: planted fabrications (expected: catches Petrova, *elementary cells*, 240M rod cells).
demo-fab: build
    ./scribe-check ./examples/article-with-fabrications.md ./examples/sources/

# Run both demos back-to-back.
demo: demo-clean demo-fab

# Regenerate the captured JSON + transcript artifacts under examples/.
refresh-examples: build
    ./scribe-check --json ./examples/article-with-fabrications.md ./examples/sources/ > examples/output-fabrications.json
    ./scribe-check --json ./examples/good-article.md ./examples/sources/ > examples/output-clean.json
    ./scribe-check ./examples/article-with-fabrications.md ./examples/sources/ > examples/transcript-fabrications.txt 2>&1 || true

# Verify Ollama has gemma4:latest pulled.
check-ollama:
    @curl -sf http://localhost:11434/api/tags | grep -q '"gemma4:latest"' && echo "gemma4:latest is present" || (echo "gemma4:latest not found — run: ollama pull gemma4:latest" && exit 1)
