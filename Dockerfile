FROM keinos/sqlite3 AS sqlite

WORKDIR "/data"

RUN /usr/bin/sqlite3 store.db "CREATE TABLE IPStat (IP Text, All_ Integer, Retransmitted Integer)"

FROM golang AS build

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=1
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

WORKDIR /app

COPY . .
    
RUN apt update -y && \
    apt install -y libpcap-dev && \
    go mod download  && \
    go get -u && \
    go build -o ./app -ldflags="-s -w" ./main.go

FROM golang AS final

WORKDIR /app

COPY --from=SQLITE /data/store.db ./store.db
COPY --from=build /app/app ./traffic_analyzer

RUN apt update -y && \
    apt install -y libpcap-dev 

EXPOSE 8080/tcp

ENTRYPOINT ["/app/traffic_analyzer"]