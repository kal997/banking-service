# Banking Service

A comprehensive banking operations microservice built with Go, featuring account management, transaction processing, and customer data handling. This service implements core banking functionality with ACID-compliant transactions and robust security.

![Go](https://img.shields.io/badge/Go-1.24-blue.svg)
![MySQL](https://img.shields.io/badge/Database-MySQL-orange.svg)
![Architecture](https://img.shields.io/badge/Architecture-Clean-green.svg)
![Testing](https://img.shields.io/badge/Coverage-95%25-brightgreen.svg)

## 🎓 Part of Banking Microservices System

This service is the core component of a comprehensive banking microservices project inspired by the **[Building Microservices API in Go](https://www.coursera.org/learn/packt-building-microservices-api-in-go-bq6wv)** course on Coursera.

**Related Services:**
- 🔐 **[Banking Auth](https://github.com/kal997/banking-auth)** - JWT authentication and authorization
- 📚 **[Banking Lib](https://github.com/kal997/banking-lib)** - Shared utilities and error handling
- 📖 **[Complete System Documentation](https://github.com/kal997/banking-microservices-docs)** - Full architecture guide

## 🏦 Features

### **Customer Management**
- **Customer data retrieval** with filtering options
- **Status-based filtering** (active/inactive customers)
- **Individual customer lookup** by ID
- **Role-based access control** (admin/user permissions)

### **Account Management**
- **Account creation** with validation (minimum $5,000 deposit)
- **Account type support** (saving/checking accounts)
- **Account ownership** verification for users
- **Database-generated** unique account IDs

### **Transaction Processing**
- **Deposit and withdrawal** operations
- **ACID-compliant transactions** with proper rollback
- **Insufficient funds** validation for withdrawals
- **Real-time balance updates** with transaction history
- **Atomic operations** ensuring data consistency

### **Security & Integration**
- **JWT-based authentication** via Auth Service
- **Cross-service authorization** validation
- **Request validation** and input sanitization
- **Structured error handling** with proper HTTP status codes

### **Architecture Excellence**
- **Clean Architecture** with clear layer separation
- **Hexagonal Architecture** (ports and adapters)
- **Repository pattern** with interface abstraction
- **Dependency injection** for testability
- **Comprehensive testing** with mock implementations

## 🏗️ Architecture

### **Service Architecture**
```
┌─────────────────────┐
│   HTTP Handlers     │ ← Customer, Account endpoints
├─────────────────────┤
│   Auth Middleware   │ ← JWT validation via Auth Service
├─────────────────────┤
│   Business Services │ ← Account, Customer, Transaction logic
├─────────────────────┤
│   Domain Models     │ ← Customer, Account, Transaction entities
├─────────────────────┤
│   Repository Layer  │ ← Database operations with sqlx
└─────────────────────┘
```

### **Database Schema**
```sql
customers (customer_id, name, city, zipcode, date_of_birth, status)
    ↓
accounts (account_id, customer_id, account_type, amount, opening_date, status)
    ↓
transactions (transaction_id, account_id, amount, transaction_type, transaction_date)
```

### **Dependencies**
- **[Banking Auth](https://github.com/kal997/banking-auth)**: Authentication and authorization
- **[Banking Lib](https://github.com/kal997/banking-lib)**: Shared logging and error handling
- **MySQL**: Core banking data storage
- **Gorilla Mux**: HTTP routing and middleware

## 🚀 Quick Start

### Prerequisites
- **Go 1.24+**
- **MySQL 8.0+**
- **Banking Auth Service** running on port 8080
- **Banking Lib** dependency

### 1. Clone and Setup
```bash
git clone https://github.com/kal997/banking-service.git
cd banking-service
go mod tidy
```

### 2. Database Setup
```sql
-- Create database
CREATE DATABASE banking_db;
USE banking_db;

-- Create tables
CREATE TABLE customers (
    customer_id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    zipcode VARCHAR(10) NOT NULL,
    date_of_birth DATE NOT NULL,
    status TINYINT(1) DEFAULT 1
);

CREATE TABLE accounts (
    account_id INT PRIMARY KEY AUTO_INCREMENT,
    customer_id INT NOT NULL,
    opening_date DATETIME NOT NULL,
    account_type VARCHAR(10) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status TINYINT(1) DEFAULT 1,
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
);

CREATE TABLE transactions (
    transaction_id INT PRIMARY KEY AUTO_INCREMENT,
    account_id INT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    transaction_type VARCHAR(10) NOT NULL,
    transaction_date DATETIME NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- Insert sample data
INSERT INTO customers (customer_id, name, city, zipcode, date_of_birth, status) VALUES
(2000, 'Khaled Hussein', 'Cairo', '12577', '1997-06-04', 1),
(2001, 'Ahmed Ali', 'Giza', '12345', '1990-03-15', 1);

INSERT INTO accounts (account_id, customer_id, opening_date, account_type, amount, status) VALUES
(95470, 2000, '2023-01-15 10:30:00', 'saving', 25000.00, 1),
(95471, 2000, '2023-02-20 14:15:00', 'checking', 5000.00, 1);
```

### 3. Configuration
Create environment variables:
```bash
export SERVER_ADDRESS=localhost
export SERVER_PORT=8000
export AUTH_SERVER_PORT=8080
export DB_USER=your_db_user
export DB_PASSWD=your_db_password
export DB_ADDR=localhost
export DB_PORT=3306
export DB_NAME=banking_db
```

### 4. Start Dependencies
```bash
# Start Banking Auth Service first
cd ../banking-auth && go run main.go

# In another terminal, start Banking Service
cd banking-service && go run main.go
```

Service will start on `http://localhost:8000`

## 📖 API Documentation

### **Authentication**
All endpoints require a valid JWT token from Banking Auth Service:
```
Authorization: Bearer <your-jwt-token>
```

Get token from: `POST http://localhost:8080/auth/login`

### **Customer Endpoints**

#### **Get All Customers**
```http
GET /customers
Authorization: Bearer <jwt-token>
```

**Query Parameters:**
- `status` (optional): Filter by `active` or `inactive`

**Response:**
```json
[
  {
    "customer_id": "2000",
    "name": "Khaled Hussein",
    "city": "Cairo",
    "zipcode": "12577",
    "date_of_birth": "1997-06-04",
    "status": "Active"
  }
]
```

#### **Get Customer by ID**
```http
GET /customers/{customer_id}
Authorization: Bearer <jwt-token>
```

### **Account Endpoints**

#### **Create New Account**
```http
POST /customers/{customer_id}/account
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "account_type": "saving",
  "amount": 10000.00
}
```

**Validation Rules:**
- `amount` must be ≥ $5,000
- `account_type` must be `saving` or `checking`

**Response:**
```json
{
  "account_id": "95472"
}
```

### **Transaction Endpoints**

#### **Make Transaction**
```http
POST /customers/{customer_id}/account/{account_id}
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "amount": 1000.00,
  "transaction_type": "deposit"
}
```

**Transaction Types:**
- `deposit`: Add money to account
- `withdrawal`: Remove money from account (with balance validation)

**Response:**
```json
{
  "transaction_id": "12345",
  "account_id": "95470",
  "new_balance": 26000.00,
  "transaction_type": "deposit",
  "transaction_date": "2024-01-15 14:30:00"
}
```

## 🔐 Security & Authorization

### **Authentication Flow**
1. Client sends JWT token in `Authorization: Bearer <token>` header
2. Auth Middleware extracts token and calls Banking Auth Service
3. Auth Service validates token and checks route permissions
4. If authorized, request proceeds; otherwise returns 403 Forbidden

### **Permission Levels**
- **Admin**: Can access all customers and perform all operations
- **User**: Can only access their own customer data and accounts

### **Route Protection**
```go
// All routes protected by Auth Middleware
router.Use(am.authorizationHandler())

// Route names used for permission checking:
// - "GetAllCustomers" (admin only)
// - "GetCustomer" (admin or own customer)
// - "NewAccount" (admin or own customer)
// - "NewTransaction" (admin or own account)
```

## 🧪 Testing

### **Run All Tests**
```bash
go test ./... -v -cover
```

### **Test Coverage**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### **Test Categories**

#### **Unit Tests**
- **Service layer tests** with mocked repositories
- **Domain model validation** tests
- **DTO validation** tests
- **Business logic** edge cases


#### **Example Test Results**
```
=== RUN   Test_should_return_customer_with_status_code_200
--- PASS: Test_should_return_customer_with_status_code_200 (0.00s)
=== RUN   Test_MakeTransaction_should_return_a_validation_error_response_when_the_request_balance_is_not_sufficient
--- PASS: Test_MakeTransaction_should_return_a_validation_error_response_when_the_request_balance_is_not_sufficient (0.00s)

PASS
coverage: 95.2% of statements
```

## 🗄️ Database Design

### **Transaction Handling**
```go
// Atomic transaction with proper rollback
tx, err := d.client.Begin()
if err != nil {
    return nil, errs.NewUnexpectedError("Unexpected database error")
}

// Insert transaction record
result, err := tx.Exec(sqlNewTransaction, transaction.AccountId, transaction.Amount, transaction.TransactionType, transaction.TransactionDate)
if err != nil {
    tx.Rollback()
    return nil, errs.NewUnexpectedError(err.Error())
}

// Update account balance atomically
if transaction.IsWithdrawal() {
    _, err = tx.Exec(`UPDATE accounts SET amount = amount - ? where account_id = ?`, transaction.Amount, transaction.AccountId)
} else {
    _, err = tx.Exec(`UPDATE accounts SET amount = amount + ? where account_id = ?`, transaction.Amount, transaction.AccountId)
}

if err != nil {
    tx.Rollback()
    return nil, errs.NewUnexpectedError(err.Error())
}

// Commit transaction
err = tx.Commit()
```


## 🚀 Deployment

### **Docker**
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o banking-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/banking-service .
CMD ["./banking-service"]
```

### **Docker Compose**
```yaml
version: '3.8'
services:
  banking-service:
    build: .
    ports:
      - "8000:8000"
    environment:
      - SERVER_ADDRESS=0.0.0.0
      - AUTH_SERVER_PORT=8080
      - DB_ADDR=mysql
    depends_on:
      - mysql
      - banking-auth
```


## 🔧 Configuration

### **Environment Variables**

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_ADDRESS` | Server bind address | `localhost` |
| `SERVER_PORT` | Server port | `8000` |
| `AUTH_SERVER_PORT` | Auth service port | `8080` |
| `DB_USER` | Database username | - |
| `DB_PASSWD` | Database password | - |
| `DB_ADDR` | Database address | `localhost` |
| `DB_PORT` | Database port | `3306` |
| `DB_NAME` | Database name | - |

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/new-feature`
3. Write tests for new functionality
4. Ensure all tests pass: `go test ./... -cover`
5. Update documentation as needed
6. Commit changes: `git commit -am 'Add new feature'`
7. Push branch: `git push origin feature/new-feature`
8. Create Pull Request

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details.

## 👤 Author

**Khaled Hussein** - *Backend Engineer*
- 📧 Email: khaled.soliman97@gmail.com
- 💼 LinkedIn: [linkedin.com/in/khaled-soliman-ali](https://linkedin.com/in/khaled-soliman-ali)
- 🐙 GitHub: [@kal997](https://github.com/kal997)

---

**Built with 🏦 and Go | Core of Banking Microservices System**
