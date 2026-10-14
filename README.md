# Simple-Auction-system

REST API cho hệ thống đấu giá đơn giản: đăng ký/đăng nhập bằng JWT, tạo phiên đấu giá, đặt giá, kết thúc phiên và xác định người thắng.

## Công nghệ

| Thành phần    | Sử dụng                                                   |
|---------------|-----------------------------------------------------------|
| Ngôn ngữ      | Go 1.27                                                   |
| HTTP router   | `net/http` ServeMux (routing theo method, Go 1.22+)       |
| ORM           | GORM                                                      |
| Database      | SQLite (`glebarez/sqlite` – pure Go, không cần CGO)       |
| Xác thực      | JWT HS256 (`golang-jwt/jwt/v5`), mật khẩu băm bằng bcrypt |
| Tài liệu API  | Swagger (`swaggo`)                                        |
| Triển khai    | Binary đa nền tảng, Docker image (GHCR), GitHub Actions   |

## Chạy dự án

Có thể chạy bằng binary (không cần cài Go) hoặc Docker image (không cần build). Sau khi chạy:

- API: `http://localhost:8080/api`
- Swagger: `http://localhost:8080/docs/`

### Biến môi trường

| Biến          | Mặc định          | Mô tả                                                |
|---------------|-------------------|------------------------------------------------------|
| `JWT_SECRET`  | (nên đặt)         | Khóa ký JWT. Nếu thiếu, server dùng khóa mặc định và in cảnh báo |
| `DB_PATH`     | `data/auction.db` | Đường dẫn file SQLite (thư mục tự được tạo)          |
| `SERVER_PORT` | `8080`            | Cổng HTTP server                                     |

Binary tự đọc file `.env` trong thư mục đang chạy. Mẫu: `.env.example`.

### Cách 1: Chạy 1 lệnh từ URL (khuyên dùng)

Lệnh sẽ tải binary đúng hệ điều hành từ bản Release mới nhất vào thư mục `simple-auction/`, tạo `.env` với `JWT_SECRET` ngẫu nhiên (nếu chưa có) rồi chạy server.

Linux / macOS:

```bash
curl -fsSL https://raw.githubusercontent.com/dominhhieu1405/Simple-Auction-system/main/scripts/run.sh | sh
```

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/dominhhieu1405/Simple-Auction-system/main/scripts/run.ps1 | iex
```

- Chạy lại lệnh để cập nhật bản mới, dữ liệu và `.env` được giữ nguyên.
- Chọn phiên bản: `VERSION=v1.0.0` (Linux/macOS: `curl ... | VERSION=v1.0.0 sh`, Windows: `$env:VERSION="v1.0.0"` trước khi chạy).
- Đổi thư mục cài: biến `INSTALL_DIR`.

### Cách 2: Docker 1 lệnh

```bash
docker run -d --name simple-auction-api \
  -p 8080:8080 \
  -e JWT_SECRET="change-this-secret" \
  -v auction_data:/app/data \
  ghcr.io/dominhhieu1405/simple-auction-system:latest
```

Image hỗ trợ `linux/amd64` và `linux/arm64`.

### Cách 3: Docker Compose

```bash
curl -fsSLO https://raw.githubusercontent.com/dominhhieu1405/Simple-Auction-system/main/docker-compose.yml
echo "JWT_SECRET=change-this-secret" > .env
docker compose up -d
```

| Thao tác             | Lệnh                                     |
|----------------------|------------------------------------------|
| Cập nhật bản mới     | `docker compose pull && docker compose up -d` |
| Dùng phiên bản cụ thể | Thêm `APP_VERSION=1.0.0` vào `.env`     |
| Đổi cổng ngoài       | Thêm `SERVER_PORT=9000` vào `.env`       |
| Xem log              | `docker compose logs -f api`             |
| Dừng                 | `docker compose down`                    |

Compose dùng image có sẵn trên GHCR, không build. Dữ liệu SQLite lưu trong volume `auction_data`.

### Cách 4: Tải binary thủ công

Tải từ trang [Releases](https://github.com/dominhhieu1405/Simple-Auction-system/releases/latest):

| Hệ điều hành          | File                                  |
|-----------------------|---------------------------------------|
| Windows x64           | `auction-server-windows-amd64.exe`    |
| Windows ARM           | `auction-server-windows-arm64.exe`    |
| Linux x64             | `auction-server-linux-amd64`          |
| Linux ARM (Pi, VPS ARM) | `auction-server-linux-arm64`        |
| macOS Intel           | `auction-server-darwin-amd64`         |
| macOS Apple Silicon   | `auction-server-darwin-arm64`         |

`checksums.txt` chứa SHA-256 để kiểm tra file.

Linux / macOS:

```bash
chmod +x auction-server-linux-amd64
echo "JWT_SECRET=change-this-secret" > .env
./auction-server-linux-amd64
```

Windows: đặt file `.env` cạnh file `.exe` rồi chạy `.\auction-server-windows-amd64.exe`.

### Cách 5: Chạy từ source (phát triển)

Yêu cầu Go 1.27+.

```bash
git clone https://github.com/dominhhieu1405/Simple-Auction-system.git
cd Simple-Auction-system
cp .env.example .env
go run .
```

Sinh lại tài liệu Swagger sau khi sửa annotation:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### Build Docker image ở máy (phát triển)

`Dockerfile` không compile mà chỉ đóng gói binary trong `dist/`, nên cần build binary Linux trước.

Linux / macOS:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/auction-server-linux-amd64 .
docker build -t simple-auction .
```

Windows (PowerShell):

```powershell
$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"
go build -trimpath -ldflags="-s -w" -o dist/auction-server-linux-amd64 .
Remove-Item Env:GOOS, Env:GOARCH
docker build -t simple-auction .
```

Máy ARM (Apple Silicon) thay `amd64` bằng `arm64`.

## Kiến trúc

Dự án chia theo tầng, mỗi tầng chỉ gọi tầng ngay dưới nó:

```
Request
  → middleware   (xác thực JWT, gắn user ID vào context)
  → handlers     (đọc request, validate định dạng, trả response JSON)
  → services     (logic nghiệp vụ, kiểm tra quy tắc)
  → repositories (truy vấn database qua GORM)
  → SQLite
```

- `models` định nghĩa struct dùng chung cho tất cả các tầng.
- `config` đọc biến môi trường.
- Các dependency được khởi tạo và truyền thủ công trong `main.go` (repository → service → handler).
- Mỗi tầng giao tiếp qua interface. Lỗi nghiệp vụ được định nghĩa ở `services`, handler map sang HTTP status.
- Database và các bảng được tạo tự động khi khởi động (GORM AutoMigrate).
- Swagger được nhúng vào binary qua package `docs`.

## Xác thực

Đăng nhập trả về JWT, gửi kèm ở các endpoint cần xác thực:

```
Authorization: Bearer <token>
```

Token có hiệu lực **24 giờ**.

## REST API Endpoints

Base URL: `/api`

| Method | Endpoint               | Mô tả                                        | Xác thực        |
|--------|------------------------|----------------------------------------------|-----------------|
| POST   | `/auth/register`       | Đăng ký tài khoản                            | Không           |
| POST   | `/auth/login`          | Đăng nhập và nhận JWT                        | Không           |
| GET    | `/auctions`            | Lấy danh sách phiên đấu giá                  | Không           |
| GET    | `/auctions/{id}`       | Xem chi tiết phiên đấu giá                   | Không           |
| POST   | `/auctions`            | Tạo phiên đấu giá mới                        | JWT             |
| DELETE | `/auctions/{id}`       | Xóa phiên đấu giá chưa có lượt đặt giá       | JWT (chủ phiên) |
| POST   | `/auctions/{id}/bids`  | Đặt giá cho phiên đấu giá                    | JWT             |
| GET    | `/auctions/{id}/bids`  | Xem lịch sử đặt giá                          | Không           |
| POST   | `/auctions/{id}/close` | Kết thúc phiên đấu giá, xác định người thắng | JWT (chủ phiên) |
| GET    | `/users/me`            | Xem thông tin tài khoản hiện tại             | JWT             |
| GET    | `/users/me/bids`       | Xem lịch sử tham gia đấu giá của tài khoản   | JWT             |

### Định dạng lỗi

Mọi lỗi đều trả về:

```json
{ "error": "mô tả lỗi" }
```

Request body không đúng JSON, có field lạ hoặc lớn hơn 1 MB sẽ trả `400 invalid request payload`.

### Chi tiết endpoint

#### POST `/auth/register`

```json
{ "username": "alice", "password": "password123" }
```

| Status | Ý nghĩa                                          |
|--------|--------------------------------------------------|
| 201    | `{ "message": "user registered successfully" }`  |
| 400    | Username/password không hợp lệ                   |
| 409    | Username đã tồn tại                              |

#### POST `/auth/login`

```json
{ "username": "alice", "password": "password123" }
```

| Status | Ý nghĩa                    |
|--------|----------------------------|
| 200    | `{ "token": "<jwt>" }`     |
| 401    | Sai username hoặc password |

#### GET `/users/me`

| Status | Ý nghĩa                            |
|--------|------------------------------------|
| 200    | `{ "id": 1, "username": "alice" }` |
| 401    | Thiếu hoặc sai token               |
| 404    | User không tồn tại                 |

#### POST `/auctions`

```json
{
  "title": "Đồng hồ cổ",
  "description": "Mô tả sản phẩm",
  "starting_price": 1000000,
  "min_increment": 50000,
  "end_time": "2026-12-31T23:59:59+07:00"
}
```

`end_time` theo định dạng RFC 3339.

| Status | Ý nghĩa                             |
|--------|-------------------------------------|
| 201    | Trả về object phiên đấu giá vừa tạo |
| 400    | Dữ liệu không hợp lệ                |
| 401    | Chưa đăng nhập                      |

Object phiên đấu giá:

```json
{
  "id": 1,
  "owner_id": 1,
  "title": "Đồng hồ cổ",
  "description": "Mô tả sản phẩm",
  "starting_price": 1000000,
  "min_increment": 50000,
  "end_time": "2026-12-31T23:59:59+07:00",
  "status": "active",
  "created_at": "...",
  "updated_at": "..."
}
```

#### GET `/auctions`

`200` – mảng phiên đấu giá, mới nhất trước.

#### GET `/auctions/{id}`

| Status | Ý nghĩa              |
|--------|----------------------|
| 200    | Object phiên đấu giá |
| 400    | ID không hợp lệ      |
| 404    | Không tìm thấy       |

#### DELETE `/auctions/{id}`

| Status | Ý nghĩa                                         |
|--------|-------------------------------------------------|
| 200    | `{ "message": "auction deleted successfully" }` |
| 400    | ID không hợp lệ                                 |
| 401    | Chưa đăng nhập                                  |
| 403    | Không phải chủ phiên                            |
| 404    | Không tìm thấy                                  |
| 409    | Phiên đã có lượt đặt giá                        |

#### POST `/auctions/{id}/bids`

```json
{ "amount": 1050000 }
```

| Status | Ý nghĩa                                      |
|--------|----------------------------------------------|
| 201    | Trả về object bid vừa tạo                    |
| 400    | Số tiền ≤ 0 hoặc thấp hơn mức tối thiểu      |
| 401    | Chưa đăng nhập                               |
| 404    | Không tìm thấy phiên                         |
| 409    | Phiên đã đóng hoặc đã quá thời gian kết thúc |

Object bid:

```json
{
  "id": 1,
  "auction_id": 1,
  "user_id": 2,
  "amount": 1050000,
  "created_at": "..."
}
```

#### GET `/auctions/{id}/bids`

| Status | Ý nghĩa                  |
|--------|--------------------------|
| 200    | Mảng bid, mới nhất trước |
| 400    | ID không hợp lệ          |
| 404    | Không tìm thấy phiên     |

#### POST `/auctions/{id}/close`

| Status | Ý nghĩa                                      |
|--------|----------------------------------------------|
| 200    | `{ "auction": {...}, "winning_bid": {...} }` |
| 400    | ID không hợp lệ                              |
| 401    | Chưa đăng nhập                               |
| 403    | Không phải chủ phiên                         |
| 404    | Không tìm thấy                               |
| 409    | Phiên đã đóng trước đó                       |

`winning_bid` là `null` nếu phiên không có lượt đặt giá nào.

#### GET `/users/me/bids`

`200` – mảng bid của người dùng hiện tại, mới nhất trước.

## Quy tắc nghiệp vụ

### Tài khoản

- Username dài 3–32 ký tự, không trùng, khoảng trắng đầu/cuối được loại bỏ.
- Password dài 6–72 byte, lưu dưới dạng bcrypt hash, không bao giờ trả về trong response.

### Phiên đấu giá

- `title` không được rỗng.
- `starting_price > 0`, `min_increment > 0`.
- `end_time` phải lớn hơn thời điểm hiện tại.
- Phiên mới tạo có trạng thái `active`.
- Trạng thái: `active` (đang nhận giá), `closed` (đã kết thúc).
- Tiền lưu dạng số nguyên (`int64`), không dùng số thực.

### Đặt giá

- `amount > 0`.
- Bid đầu tiên phải ≥ `starting_price`.
- Các bid sau phải ≥ giá cao nhất hiện tại + `min_increment`.
- Không thể đặt giá khi phiên đã `closed` hoặc đã quá `end_time`.

### Kết thúc phiên

- Chỉ chủ phiên được kết thúc.
- Người thắng là người có lượt đặt giá cao nhất.
- Phiên chỉ kết thúc khi chủ phiên gọi API `/close`.

### Xóa phiên

- Chỉ chủ phiên được xóa.
- Chỉ xóa được khi phiên chưa có lượt đặt giá nào.

## Database

| Bảng     | Các cột                                                                                                   | Mô tả             |
|----------|-----------------------------------------------------------------------------------------------------------|-------------------|
| users    | id, username, password                                                                                    | Người dùng        |
| auctions | id, owner_id, title, description, starting_price, min_increment, end_time, status, created_at, updated_at | Các phiên đấu giá |
| bids     | id, auction_id, user_id, amount, created_at                                                               | Lịch sử trả giá   |

Ràng buộc:

- `users.username`: unique, not null.
- `auctions.owner_id` → `users.id` (có index).
- `auctions.end_time`: có index.
- `bids.auction_id` → `auctions.id` (có index).
- `bids.user_id` → `users.id` (có index).

## Kiểm thử

Test là integration test gọi trực tiếp API, **server phải đang chạy** ở `localhost:8080`:

```bash
go run .            # terminal 1 (hoặc docker compose up -d)
go test ./tests -v  # terminal 2
```

| File                         | Nội dung                                                                 |
|------------------------------|--------------------------------------------------------------------------|
| `tests/auth_api_test.go`     | Đăng ký, đăng nhập, lấy thông tin user                                   |
| `tests/auction_api_test.go`  | Tạo phiên, xem danh sách/chi tiết, đặt giá, lịch sử, kết thúc phiên, chặn bid sau khi đóng, chặn xóa phiên có bid |

## CI/CD

Một workflow duy nhất: `.github/workflows/cicd.yml`.

| Job             | Chạy khi                         | Nội dung                                                                 |
|-----------------|----------------------------------|--------------------------------------------------------------------------|
| `build`         | Mọi lần push, PR, chạy tay       | `go vet`, build binary cho 6 nền tảng vào `dist/` kèm `checksums.txt`    |
| `test`          | Sau `build`                      | Đóng gói binary thành Docker image, chạy bằng Docker Compose, chạy API test |
| `publish-image` | Push lên `main`/`master` hoặc tag `v*`, test pass | Đẩy image `linux/amd64` + `linux/arm64` lên GHCR                         |
| `release`       | Push tag `v*`, test pass         | Tạo GitHub Release, đính kèm binary, `checksums.txt`, `docker-compose.yml` |

Go chỉ compile một lần ở job `build`. Các job sau dùng lại binary đó.

Tag image:

| Sự kiện               | Tag                         |
|-----------------------|-----------------------------|
| Push `main`           | `latest`                    |
| Push tag `v1.2.3`     | `1.2.3`, `1.2`              |

### Phát hành phiên bản mới

```bash
git tag v1.0.0
git push origin v1.0.0
```

### Lưu ý GHCR

Package trên GHCR mặc định là **private** sau lần đẩy đầu tiên. Để mọi người `docker pull` được mà không cần đăng nhập: GitHub → Packages → `simple-auction-system` → Package settings → Change visibility → Public.

## Cấu trúc thư mục

```
Simple-Auction-System/
├── main.go                     # Khởi chạy ứng dụng, DI, đăng ký routes
│
├── internal/
│   ├── config/
│   │   └── config.go           # Đọc biến môi trường
│   │
│   ├── handlers/               # Tầng API
│   │   ├── auth.go             # Đăng ký, đăng nhập, thông tin cá nhân
│   │   ├── auction.go          # Tạo, xem, xóa, kết thúc phiên đấu giá
│   │   ├── bid.go              # Đặt giá, lịch sử đặt giá
│   │   └── response.go         # Hàm ghi JSON, đọc request, kiểu response chung
│   │
│   ├── services/               # Logic nghiệp vụ
│   │   ├── auth.go             # Xác thực, JWT, xử lý thông tin người dùng
│   │   ├── auction.go          # Nghiệp vụ quản lý phiên đấu giá
│   │   └── bid.go              # Nghiệp vụ đặt giá, kiểm tra giá hợp lệ
│   │
│   ├── repositories/           # Truy cập dữ liệu
│   │   ├── user.go             # Truy vấn dữ liệu người dùng
│   │   ├── auction.go          # Truy vấn dữ liệu phiên đấu giá
│   │   └── bid.go              # Truy vấn và lưu lượt đặt giá
│   │
│   ├── models/                 # Định nghĩa cấu trúc dữ liệu
│   │   ├── user.go             # Struct User
│   │   ├── auction.go          # Struct Auction
│   │   └── bid.go              # Struct Bid
│   │
│   └── middleware/
│       └── jwt.go              # Middleware xác thực JWT
│
├── tests/
│   ├── auth_api_test.go        # Test API xác thực
│   └── auction_api_test.go     # Test luồng đấu giá
│
├── scripts/
│   ├── run.sh                  # Tải binary và chạy (Linux/macOS)
│   └── run.ps1                 # Tải binary và chạy (Windows)
│
├── docs/                       # Swagger/OpenAPI sinh tự động bằng `swag init`
├── data/                       # File SQLite (không commit)
├── dist/                       # Binary đã build (không commit)
│
├── .github/workflows/
│   └── cicd.yml                # Build, test, đẩy image, tạo release
│
├── .env.example                # Mẫu biến môi trường
├── Dockerfile                  # Đóng gói binary trong dist/ thành image
├── docker-compose.yml          # Chạy image từ GHCR, lưu SQLite bằng volume
├── go.mod / go.sum             # Go module và dependencies
└── README.md                   # Tài liệu
```