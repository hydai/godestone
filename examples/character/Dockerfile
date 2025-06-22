# 多階段構建，優化最終鏡像大小
FROM golang:1.21 AS builder

# 設定工作目錄
WORKDIR /app

# 安裝構建依賴
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用程序
RUN CGO_ENABLED=1 go build -o ff14-server .

# 運行階段：使用最小的 Ubuntu 鏡像
FROM ubuntu:22.04

# 安裝運行時依賴
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libsqlite3-0 \
    wget \
    && rm -rf /var/lib/apt/lists/*

# 創建非 root 用戶
RUN groupadd -g 1000 appuser && \
    useradd -d /app -u 1000 -g appuser appuser

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