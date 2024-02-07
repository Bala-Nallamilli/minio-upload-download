#!/bin/bash

# Number of requests to send
REQUEST_COUNT=300

# Base URL of your server
BASE_URL="http://localhost:6000"

# Temporary files to store responses
BLOCK_RESPONSES=$(mktemp)
UNBLOCK_RESPONSES=$(mktemp)

# Function to send block requests and store responses
send_block_requests() {
    for ((i=1; i<=REQUEST_COUNT; i++))
    do
        (curl -s -X POST "$BASE_URL/blockRequest/$i" && echo " $i: Success") >> $BLOCK_RESPONSES || echo " $i: Fail" >> $BLOCK_RESPONSES &
    done
}

# Function to send unblock requests and store responses
send_unblock_requests() {
    for ((i=1; i<=REQUEST_COUNT; i++))
    do
        (curl -s -X POST "$BASE_URL/unblockRequest/$i" && echo " $i: Success") >> $UNBLOCK_RESPONSES || echo " $i: Fail" >> $UNBLOCK_RESPONSES &
    done
}

echo "Sending block requests..."
send_block_requests
wait
echo "Sending unblock requests..."
send_unblock_requests
wait

echo "Block Requests Results:"
cat $BLOCK_RESPONSES | sort -n

echo "Unblock Requests Results:"
cat $UNBLOCK_RESPONSES | sort -n

# Cleanup
rm $BLOCK_RESPONSES
rm $UNBLOCK_RESPONSES
