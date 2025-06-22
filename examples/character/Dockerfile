# 多階段構建，優化最終鏡像大小
FROM golang:1.21-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 安裝構建依賴
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用程序
RUN CGO_ENABLED=1 GOOS=linux CGO_CFLAGS="-D_GNU_SOURCE" go build -a -installsuffix cgo -tags "sqlite_omit_load_extension" -o ff14-server .

# 運行階段：使用最小的 Alpine 鏡像
FROM alpine:latest

# 安裝運行時依賴
RUN apk --no-cache add ca-certificates sqlite

# 創建非 root 用戶
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 設定工作目錄
WORKDIR /app

# 從構建階段複製編譯好的程序
COPY --from=builder /app/ff14-server .

# 創建數據庫目錄並設定權限
RUN mkdir -p /app/data && \
    chown -R appuser:appuser /app

# 切換到非 root 用戶
USER appuser

# 暴露端口
EXPOSE 8080

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 啟動應用程序
CMD ["./ff14-server", "-port", "8080", "-db", "/app/data/characters.db"]