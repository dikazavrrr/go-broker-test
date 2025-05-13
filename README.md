# 📈 Go Broker Test

A simple Go server that handles trade data, stores it in SQLite, and calculates per-account profit statistics.

---

## 🚀 How to Run

### 1. Build & Start Services

```bash
docker-compose build
docker-compose up
```

### 2. Start the Server & Worker
Start API server in terminal 1
```
go run ./cmd/server --db data.db --listen 8080
```
or
``` 
make run-server
``` 

Start background worker in terminal 2
```
go run ./cmd/worker --db data.db --poll 100ms
```
```
make run-worker
```


### 3. Making Requests
## Sample POST request
To add a new trade, use the following curl command
```
curl -X POST http://localhost:8080/trades \
     -H 'Content-Type: application/json' \
     -d '{
           "account": "127",
           "symbol": "EURUSD",
           "volume": 1.1,
           "open": 1.1000,
           "close": 1.1050,
           "side": "buy"
         }'
```

## Sample GET request
To fetch account stats

```
curl -X GET http://localhost:8080/stats/127
# {"Account":"127","Trades":1,"Profit":550}
```

## Sample healthz checker
To check the health of the API
```
curl -X GET http://localhost:8080/healthz
```

## Sample worker request
To query the worker for statistics
```
curl http://localhost:8080/stats/127
# {"Account":"127","Trades":1,"Profit":550}
```