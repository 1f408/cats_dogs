#!/bin/bash
printf "types {\n"

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
    if [[ -z ${ext_uniq[${e}]} ]]; then
      e_lst+=("${e}")
    fi

    ext_uniq["${e}"]="ok"
  done <<<"${ext} "

  if [[ ${#e_lst[@]} -eq 0 ]]; then
    continue
  fi

  printf "\t%s\t%s;\n" "${mtype}" "${e_lst[*]}"
done | sort -u | expand -t4,64
printf "}\n"
