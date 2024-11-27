
# Receipt Processor API

This project implements a web service that processes receipts, calculates points based on specific rules, and allows you to retrieve points for a given receipt.

---

## Endpoints

### 1. **Process Receipts**

- **Path**: `/receipts/process`
- **Method**: `POST`
- **Description**: Processes a receipt, calculates points based on its data, and returns a unique ID for the receipt.

#### Request Body (Example):
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },
    {
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    }
  ],
  "total": "35.35"
}
```

#### Response (Example):
```json
{
  "id": "7fb1377b-b223-49d9-a31a-5a02701dd310"
}
```

---

### 2. **Get Points**

- **Path**: `/receipts/{id}/points`
- **Method**: `GET`
- **Description**: Retrieves the total points awarded to a receipt based on its unique ID.

#### Response (Example):
```json
{
  "points": 28
}
```

#### Notes:
- Replace `{id}` in the path with the ID returned by the `/receipts/process` endpoint.

---

## Running the Application with Docker

### 1. **Build the Docker Image**
Run the following command to build the Docker image:
```bash
docker build -t receipt-processor .
```

### 2. **Run the Docker Container**
Run the container with the following command:
```bash
docker run --rm -p 8080:8080 receipt-processor
```

### 3. **Access the API**
The API will be accessible at `http://localhost:8080`. You can use tools like `curl` or Postman to call the endpoints.

---

Enjoy using the Receipt Processor API! 🚀
