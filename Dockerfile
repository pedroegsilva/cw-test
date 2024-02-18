#==================================================
# Build Layer
FROM golang:1.21.7-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY parser/*.go ./parser/
COPY reports/*.go ./reports/
COPY main.go ./

RUN go build -o ./main

#==================================================
# Run Layer

FROM golang:1.21.7-alpine 

WORKDIR /app

COPY --from=build /app/main /app/main

COPY input/* /app/input/

CMD [ "/app/main", "-i", "/app/input/qgames.log" ]
