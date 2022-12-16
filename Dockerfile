FROM public.ecr.aws/docker/library/golang:1.19 as golang

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=js GOARCH=wasm go build -o chip-8.wasm backend/wasm/main.go
RUN GOOS=linux GOARCH=amd64 go build -o webserver backend/server/main.go

FROM node:18 as node

COPY frontend .
RUN npm install

COPY --from=golang /go/src/app/chip-8.wasm public/

RUN npm test -- --watchAll=false --passWithNoTests
RUN npm run build

FROM ubuntu:22.04
WORKDIR /home/

COPY --from=node build ./build
COPY --from=golang /go/src/app/webserver .

CMD ./webserver
