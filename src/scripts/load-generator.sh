#!/bin/bash

# Number of transfer requests to send
TOTAL_REQUESTS=400
SLEEP_SECONDS=0.5

echo "Sending $TOTAL_REQUESTS transfer requests to /transfer every $SLEEP_SECONDS seconds..."

for ((i = 1; i <= TOTAL_REQUESTS; i++)); do
  curl -s -X POST http://localhost:30010/transfer \
    -H "Content-Type: application/json" \
    -d '{
      "fromAccount": "user1",
      "toAccount": "user2",
      "amount": 100,
      "currency": "USD"
    }' > /dev/null

  echo "Sent request #$i"
  sleep $SLEEP_SECONDS
done

echo "Done. $TOTAL_REQUESTS requests sent."
