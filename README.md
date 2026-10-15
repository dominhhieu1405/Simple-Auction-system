# Simple-Auction-system

REST API cho hệ thống đấu giá đơn giản, viết bằng Go + SQLite.

## Nghiệp vụ

**Tài khoản**
- Người dùng đăng ký bằng username (3–32 ký tự, không trùng) và password (6–72 kí tự).
- Đăng nhập nhận JWT, dùng cho các API cần xác thực. Token có hiệu lực 24 giờ.

**Phiên đấu giá**
- Người dùng đã đăng nhập có thể tạo phiên với: tiêu đề, mô tả, giá khởi điểm, bước giá tối thiểu, thời gian kết thúc.
- Giá khởi điểm và bước giá phải > 0, thời gian kết thúc phải ở tương lai.
- Phiên mới tạo ở trạng thái `active`.
- Chỉ chủ phiên được **kết thúc** phiên. Khi kết thúc, người có giá cao nhất là người thắng.
- Chỉ chủ phiên được **xóa** phiên, và chỉ khi phiên chưa có lượt đặt giá nào.

**Đặt giá**
- Lượt đặt đầu tiên phải ≥ giá khởi điểm.
- Các lượt sau phải ≥ giá cao nhất hiện tại + bước giá.
- Không đặt giá được khi phiên đã kết thúc hoặc đã quá thời gian kết thúc.
- Người dùng xem được lịch sử đặt giá của một phiên và của chính mình.

## API

Base URL: `http://localhost:8080/api`

Swagger: `http://localhost:8080/docs/`

API cần JWT gửi header: `Authorization: Bearer <token>`

| Method | Path                    | Mô tả                                | Xác thực        |
|--------|-------------------------|--------------------------------------|-----------------|
| POST   | `/auth/register`        | Đăng ký                              | –               |
| POST   | `/auth/login`           | Đăng nhập, trả JWT                   | –               |
| GET    | `/users/me`             | Thông tin tài khoản hiện tại         | JWT             |
| GET    | `/users/me/bids`        | Lịch sử đặt giá của tôi              | JWT             |
| GET    | `/auctions`             | Danh sách phiên đấu giá              | –               |
| GET    | `/auctions/{id}`        | Chi tiết phiên                       | –               |
| POST   | `/auctions`             | Tạo phiên                            | JWT             |
| DELETE | `/auctions/{id}`        | Xóa phiên (chưa có lượt đặt giá)     | JWT (chủ phiên) |
| POST   | `/auctions/{id}/close`  | Kết thúc phiên, xác định người thắng | JWT (chủ phiên) |
| GET    | `/auctions/{id}/bids`   | Lịch sử đặt giá của phiên            | –               |
| POST   | `/auctions/{id}/bids`   | Đặt giá                              | JWT             |


## Cách chạy

Biến môi trường

| Biến          | Mặc định          |
|---------------|-------------------|
| `JWT_SECRET`  | (bắt buộc đặt)    |
| `DB_PATH`     | `data/auction.db` |
| `SERVER_PORT` | `8080`            |

### Docker

```
docker run -d -p 8080:8080 -e JWT_SECRET="change-this-secret" -v auction_data:/app/data ghcr.io/dominhhieu1405/simple-auction-system:latest
```

Hoặc dùng Docker Compose (cần file `.env` chứa `JWT_SECRET`):

```
docker compose up -d
```

### Từ source (cần Go 1.27+)

```
git clone https://github.com/dominhhieu1405/Simple-Auction-system.git
cd Simple-Auction-system
go run .
```

Tạo file `.env` từ `.env.example` và điền `JWT_SECRET` trước khi chạy.

### Test

Khi server đang chạy:

```
go test ./tests -v
```