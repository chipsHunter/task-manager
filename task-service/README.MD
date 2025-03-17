# To run service
``` bash
docker run -d \
  -v $(pwd)/storage/storage.db:/app/storage/storage.db \
  -p 8082:8082 \
  --name task-service-container \
  task-service
```
# To build
``` bash
mkdir storage
touch ./storage/storage.db
docker image build -t task-service .
```
