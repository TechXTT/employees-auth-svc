FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/auth-svc cmd/main.go

FROM alpine

COPY --from=0 /app/auth-svc /auth-svc

COPY --from=0 /app/pkg/config/envs/dev.env /envs/dev.env

EXPOSE 50051

CMD [ "/auth-svc" ]