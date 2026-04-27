#!/bin/sh

GREEN="\033[0;32m"
RED="\033[0;31m"
RESET="\033[0m"

set -e

# cleanup database files to reset to "factory settings" and have a clean state for integration test
echo "${RED}****************************************************************************************************************************************"
echo "${RED}   Deleting database files after new start to reset Databasus to factory settings and create a clean state to all integration tests"
echo "${RED}****************************************************************************************************************************************${RESET}"
rm -rf /databasus-data/*

# Start original entrypoint in background
/app/start.sh &
MAIN_PID=$!

# Wait for port 4005
echo "Waiting for service..."
until curl -sf http://localhost:4005 > /dev/null; do
  sleep 5
done

echo "Databasus is ready, running command to set admin pwd in 10 seconds..."
sleep 10

./main --new-password="${ADMIN_PASSWORD_TO_SET}" --email="admin"

# Fixed inner width of the box (adjust if you want)
WIDTH=83

MSG1="ADMIN PASSWORD HAS BEEN SET TO: ${ADMIN_PASSWORD_TO_SET}"
MSG2="You can now use the REST API with user 'admin' and the new password"

pad_line() {
  TEXT="$1"
  LEN=${#TEXT}
  PAD=$((WIDTH - LEN))

  # build padding spaces
  SPACES=$(printf "%${PAD}s" "")

  printf "${GREEN}***** %s%s *****\n" "$TEXT" "$SPACES"
}

echo "${GREEN}$(printf '%*s' $((WIDTH + 12)) '' | tr ' ' '*')"
pad_line "$MSG1"
pad_line "$MSG2"
echo "${GREEN}$(printf '%*s' $((WIDTH + 12)) '' | tr ' ' '*')${RESET}"

# Keep container alive and properly handle signals
wait $MAIN_PID