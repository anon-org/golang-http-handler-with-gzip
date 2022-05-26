#!/bin/sh

set -euo pipefail

echo '{"number": 99}' | gzip | \
curl -iXPOST http://localhost:8000/ \
-H "Content-Type: application/json" \
-H "Content-Encoding: gzip" \
--compressed --data-binary @-
