# Simple-Auction-system

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

## Database

| Bảng     | Các cột                                                                           | Mô tả           |
|----------|-----------------------------------------------------------------------------------|-----------------|
| users    | id, username, password                                                            | Người dùng      |
| auctions | id, owner_id, title, description, starting_price, min_increment, end_time, status | Tên đăng nhập   |
| bids     | id, auction_id, user_id, amount, create_at                                        | Lịch sử trả giá |

## Cấu trúc thư mục
```
Simple-Auction-System/
├── main.go                     # Khởi chạy ứng dụng
│
├── internal/                   # Viết code trong đây
│   │
│   ├── handlers/               # Tầng API
│   │   ├── auth.go             # Đăng ký, đăng nhập, thông tin cá nhân
│   │   ├── auction.go          # Tạo, xem, xóa, kết thúc phiên đấu giá
│   │   └── bid.go              # Đặt giá, lịch sử đặt giá
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
│   └── middleware/             # Xử lý request trước khi chuyển đến handler
│       └── jwt.go              # Middleware xác thực JWT
│
├── data/                       # Lưu dữ liệu SQLite
├── docs/                       # Tài liệu Swagger/OpenAPI được sinh tự động
├── tests/                      # Chứa script kiểm thử hệ thống
├── go.mod                      # Khai báo Go module và các dependencies
├── Dockerfile                  # Hướng dẫn build Docker image cho ứng dụng
├── docker-compose.yml          # Cấu hình chạy container và lưu trữ SQLite bằng volume
└── README.md                   # Tài liệu
 ```
