# Traffic analyzer
Simple golang app to analyze network traffic for TCP retransmission.

# Build
```sh
docker build -t traffic_analyzer .
```

# Run
```sh
docker run --network host -d traffic_analyzer <interface-to-capture>
```