# Simple-Auction-System

Dự án cung cấp hệ thống cho nền tảng đấu giá trực tuyến đơn giản. Hệ thống được viết bằng Go (Golang 1.27) kết hợp với cơ sở dữ liệu SQLite và được thiết kế theo chuẩn kiến trúc phân tầng (Layered Architecture).

---

## Logic Nghiệp vụ

**1. Tài khoản (Authentication)**
- Người dùng đăng ký bằng username (3-32 ký tự, duy nhất) và password (6-72 kí tự).
- Đăng nhập thành công sẽ nhận được thẻ JWT. Token này dùng cho các API cần xác thực và có hiệu lực trong 24 giờ.

**2. Phiên đấu giá (Auctions)**
- Người dùng đã đăng nhập có thể tạo phiên đấu giá với các thông tin: tiêu đề, mô tả, giá khởi điểm, bước giá tối thiểu và thời gian kết thúc.
- **Điều kiện:** Giá khởi điểm và bước giá phải `> 0`, thời gian kết thúc phải lớn hơn thời điểm hiện tại.
- Phiên mới tạo mặc định ở trạng thái `active`.
- **Phân quyền Chủ phiên:**
  - Chỉ chủ phiên mới được **Kết thúc (Close)** phiên. Khi kết thúc, người có lượt đặt giá cao nhất sẽ được vinh danh là người chiến thắng.
  - Chỉ chủ phiên mới được **Xóa (Delete)** phiên, với điều kiện phiên đó *chưa có bất kỳ lượt đặt giá nào*.

**3. Đặt giá (Bidding)**
- Lượt đặt giá đầu tiên phải `>= Giá khởi điểm`.
- Các lượt đặt giá tiếp theo phải `>= Giá cao nhất hiện tại + Bước giá tối thiểu`.
- Hệ thống từ chối các lượt đặt giá nếu phiên đã kết thúc (`closed`) hoặc đã vượt quá thời gian kết thúc (`end_time`).

---

## Công nghệ sử dụng

- **Ngôn ngữ lập trình:** Go (Golang) version 1.27+
- **Cơ sở dữ liệu:** SQLite (sử dụng driver thuần Go `glebarez/sqlite` không phụ thuộc CGO)
- **ORM:** GORM
- **Xác thực & Bảo mật:** JSON Web Token (`golang-jwt/v5`) và thuật toán băm mật khẩu `bcrypt`
- **Routing:** `net/http` ServeMux (Thư viện chuẩn của Go)
- **Triển khai (Deployment):** Docker, Docker Compose
- **CI/CD:** GitHub Actions
  
## Cấu trúc thư mục

Dự án được tổ chức tuân thủ cấu trúc thư mục tiêu chuẩn của Go:

```text
Simple-Auction-System/
├── data/                   # Chưa CSDL auction.db
├── internal/               # Code
│   ├── handlers/           # (Tầng Controllers) Nhận HTTP Request, Parse JSON, Gọi Services
│   ├── middleware/         # (Tầng Middleware) Chặn request, kiểm tra JWT, nhúng userID
│   ├── models/             # (Tầng Data) Định nghĩa các struct và GORM mapping
│   ├── repositories/       # (Tầng Database) Thực thi các câu lệnh truy vấn CSDL
│   └── services/           # (Tầng Nghiệp vụ) Xử lý logic, tính toán, kiểm tra quyền
├── .env.example            # File mẫu chứa các biến môi trường
├── main.go                 # Điểm bắt đầu (Entry point), khởi tạo hệ thống
└── README.md               # Tài liệu dự án
```

**Luồng dữ liệu :**
`Client` -> `Router (main.go)` -> `Middleware` -> `Handlers` -> `Services` -> `Repositories` -> `Database`

---

## Đặc tả API (API Specification)

- **Base URL:** `http://localhost:8080/api`
- **Swagger UI:** `http://localhost:8080/docs/`
- **Yêu cầu xác thực:** Các API đánh dấu "JWT" yêu cầu truyền header `Authorization: Bearer <token>`

| Method | Path | Mô tả chức năng | Xác thực |
| :--- | :--- | :--- | :--- |
| **POST** | `/auth/register` | Đăng ký tài khoản mới | - |
| **POST** | `/auth/login` | Đăng nhập, nhận chuỗi JWT | - |
| **GET** | `/users/me` | Lấy thông tin tài khoản đang đăng nhập | JWT |
| **GET** | `/users/me/bids` | Xem lịch sử các lượt đặt giá của bản thân | JWT |
| **GET** | `/auctions` | Xem danh sách tất cả phiên đấu giá | - |
| **GET** | `/auctions/{id}` | Xem thông tin chi tiết một phiên đấu giá | - |
| **POST** | `/auctions` | Tạo phiên đấu giá sản phẩm mới | JWT |
| **DELETE** | `/auctions/{id}` | Xóa phiên (chỉ khi chưa có ai trả giá) | JWT (Chủ phiên) |
| **POST** | `/auctions/{id}/close` | Đóng phiên, chốt kết quả và người thắng | JWT (Chủ phiên) |
| **GET** | `/auctions/{id}/bids` | Lấy lịch sử các lượt đặt giá của một phiên | - |
| **POST** | `/auctions/{id}/bids` | Đặt giá mới cho một phiên | JWT |

---

## Quy ước Kỹ thuật 

**1. Định dạng dữ liệu (Data Types)**
- **Tiền tệ:** Mọi trường liên quan đến tiền (`starting_price`, `min_increment`, `amount`) đều sử dụng số nguyên `int64` (không dùng số thực) để đảm bảo độ chính xác tuyệt đối.
- **Thời gian:** Các trường thời gian (`end_time`, `created_at`, `updated_at`) luôn trả về và yêu cầu định dạng chuẩn **RFC 3339 / ISO 8601** (Ví dụ: `2026-12-31T23:59:59+07:00`).

**2. Giới hạn Payload (Security)**
- Các request `POST` / `PUT` chỉ chấp nhận Body JSON có dung lượng tối đa **1MB**.
- Hệ thống áp dụng **Strict JSON**: Chặn mọi request chứa các trường dữ liệu lạ không nằm trong đặc tả.

**3. Bảng Mã Lỗi & Định dạng Lỗi (Error Handling)**
Khi hệ thống gặp lỗi, API sẽ luôn trả về một **HTTP Status Code** tương ứng kèm theo body JSON thống nhất có cấu trúc:
```json
{
  "error": "mô tả chi tiết nguyên nhân lỗi"
}
```

**Bảng giải nghĩa HTTP Status Code:**

| HTTP Code | Ý nghĩa | Tình huống phát sinh |
| :---: | :--- | :--- |
| **400** | `Bad Request` | Lỗi dữ liệu gửi lên: Sai định dạng JSON, gửi thừa trường, thiếu trường bắt buộc, giá `≤ 0`, gửi mã ID trên URL không phải số nguyên. |
| **401** | `Unauthorized` | Lỗi định danh: Sai `username`/`password` khi đăng nhập, thiếu token JWT, token hết hạn hoặc token sai chữ ký. |
| **403** | `Forbidden` | Lỗi phân quyền: Đã có token nhưng cố tình xóa hoặc đóng phiên đấu giá của người khác. |
| **404** | `Not Found` | Không tìm thấy tài nguyên: Phiên đấu giá (Auction ID) không tồn tại. |
| **409** | `Conflict` | Xung đột nghiệp vụ: Đăng ký trùng username, đặt giá vào phiên đã đóng/hết hạn, cố xóa phiên đã có lượt trả giá. |
| **500** | `Internal Server Error` | Lỗi hệ thống: Lỗi kết nối DB hoặc lỗi không xác định từ máy chủ. |

---

## Test
## Kết quả load test

**Cấu hình:** `127.0.0.1:8636` | ramp 5s + chạy 30s mỗi mức | nghỉ ~1s giữa 2 request | client và server chạy cùng VPS

### Tổng quan

| Users | Request | Req/s | Lỗi | Lỗi % |
|---:|---:|---:|---:|---:|
| 10 | 336 | 9.4 | 0 | 0.00 |
| 100 | 3,228 | 88.9 | 0 | 0.00 |
| 1,000 | 19,577 | 469.8 | 1,261 | 6.44 |
| 10,000 | 54,319 | 1,145.2 | 5,220 | 9.61 |

### Độ trễ (ms)

| Users | Min | Avg | p50 | p90 | p95 | p99 | Max |
|---:|---:|---:|---:|---:|---:|---:|---:|
| 10 | 1.1 | 9.6 | 2.3 | 44.4 | 55.7 | 64.2 | 87.8 |
| 100 | 0.6 | 25.7 | 2.7 | 60.1 | 97.0 | 235.6 | 737.2 |
| 1,000 | 0.9 | 735.9 | 195.0 | 1,655.6 | 5,289.4 | 6,008.8 | 9,100.8 |
| 10,000 | 0.9 | 4,998.1 | 3,717.5 | 10,665.8 | 15,055.8 | 21,363.1 | 30,181.0 |

### p95 theo endpoint (ms)

| Endpoint | 10 | 100 | 1,000 | 10,000 |
|---|---:|---:|---:|---:|
| `GET /auctions` | 21.2 | 84.9 | 985.6 | 6,402.7 |
| `GET /auctions/{id}` | 2.2 | 83.3 | 1,022.4 | 6,379.8 |
| `GET /auctions/{id}/bids` | 3.3 | 84.4 | 1,466.0 | 11,248.9 |
| `GET /users/me` | 2.2 | 108.5 | 876.3 | 6,686.5 |
| `POST /auctions/{id}/bids` | 66.8 | 142.6 | 6,348.6 | 22,492.5 |

### POST bid: mã trả về

| Users | 201 | 400 | 500 | Timeout | Tỷ lệ 500 |
|---:|---:|---:|---:|---:|---:|
| 10 | 44 | 0 | 0 | 0 | 0% |
| 100 | 338 | 1 | 0 | 0 | 0% |
| 1,000 | 659 | 83 | 1,261 | 0 | 63.0% |
| 10,000 | 143 | 52 | 5,214 | 1 | 96.4% |

> 400 = bid thấp hơn mức tối thiểu (từ chối nghiệp vụ, không tính là lỗi).
> 500 = `database is locked (5) (SQLITE_BUSY)`.

### Tài nguyên

| Users | CPU server (avg / max) | RAM server (avg / max) | CPU script max |
|---:|---:|---:|---:|
| 10 | 2% / 6% | 29 / 31 MB | 4% |
| 100 | 22% / 40% | 33 / 37 MB | 12% |
| 1,000 | 247% / 462% | 138 / 166 MB | 60% |
| 10,000 | 689% / 1,397% | 1,393 / 1,916 MB | 110% |

> CPU server: 100% = 1 nhân. CPU máy: tính trên tổng số nhân.

### Nguyên nhân lỗi

- **Lỗi 500 khi bid là do SQLite khóa ghi.** Log server báo `place bid: database is locked (5) (SQLITE_BUSY)`.
- SQLite chỉ cho **1 thao tác ghi tại một thời điểm**. Khi nhiều request bid cùng insert, các request còn lại phải chờ. Chờ quá thời gian cho phép thì nhận `SQLITE_BUSY` và trả về 500.
- Từ 100 users trở xuống, lượng ghi còn ít nên hầu như không phải chờ. Từ 1,000 users trở lên, request ghi dồn lại thành hàng đợi và lỗi tăng vọt (63% rồi 96%).

### Nhận xét

- **<= 100 users:** ổn định, 0 lỗi, p95 < 100 ms.
- **1,000 users:** quá tải. p95 = 5.3 s, lỗi 6.4%, toàn bộ là 500 ở POST bid.
- **10,000 users:** gần sập. p95 = 15 s, 96% lượt bid lỗi, RAM tăng ~60 lần (31 MB -> 1.9 GB), FD 18.5k.
- **Req/s không tăng tuyến tính:** users x10 thì req/s chỉ tăng x5.3 (100 -> 1k users) và x2.4 (1k -> 10k users).
- **Nút thắt là ghi DB, không phải CPU:** CPU máy tối đa chỉ 34%. Ngay từ 10 users, bid đã chậm nhất (p50 51 ms, so với ~2 ms của API đọc).
- **API đọc không lỗi** nhưng chậm dần khi tải ghi tăng, vì phải chờ khóa của DB.
- **Số liệu mức 10,000 users chưa chính xác:** CPU của script 110% nghĩa là client Python đã bị nghẽn
- **Ngưỡng chịu tải:** khoảng 100-1,000 users đồng thời.

## Hướng dẫn Cài đặt & Chạy (How to run)

### Cấu hình Biến môi trường (`.env`)

| Biến | Giải thích | Giá trị mặc định |
| :--- | :--- | :--- |
| `JWT_SECRET` | Chìa khóa ký JWT *(Bắt buộc phải đổi)* | (Không có) |
| `DB_PATH` | Đường dẫn lưu file cơ sở dữ liệu SQLite | `data/auction.db` |
| `SERVER_PORT` | Cổng HTTP mà server lắng nghe | `8080` |

### Cách 1: Chạy bằng Docker Compose (Khuyên dùng)
Yêu cầu đã cài đặt Docker. Khởi tạo file `.env` chứa biến `JWT_SECRET` trước khi chạy.
```bash
docker compose up -d
```
*Hệ thống sẽ chạy ngầm và tự động mount dữ liệu vào thư mục `data/` để không bị mất dữ liệu khi xóa container.*

### Cách 2: Chạy trực tiếp từ Source Code
Yêu cầu đã cài đặt Go 1.27+.
```bash
# 1. Clone dự án
git clone [https://github.com/dominhhieu1405/Simple-Auction-system.git](https://github.com/dominhhieu1405/Simple-Auction-system.git)
cd Simple-Auction-system

# 2. Tạo file môi trường
cp .env.example .env
# (Mở file .env và điền giá trị cho JWT_SECRET)

# 3. Chạy server
go run .
```

## Gợi ý cho P2

- **Hiệu năng và khả năng mở rộng dữ liệu:**
  - **Hạn chế:** Chưa có phân trang
  - **Phương hướng:** Triển khai cơ chế phân trang qua các tham số `limit` và `offset`.

- **Xử lý tranh chấp dữ liệu (Race Condition):**
  - **Hạn chế:** Các request đặt giá đồng thời cùng đọc một mức giá cũ, gây xung đột.
  - **Phương hướng:** Bọc toàn bộ logic kiểm tra và lưu giá mới vào trong Transaction của Database.

- **Bảo mật:**
  - **Hạn chế:** Chưa có rate limit cho đăng nhập và đăng ký. Nên user có thể spam tạo acc hoặc brute force
  - **Phương hướng:** Tích hợp Rate Limit, captcha.

