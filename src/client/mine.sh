#!/usr/bin/env bash

INPUT=$1
OUTPUT=$2
CWD=$(pwd)

python3 \
-m gspan_mining \
-s 100 \
-d True \
-l 3 \
-u 3 \
-w True \
"$CWD"/"$INPUT" >"$CWD"/"$OUTPUT"
