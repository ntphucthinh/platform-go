src/
├── .air.toml                       # Cấu hình Air cho hot reload
├── .env.local                      # Biến môi trường local (không đưa vào Git)
├── .env.test                       # Biến môi trường test (không đưa vào Git)
├── .gitignore                      # Danh sách file/thư mục không đưa vào Git
├── .golangci.yml                   # Cấu hình golangci-lint
├── .sample.env                     # File mẫu biến môi trường local
├── .sample.env.test                # File mẫu biến môi trường test
├── go.mod                          # Định nghĩa Go module và dependencies
├── go.sum                          # Checksum của Go dependencies
├── models.yaml                     # Định nghĩa model OpenAPI
├── openapitools.json               # Cấu hình OpenAPI Generator
├── server.yaml                     # Cấu hình server cho OpenAPI
│
├── cmd/                            # Entry point của ứng dụng
│   └── main/                       # Điểm khởi chạy ứng dụng
│
├── config/                         # Cấu hình hệ thống
│
├── constants/                      # Các hằng số
│
├── db/                             # Thành phần liên quan đến database
│   ├── migrations/                 # Database migrations
│   └── seed/                       # Dữ liệu khởi tạo
│
├── generated/                      # Code tự động sinh — KHÔNG chỉnh sửa trực tiếp
│   ├── api/                        # Code API tự động sinh
│   └── openapi/                    # Code OpenAPI tự động sinh
│
├── internal/                       # Package chỉ sử dụng bên trong project
│   │
│   ├── common/                     # Các xử lý dùng chung
│   │   ├── ctxkeys/                # Key sử dụng trong context
│   │   ├── mysqlcode/              # Mã lỗi MySQL
│   │   └── permissions/            # Xử lý quyền hạn
│   │
│   ├── controller/                 # HTTP Controller
│   │
│   ├── controller_test/            # Test cho Controller
│   │
│   ├── di/                         # Dependency Injection
│   │
│   ├── domain/                     # Domain logic và business rules
│   │   ├── administrator/          # Domain quản trị viên
│   │   ├── user/                   # Domain người dùng
│   │   └── .../                    # Các domain khác
│   │
│   ├── infrastructure/             # Infrastructure layer
│   │   ├── cache/                  # Cache
│   │   ├── dbmodel/                # Database models
│   │   ├── dbutil/                 # Database utilities
│   │   ├── qsdto/                  # Query Service DTO
│   │   ├── queryservice/           # Query Service
│   │   ├── repository/             # Repository implementations
│   │   ├── scope/                  # Scope implementations
│   │   └── security/               # Security
│   │
│   ├── middleware/                 # HTTP Middleware
│   │
│   ├── pkg/                        # Package tiện ích dùng chung
│   │   ├── ai/                     # AI
│   │   ├── app/                    # Application utilities
│   │   ├── cvalidator/             # Custom validator
│   │   ├── debug/                  # Debug utilities
│   │   ├── deepcopy/               # Deep copy utilities
│   │   ├── helper/                 # Helper functions
│   │   ├── log/                    # Logging
│   │   ├── safeconv/               # Chuyển đổi kiểu an toàn
│   │   ├── sliceutil/              # Slice utilities
│   │   ├── strings/                # String utilities
│   │   └── types/                  # Custom types
│   │
│   └── usecases/                   # Application use cases
│
├── logs/                           # Log files
│
├── mocks/                          # Mock objects phục vụ testing
│
├── openapi/                        # OpenAPI definitions được quản lý thủ công
│   ├── openapi.yml                 # OpenAPI definition chính
│   ├── components/                 # OpenAPI components
│   │   └── schemas/                # Schema definitions
│   └── paths/                      # API path definitions
│
├── testonly/                       # Utilities chỉ dùng cho test
│
├── tests/                          # Thành phần liên quan đến testing
│   ├── factory/                    # Test data factories
│   └── utils/                      # Test utilities
│
└── tmp/                            # File tạm thời