#!/bin/bash

_RECORD_FILE=~/.im_records
_RECORD_DATE="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
_RECORD_DESC="$@"

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
        cat "$_RECORD_FILE" | grep -oh '@\w*' | cut -d@ -f2 | sort -u
        exit 0
        ;;
    g)  # grep by tag
        cat "$_RECORD_FILE" | grep "@${OPTARG}"
        exit 0
        ;;
    s)  # show revent records
        tac ~/.im_records | head
        exit 0
        ;;
    d)  # delete previous record
        _TMP_FILE="$(mktemp)"
        head --lines=-1 "$_RECORD_FILE" > "$_TMP_FILE"
        cp "$_TMP_FILE" "$_RECORD_FILE"
        rm "$_TMP_FILE"
        exit 0
        ;;
    h|*)
        usage
        exit 1
        ;;
    esac
done

echo "$_RECORD_DATE: $_RECORD_DESC" >> "$_RECORD_FILE"




