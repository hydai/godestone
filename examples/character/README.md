# FF14 角色資料 HTTP Server

這是一個 HTTP Server，可以透過 API 獲取 Final Fantasy XIV 的角色資料，並使用 SQLite 作為快取來避免過度請求 Lodestone。

## 功能特色

- **HTTP API**: 提供 RESTful API 來獲取角色資料
- **SQLite 快取**: 自動快取角色資料 24 小時，避免頻繁請求
- **Rate Limit 保護**: 透過快取機制減少對 Lodestone 的請求
- **優雅關閉**: 支援 SIGTERM 和 SIGINT 信號處理

## 安裝與使用

### 1. 安裝依賴
```bash
go mod tidy
```

### 2. 建置程式
```bash
go build -o ff14-server .
```

### 3. 啟動伺服器
```bash
# 使用預設設定 (port: 8080, db: characters.db)
./ff14-server

# 自訂端口和資料庫路徑
./ff14-server -port 9000 -db /path/to/characters.db
```

## API 使用方法

### 獲取角色資料
```bash
GET /api/character/{id}
```

範例：
```bash
curl http://localhost:8080/api/character/46441202
```

回應格式：
```json
{
  "data": {
    "ActiveClassJob": {...},
    "Avatar": "https://...",
    "Bio": "...",
    "ClassJobs": [...],
    ...
  },
  "cached": false,
  "timestamp": "2025-06-22T22:10:33Z"
}
```

### 健康檢查
```bash
GET /health
```

### 根路徑說明
```bash
GET /
```

## 資料庫結構

SQLite 資料庫包含以下欄位：
- `id`: 角色 ID (PRIMARY KEY)
- `character_data`: JSON 格式的角色資料
- `cached_at`: 快取時間
- `updated_at`: 更新時間

快取有效期為 24 小時。

## 命令行參數

- `-port`: 伺服器端口 (預設: 8080)
- `-db`: SQLite 資料庫檔案路徑 (預設: characters.db)

## 錯誤處理

- **400 Bad Request**: 無效的角色 ID
- **404 Not Found**: 角色不存在或無法抓取
- **500 Internal Server Error**: 伺服器內部錯誤

## 注意事項

1. 第一次請求角色資料會從 Lodestone 抓取，可能需要較長時間
2. 後續 24 小時內的請求會直接從快取返回
3. 請適度使用，避免對 Lodestone 造成過大負載
4. 資料庫檔案會在第一次啟動時自動建立