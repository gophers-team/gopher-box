curl -v -X POST http://127.0.0.1:8000/status \
    -H "Contenty-Type: application/json" \
    -d '{"device_id": 1}' 

sqlite3 gopher-box.db 'SELECT * from heartbeats WHERE device_id = 1'

