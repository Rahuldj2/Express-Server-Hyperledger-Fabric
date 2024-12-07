#!/bin/bash

# Define paths
CURRENT_DIR=$(pwd)
PARENT_DIR=$(dirname "$CURRENT_DIR")
ORG1_SOURCE="$PARENT_DIR/organizations/peerOrganizations/org1.example.com/connection-org1.json"
ORG2_SOURCE="$PARENT_DIR/organizations/peerOrganizations/org2.example.com/connection-org2.json"
ORG1_DEST="$CURRENT_DIR/config/connection-org1.json"
ORG2_DEST="$CURRENT_DIR/config/connection-org2.json"

# Overwrite the destination files with the source files
if [[ -f "$ORG1_SOURCE" ]]; then
  cp "$ORG1_SOURCE" "$ORG1_DEST"
  echo "Copied connection-org1.json successfully."
else
  echo "Source file for Org1 not found: $ORG1_SOURCE"
fi

if [[ -f "$ORG2_SOURCE" ]]; then
  cp "$ORG2_SOURCE" "$ORG2_DEST"
  echo "Copied connection-org2.json successfully."
else
  echo "Source file for Org2 not found: $ORG2_SOURCE"
fi
