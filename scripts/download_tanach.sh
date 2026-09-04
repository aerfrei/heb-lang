#!/usr/bin/env bash
# Downloads the individual Tanach book XML files from tanach.us into tanach/.
set -euo pipefail

BOOKS=(
  Genesis Exodus Leviticus Numbers Deuteronomy
  Joshua Judges Samuel_1 Samuel_2 Kings_1 Kings_2
  Isaiah Jeremiah Ezekiel
  Hosea Joel Amos Obadiah Jonah Micah Nahum Habakkuk Zephaniah Haggai Zechariah Malachi
  Chronicles_1 Chronicles_2
  Psalms Job Proverbs Ruth Song_of_Songs Ecclesiastes Lamentations Esther Daniel
  Ezra Nehemiah
)

OUT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/tanach"
mkdir -p "$OUT_DIR"

for book in "${BOOKS[@]}"; do
  echo "Downloading ${book}..."
  curl -fsSL "https://www.tanach.us/Books/${book}.xml" -o "${OUT_DIR}/${book}.xml"
done

echo "Done. Saved ${#BOOKS[@]} books to ${OUT_DIR}"
