# WebSocket Loadtester
Command line tool written in Go to open many connections with a websocket server for load testing. 
This tool is designed for read-only WebSockets (ie where the client does not send any messages to the server). 
It allows you to specify the URL to connect to, and the number of concurrent connections to open.

## Usage

Build the docker image with:
```
docker build -t ws-loadtester .
```
The docker build uses multi-stage build to give a small image size (~20mb).

Run a container with:
```
docker run --rm ws-loadtester
```
The executable in the container is named `/loadtester`, so more command line options can be seen by running:
```
docker run --rm ws-loadtester /loadtester --help
```
```
Usage of /loadtester:
  -duration string
        Specify duration of test. Each connection will stay connected for this duration. (default "10m")
  -print
        Use --print to print a dot for each message received on each connection. False to only print for one channel.
  -qty int
        Specify quantity of concurrent connections. (default 100)
  -url string
        Specify the url of the WebSocket. Should begin with ws:// or wss://. (default "wss://companies.stream/events")
```
