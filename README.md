Simple-Auction-SystemDự án cung cấp hệ thống Backend RESTful API cho nền tảng đấu giá trực tuyến đơn giản. Hệ thống được viết bằng Go (Golang 1.27) kết hợp với cơ sở dữ liệu SQLite và được thiết kế theo chuẩn kiến trúc phân tầng (Layered Architecture).Nghiệp vụ (Business Logic)1. Tài khoản (Authentication)Người dùng đăng ký bằng username (3–32 ký tự, duy nhất) và password (6–72 kí tự).Đăng nhập thành công sẽ nhận được thẻ JWT. Token này dùng cho các API cần xác thực và có hiệu lực trong 24 giờ.2. Phiên đấu giá (Auctions)Người dùng đã đăng nhập có thể tạo phiên đấu giá với các thông tin: tiêu đề, mô tả, giá khởi điểm, bước giá tối thiểu và thời gian kết thúc.Điều kiện: Giá khởi điểm và bước giá phải $> 0$, thời gian kết thúc phải lớn hơn thời điểm hiện tại.Phiên mới tạo mặc định ở trạng thái active.Phân quyền Chủ phiên:Chỉ chủ phiên mới được Kết thúc (Close) phiên. Khi kết thúc, người có lượt đặt giá cao nhất sẽ được vinh danh là người chiến thắng.Chỉ chủ phiên mới được Xóa (Delete) phiên, với điều kiện phiên đó chưa có bất kỳ lượt đặt giá nào.3. Đặt giá (Bidding)Lượt đặt giá đầu tiên phải $\ge$ Giá khởi điểm.Các lượt đặt giá tiếp theo phải $\ge$ Giá cao nhất hiện tại + Bước giá tối thiểu.Hệ thống từ chối các lượt đặt giá nếu phiên đã kết thúc (closed) hoặc đã vượt quá thời gian kết thúc (end_time).Cấu trúc thư mục (Directory Structure)Dự án được tổ chức tuân thủ cấu trúc thư mục tiêu chuẩn của Go:Simple-Auction-System/
├── .github/workflows/      # Cấu hình CI/CD GitHub Actions (Tự động build & test)
├── data/                   # Thư mục chứa file SQLite (auction.db) khi chạy thật
├── internal/               # Chứa mã nguồn nội bộ (không thể import từ project khác)
│   ├── handlers/           # (Tầng Controllers) Nhận HTTP Request, Parse JSON, Gọi Services
│   ├── middleware/         # (Tầng Gác cổng) Chặn request, kiểm tra JWT, nhúng userID
│   ├── models/             # (Tầng Data) Định nghĩa các struct và GORM mapping
│   ├── repositories/       # (Tầng Database) Thực thi các câu lệnh truy vấn CSDL
│   └── services/           # (Tầng Nghiệp vụ) "Bộ não" xử lý logic, tính toán, kiểm tra quyền
├── tests/                  # Thư mục chứa các file kiểm thử (Integration Tests)
├── .env.example            # File mẫu chứa các biến môi trường
├── docker-compose.yml      # Cấu hình triển khai hệ thống bằng Docker Compose
├── Dockerfile              # Chỉ dẫn đóng gói ứng dụng thành Docker Image
├── go.mod / go.sum         # Quản lý thư viện phụ thuộc của Go
├── main.go                 # Điểm bắt đầu (Entry point), khởi tạo hệ thống
└── README.md               # Tài liệu dự án
Luồng dữ liệu (Data Flow):
Client ➔ Router (main.go) ➔ Middleware ➔ Handlers ➔ Services ➔ Repositories ➔ DatabaseĐặc tả API (API Specification)Base URL: http://localhost:8080/apiSwagger UI: http://localhost:8080/docs/Yêu cầu xác thực: Các API đánh dấu "JWT" yêu cầu truyền header Authorization: Bearer <token>MethodPathMô tả chức năngXác thựcPOST/auth/registerĐăng ký tài khoản mới–POST/auth/loginĐăng nhập, nhận chuỗi JWT–GET/users/meLấy thông tin tài khoản đang đăng nhậpJWTGET/users/me/bidsXem lịch sử các lượt đặt giá của bản thânJWTGET/auctionsXem danh sách tất cả phiên đấu giá–GET/auctions/{id}Xem thông tin chi tiết một phiên đấu giá–POST/auctionsTạo phiên đấu giá sản phẩm mớiJWTDELETE/auctions/{id}Xóa phiên (chỉ khi chưa có ai trả giá)JWT (Chủ phiên)POST/auctions/{id}/closeĐóng phiên, chốt kết quả và người thắngJWT (Chủ phiên)GET/auctions/{id}/bidsLấy lịch sử các lượt đặt giá của một phiên–POST/auctions/{id}/bidsĐặt giá mới cho một phiênJWTQuy ước Kỹ thuật (Dành cho Frontend / QA)1. Định dạng dữ liệu (Data Types)Tiền tệ: Mọi trường liên quan đến tiền (starting_price, min_increment, amount) đều sử dụng số nguyên int64 (không dùng số thực) để đảm bảo độ chính xác tuyệt đối.Thời gian: Các trường thời gian (end_time, created_at, updated_at) luôn trả về và yêu cầu định dạng chuẩn RFC 3339 / ISO 8601 (Ví dụ: 2026-12-31T23:59:59+07:00).2. Giới hạn Payload (Security)Các request POST / PUT chỉ chấp nhận Body JSON có dung lượng tối đa 1MB.Hệ thống áp dụng Strict JSON: Chặn mọi request chứa các trường dữ liệu lạ không nằm trong đặc tả.3. Bảng Mã Lỗi & Định dạng Lỗi (Error Handling)
Khi hệ thống gặp lỗi, API sẽ luôn trả về một HTTP Status Code tương ứng kèm theo body JSON thống nhất có cấu trúc:{
  "error": "mô tả chi tiết nguyên nhân lỗi"
}
Bảng giải nghĩa HTTP Status Code:HTTP CodeÝ nghĩaTình huống phát sinh400Bad RequestLỗi dữ liệu gửi lên: Sai định dạng JSON, gửi thừa trường, thiếu trường bắt buộc, giá $\le 0$, gửi mã ID trên URL không phải số nguyên.401UnauthorizedLỗi định danh: Sai username/password khi đăng nhập, thiếu token JWT, token hết hạn hoặc token sai chữ ký.403ForbiddenLỗi phân quyền: Đã có token nhưng cố tình xóa hoặc đóng phiên đấu giá của người khác.404Not FoundKhông tìm thấy tài nguyên: Phiên đấu giá (Auction ID) không tồn tại.409ConflictXung đột nghiệp vụ: Đăng ký trùng username, đặt giá vào phiên đã đóng/hết hạn, cố xóa phiên đã có lượt trả giá.500Internal ServerLỗi hệ thống: Lỗi kết nối DB hoặc lỗi không xác định từ máy chủ.Hướng dẫn Cài đặt & Chạy (How to run)Cấu hình Biến môi trường (.env)Biến môi trườngGiải thíchGiá trị mặc địnhJWT_SECRETChìa khóa ký JWT (Bắt buộc phải đổi)(Không có)DB_PATHĐường dẫn lưu file cơ sở dữ liệu SQLitedata/auction.dbSERVER_PORTCổng HTTP mà server lắng nghe8080Cách 1: Chạy bằng Docker Compose (Khuyên dùng)Yêu cầu đã cài đặt Docker. Khởi tạo file .env chứa biến JWT_SECRET trước khi chạy.docker compose up -d
Hệ thống sẽ chạy ngầm và tự động mount dữ liệu vào thư mục data/ để không bị mất dữ liệu khi xóa container.Cách 2: Chạy trực tiếp từ Source CodeYêu cầu đã cài đặt Go 1.27+.# 1. Clone dự án
git clone https://github.com/dominhhieu1405/Simple-Auction-system.git
cd Simple-Auction-system

# 2. Tạo file môi trường
cp .env.example .env
# (Mở file .env và điền giá trị cho JWT_SECRET)

# 3. Chạy server
go run .
Chạy kịch bản kiểm thử (Integration Tests)Đảm bảo server đang chạy ở port 8080, mở một Terminal khác và gõ:go test ./tests -v
