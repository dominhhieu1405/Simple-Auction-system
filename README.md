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
# 📊 Báo Cáo Kiểm Thử Hiệu Năng Hệ Thống (Load Testing Report)

Dưới đây là kết quả kiểm thử tải trên hệ thống Backend (Go + GORM + SQLite) cho 11 endpoint nghiệp vụ:

###  Bảng Thống Kê Hiệu Năng & Tài Nguyên

| Users | 201 (Thành công) | 400 / 401 | 500 (Lỗi Server) | Tỷ lệ 500 (%) | p50 (ms) | p95 (ms) | CPU Trung bình | CPU Max | RAM (MB) |
| :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| **10** | ~580 | ~490 | 0 | **0.00%** | ~3.1ms | ~58.9ms | 0.3% | 18.5% | 44.4 MB |
| **100** | ~1,568 | ~2,600 | 0 | **0.00%** | ~131.4ms | ~565.7ms | 0.2% | 5.5% | 44.9 MB |
| **1,000** | ~1,079 | ~4,240 | 260 | **4.33%** | ~2,197.9ms | ~4,908.2ms | 0.0% | 0.0% | 45.3 MB |

---

###  Nhận Xét & Đánh Giá

* **Mức <= 100 users:** Hệ thống hoạt động hoàn toàn ổn định, 0 lỗi, độ trễ `p95` duy trì ở mức tốt (< 600 ms). Tiến trình Go cho thấy khả năng tối ưu bộ nhớ tuyệt vời khi RAM gần như không đổi (~44.9 MB).
* **Mức 1,000 users:** Hệ thống bắt đầu quá tải. Độ trễ `p95` vọt lên ~4.9 giây. Bắt đầu xuất hiện lỗi với tỷ lệ 4.33% (toàn bộ là mã lỗi 500 ở các request tương tác dữ liệu).
* **Nút thắt hệ thống (Bottleneck):** Nằm ở các thao tác ghi cơ sở dữ liệu (I/O Bound) do cơ chế khóa toàn cục của file SQLite, hoàn toàn không phải do thiếu hụt sức mạnh tính toán. CPU thực tế rất thấp (đạt đỉnh 18.5% do khởi động lạnh ở mức 10 users, sau đó duy trì quanh mức 0% đến 0.2%).
* **Ngưỡng chịu tải (Capacity Limit):** Hệ thống nguyên bản đáp ứng mượt mà ở khoảng 100 users đồng thời và bắt đầu bộc lộ giới hạn khi tiệm cận ngưỡng 1,000 users.

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

