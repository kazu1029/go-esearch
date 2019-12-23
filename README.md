# go-esearch

## How to use

`docker-compose up`

`curl -X POST http://localhost:8080/documents -d @sample_japanese.json -H "Content-Type: application/json"`

`curl http://localhost:8080/search?query=おはよう+タイトル`
