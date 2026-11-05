#!/usr/bin/env bash
set -e

PROJECT_NAME="rbooks"
YEAR="2026"
AUTHOR="Madang Business Solutions LLC"

OUTPUT="NOTICE"

echo "Generating NOTICE file..."

# Header
cat > $OUTPUT <<EOF
$PROJECT_NAME
Copyright (c) $YEAR $AUTHOR

This product includes software developed by $AUTHOR.

This product includes third-party components licensed under the
Apache License 2.0 and MIT License. Attribution for these components
is provided through their respective LICENSE files within the source tree.

Notable third-party components include:
EOF

# Generate license list using cargo-license
cargo license --json | jq -r '.[] | "\(.name) (\(.license))"' | while read -r line; do
    NAME=$(echo "$line" | cut -d '(' -f 1 | xargs)
    LICENSE=$(echo "$line" | sed -n 's/.*(\(.*\)).*/\1/p')

    # Only include MIT / Apache-2.0 crates
    if [[ "$LICENSE" == *"MIT"* || "$LICENSE" == *"Apache-2.0"* ]]; then
        echo "- $NAME ($LICENSE)" >> $OUTPUT
    fi
done

# Footer
cat >> $OUTPUT <<EOF

The contents of this NOTICE file are for informational purposes only
and do not modify the Apache License 2.0.
EOF

echo "NOTICE file generated: $OUTPUT"
