#!/usr/bin/env bash

get_env() {
  local key="$1"
  local default="$2"
  if [[ -f "${ENV_FILE}" ]]; then
    local value
    value="$(grep -E "^${key}=" "${ENV_FILE}" | tail -1 | cut -d= -f2-)"
    if [[ -n "${value}" ]]; then
      printf '%s' "${value}"
      return
    fi
  fi
  printf '%s' "${default}"
}

require_env() {
  local key="$1"
  local value
  value="$(get_env "${key}" "")"
  if [[ -z "${value}" ]]; then
    echo "missing ${key} in ${ENV_FILE}" >&2
    exit 1
  fi
  printf '%s' "${value}"
}
