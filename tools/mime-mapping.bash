#!/bin/bash
set -eu

MIME_TYPES="/etc/mime.types"

if [ $# -gt 1 ]; then
    echo "Usage: $0 [<mime.types file>]"
    exit 1
fi

if [ $# -eq 1 ]; then
    MIME_TYPES="$1"
fi

if [ ! -f "${MIME_TYPES}" ]; then
    echo "Not found mime.types file: ${MIME_TYPES}"
    exit 1
fi

declare -A ext_uniq

sed '
s/#.*$//;
s/^[ \t]*//;
s/[ \t]*$//;
/^$/d;
s/[ \t][ \t]*/ /g;
/^\(text\|application\|image\|audio\|video\)\//!d;
' /etc/mime.types | while read mtype ext; do
  [ -z "${mtype}" ] && continue
  [ -z "${ext}" ] && continue

  declare -a e_lst=()
  while read -d ' ' e; do
    if [[ -z ${ext_uniq[${e}]:-} ]]; then
      e_lst+=("${e}")
    fi

    ext_uniq["${e}"]="ok"
  done <<<"${ext} "

  if [[ ${#e_lst[@]} -eq 0 ]]; then
    continue
  fi

  pat='\.'
  for ext in "${e_lst[@]}"; do
    if [[ ! ${ext} =~ ${pat} ]]; then
      printf "%s %s\n" "${ext}" "${mtype}"
    fi
  done
done | sort -u | sort -k 2 -t " "
