#!/bin/bash

_FILE=~/.im_records
_DATE="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
_DESC="$@"

usage () {
  echo "IM: Tracks what I'm doing"
  echo
  echo "  Usage:"
  echo
  echo "    \$ im <description of activity>"
  echo
  echo "  Flags:"
  echo
  echo "    -t          list tags"
  echo "    -g <tag>    filter by tag"
  echo "    -s          show recent records"
  echo "    -d          delete most recent record"
  echo "    -h          show this help"
  echo
  echo "  Examples:"
  echo
  echo "    \$ im working on my @app"
  echo "    \$ im slacking"
  echo
}

if [ $# -eq 0 ]; then
  usage
  exit 1
fi

while getopts "g:tsdh" opt; do
    case "$opt" in
    t)  # list all tags
        cat "$_FILE" | grep -oh '@\w*' | cut -d@ -f2 | sort -u
        exit 0
        ;;
    g)  # grep by tag
        cat "$_FILE" | grep "@${OPTARG}"
        exit 0
        ;;
    s)  # show revent records
        tac ~/.im_records | head
        exit 0
        ;;
    d)  # delete previous record
        _TMP_FILE="$(mktemp)"
        head --lines=-1 "$_FILE" > "$_TMP_FILE"
        cp "$_TMP_FILE" "$_FILE"
        rm "$_TMP_FILE"
        exit 0
        ;;
    h|*)
        usage
        exit 1
        ;;
    esac
done

echo "$_DATE: $_DESC" >> "$_FILE"




