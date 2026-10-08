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


## Database Schema


| Bảng     | Các cột                                                                           | Mô tả           |
|----------|-----------------------------------------------------------------------------------|-----------------|
| users    | id, username, password                                                            | Người dùng      |
| auctions | id, owner_id, title, description, starting_price, min_increment, end_time, status | Tên đăng nhập   |
| bids     | id, auction_id, user_id, amount, create_at                                        | Lịch sử trả giá |
