# Receipt Processor API

This is a simple REST API that processes receipts and calculates reward points based on specific rules.

## Features:
- Submit a receipt for processing (`POST /receipts/process`)
- Retrieve points for a receipt (`GET /receipts/{id}/points`)


# Run the serever
go run main.go


# Hitting the endpoints

curl -X POST "http://localhost:8080/receipts/process" \
     -H "Content-Type: application/json" \
     -d '{
           "retailer": "Target",
           "purchaseDate": "2022-01-01",
           "purchaseTime": "13:01",
           "items": [
             {"shortDescription": "Mountain Dew 12PK", "price": "6.49"},
             {"shortDescription": "Emils Cheese Pizza", "price": "12.25"}
           ],
           "total": "35.35"
         }'



curl -X GET "http://localhost:8080/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points"
