<p align="center">
<h1 align="center">Shorty</h1>
<p align="center">Simple link shortener </p>
<p align="center">This project was created for training purposes to practice practical tasks and showcase my skills</p>

## Features

  * Materialistic http rest api
    * Creating short links
    * View statistics on short link visits
    * Redirecting users to real url addresses

## Build
Clone the repository and compile the application

```bash
# Clone repo
git clone https://github.com/need-o/shorty.git

#Build app
cd shorty
make build
```

## Usage
You can use environment variables to change the default address to the database and the server address
```bash
# Export environment variables
export SHORTY_DB_PATH=shorty.db
export SHORTY_ADDRESS=:1323

# Run shorty
./shorty
```

#### Creating a short link

```bash
curl -d '{"url": "https://example.com" }' -H "Content-Type: application/json" -X POST http://localhost:1323/api/shorty
```

```json
{
    "id":"sE5qwb",
    "address":"localhost:1323/sE5qwb"
}
```

#### Follow the short link

```bash
curl http://localhost:1323/sE5qwb -verbose 
```

```
* processing: http://localhost:1323/sE5qwb
*   Trying [::1]:1323...
* Connected to localhost (::1) port 1323
> GET /sE5qwb HTTP/1.1
> Host: localhost:1323
> User-Agent: curl/8.2.1
> Accept: */*
> Referer: rbose
> 
< HTTP/1.1 301 Moved Permanently
< Location: https://example.com
< X-Request-Id: ex3lMGKXYoOO5vw1ZjfPMmxUsVMg4MqG
< Date: Mon, 28 Aug 2023 11:28:50 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```
#### View visit statistics

```bash
curl http://localhost:1323/api/shorty/sE5qwb
```

```json
{
    "id": "sE5qwb",
    "url": "https://example.com",
    "visits": [
        {
            "shorty_id": "sE5qwb",
            "referer": "",
            "user_ip": "127.0.0.1",
            "user_agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/116.0",
            "created_at": "2023-08-28T14:56:47.898469006+03:00",
            "updated_at": "2023-08-28T14:56:47.898469006+03:00"
        }
    ],
    "created_at": "2023-08-28T14:23:31.413414468+03:00",
    "updated_at": "2023-08-28T14:23:31.413414468+03:00"
}
```

## Testing

To run tests, use the default `go test` command:
```sh
go test ./...
```

## License

Shprty released under GNU GPL license, refer [LICENSE](LICENSE) file.
