#!/usr/bin/env sh

printf "[client] order sent, waiting for response...\n"
curl -X POST --location "http://localhost:8080/api/v1/orders" \
    -H "Content-Type: application/json" \
    -d '{"product_id": 1234, "price": "50000", "count": 2}'