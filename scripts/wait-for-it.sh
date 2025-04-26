#!/usr/bin/env bash
# wait-for-it.sh: Wait for a service to become available

set -e

TIMEOUT=15
QUIET=0

usage() {
  echo "Usage: $0 host:port [-t timeout] [-- command args]"
  exit 1
}

wait_for() {
  local host=$1
  local port=$2
  local start_ts=$(date +%s)

  while :
  do
    if nc -z "$host" "$port"; then
      local end_ts=$(date +%s)
      echo "wait-for-it: $host:$port is available after $((end_ts - start_ts)) seconds"
      break
    fi

    local now_ts=$(date +%s)
    if ((now_ts - start_ts >= TIMEOUT)); then
      echo "wait-for-it: timeout occurred after waiting $TIMEOUT seconds for $host:$port"
      exit 1
    fi

    sleep 1
  done
}

while [[ $# -gt 0 ]]
do
  case "$1" in
    *:* )
    host_port=(${1//:/ })
    HOST=${host_port[0]}
    PORT=${host_port[1]}
    shift 1
    ;;
    -t)
    TIMEOUT=$2
    shift 2
    ;;
    --quiet)
    QUIET=1
    shift 1
    ;;
    --)
    shift
    CMD="$@"
    break
    ;;
    *)
    usage
    ;;
  esac
done

if [[ -z "$HOST" || -z "$PORT" ]]; then
  usage
fi

wait_for "$HOST" "$PORT"

if [[ -n "$CMD" ]]; then
  exec $CMD
fi