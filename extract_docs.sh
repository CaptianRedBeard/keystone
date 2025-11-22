#!/usr/bin/env bash

set -euo pipefail

OUTPUT="PROJECT_DOCS.md"
echo "# Project Documentation" > "$OUTPUT"
echo "" >> "$OUTPUT"

# Walk through all .go files while ignoring vendor and generated dirs
FILES=$(find . -type f -name "*.go" ! -path "./vendor/*" ! -path "*/mocks/*" ! -path "*/generated/*")

for file in $FILES; do
    echo "Processing $file"

    echo "## ${file#./}" >> "$OUTPUT"
    echo '```go' >> "$OUTPUT"

    # Extract only doc-comments: blocks of // comments followed immediately by a declaration
    awk '
        BEGIN { doc=""; inComment=0 }
        {
            # If we see a comment, accumulate it
            if ($0 ~ /^[[:space:]]*\/\//) {
                doc = doc $0 "\n"
                inComment=1
            }
            # If previous lines were comments and now we see a declaration, print them
            else if (inComment == 1 && $0 ~ /^[[:space:]]*(type|func|var|const|package)[[:space:]]+/) {
                print doc $0 "\n"
                doc=""
                inComment=0
            }
            # Reset if we hit a non-doc-comment line
            else {
                doc=""
                inComment=0
            }
        }
    ' "$file" >> "$OUTPUT"

    echo '```' >> "$OUTPUT"
    echo "" >> "$OUTPUT"
done

echo "Documentation extracted to $OUTPUT"
