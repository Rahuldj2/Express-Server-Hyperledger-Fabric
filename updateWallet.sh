#!/bin/bash

# Define paths
CURRENT_DIR=$(pwd)
WALLET_DIR="$CURRENT_DIR/wallet"

# Check if wallet directory exists
if [[ -d "$WALLET_DIR" ]]; then
  # Delete all files inside wallet directory
  rm -rf "$WALLET_DIR"/*
  echo "Deleted all files in the wallet directory."
else
  echo "Wallet directory does not exist: $WALLET_DIR"
  exit 1
fi

# Run the Node.js scripts
echo "Running importorg1admin.js..."
node importorg1admin.js
if [[ $? -ne 0 ]]; then
  echo "Error running importorg1admin.js"
  exit 1
fi

echo "Running importorg2admin.js..."
node importorg2admin.js
if [[ $? -ne 0 ]]; then
  echo "Error running importorg2admin.js"
  exit 1
fi

echo "Scripts executed successfully!"
