# ---------- Этап 1: сборка ----------
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Сначала только файлы зависимостей: слой кэшируется,
# пока go.mod/go.sum не менялись
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Затем исходники и сборка статического бинарника
COPY backend/ ./

# CGO_ENABLED=0 -> статический бинарник (в distroless/static нет glibc)
# -ldflags="-s -w" убирает отладочную информацию
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/app ./main.go

# ---------- Этап 2: исполнение ----------
# distroless/static: нет shell и пакетного менеджера; тег :nonroot — не root
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Из первого этапа берём только бинарник, компилятор в образ не попадает
COPY --from=builder /out/app /app/app

# Веб-приложение (frontend) упаковывается прямо в образ
COPY frontend/ /app/static/

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/app"]
