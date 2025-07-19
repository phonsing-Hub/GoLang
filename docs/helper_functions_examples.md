# Helper Functions Usage Examples

เอกสารนี้แสดงตัวอย่างการใช้งานฟังก์ชัน Helper ต่างๆ ใน `internal/utils/helper` package

## Table of Contents
- [POST Operations](#post-operations)
  - [Create](#create)
  - [CreateWithValidation](#createwithvalidation)
  - [CreateWithTransaction](#createwithtransaction)
- [PUT Operations](#put-operations)
  - [UpdateByID](#updatebyid)
  - [UpdateByIDWithValidation](#updatebyidwithvalidation)

---

## POST Operations

### Create

ฟังก์ชัน `Create[M any, S any]` ใช้สำหรับการสร้างข้อมูลใหม่แบบพื้นฐาน

#### ตัวอย่างการใช้งาน:

```go
// Model struct สำหรับฐานข้อมูล
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Request struct สำหรับรับ input
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

// Handler function
func CreateUser(c *fiber.Ctx) error {
    return helper.Create[User, CreateUserRequest](c, database.DB)
}
```

#### Request Example:
```json
POST /users
{
    "name": "John Doe",
    "email": "john@example.com"
}
```

#### Response Example:
```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "created_at": "2025-01-17T10:30:00Z",
        "updated_at": "2025-01-17T10:30:00Z"
    }
}
```

---

### CreateWithValidation

ฟังก์ชัน `CreateWithValidation[M any, S any]` ใช้เมื่อต้องการ validation เพิ่มเติมนอกเหนือจาก struct tag validation

#### ตัวอย่างการใช้งาน:

```go
type Product struct {
    ID          uint    `json:"id" gorm:"primaryKey"`
    Name        string  `json:"name" gorm:"not null"`
    Price       float64 `json:"price" gorm:"not null"`
    CategoryID  uint    `json:"category_id" gorm:"not null"`
    Stock       int     `json:"stock" gorm:"default:0"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
    Name       string  `json:"name" validate:"required,min=2,max=200"`
    Price      float64 `json:"price" validate:"required,gt=0"`
    CategoryID uint    `json:"category_id" validate:"required"`
    Stock      int     `json:"stock" validate:"min=0"`
}

// Custom validation function
func validateProduct(product Product) error {
    // ตรวจสอบว่า category มีอยู่จริงในฐานข้อมูล
    var category Category
    if err := database.DB.First(&category, product.CategoryID).Error; err != nil {
        return fmt.Errorf("category with ID %d not found", product.CategoryID)
    }
    
    // ตรวจสอบว่าชื่อสินค้าไม่ซ้ำในหมวดหมู่เดียวกัน
    var existingProduct Product
    if err := database.DB.Where("name = ? AND category_id = ?", product.Name, product.CategoryID).First(&existingProduct).Error; err == nil {
        return fmt.Errorf("product with name '%s' already exists in this category", product.Name)
    }
    
    return nil
}

// Handler function
func CreateProduct(c *fiber.Ctx) error {
    return helper.CreateWithValidation[Product, CreateProductRequest](c, database.DB, validateProduct)
}
```

#### Request Example:
```json
POST /products
{
    "name": "MacBook Pro M3",
    "price": 89900.00,
    "category_id": 1,
    "stock": 10
}
```

---

### CreateWithTransaction

ฟังก์ชัน `CreateWithTransaction[M any, S any]` ใช้เมื่อต้องการ transaction เพื่อความปลอดภัยของข้อมูล

#### ตัวอย่างการใช้งาน:

```go
type Order struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    UserID     uint      `json:"user_id" gorm:"not null"`
    TotalPrice float64   `json:"total_price" gorm:"not null"`
    Status     string    `json:"status" gorm:"default:'pending'"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
    UserID     uint    `json:"user_id" validate:"required"`
    TotalPrice float64 `json:"total_price" validate:"required,gt=0"`
}

// Handler function สำหรับการสั่งซื้อที่ต้องการความปลอดภัยสูง
func CreateOrder(c *fiber.Ctx) error {
    return helper.CreateWithTransaction[Order, CreateOrderRequest](c, database.DB)
}
```

#### Use Cases สำหรับ Transaction:
- การสร้าง Order พร้อมกับ OrderItems
- การโอนเงินระหว่างบัญชี
- การอัพเดทข้อมูลหลายตารางพร้อมกัน
- การดำเนินการที่ต้องการ ACID properties

---

## PUT Operations

### UpdateByID

ฟังก์ชัน `UpdateByID[M any, S any]` ใช้สำหรับการอัพเดทข้อมูลตาม ID

#### ตัวอย่างการใช้งาน:

```go
type UpdateUserRequest struct {
    Name  string `json:"name" validate:"omitempty,min=2,max=100"`
    Email string `json:"email" validate:"omitempty,email"`
}

// Handler function
func UpdateUser(c *fiber.Ctx) error {
    return helper.UpdateByID[User, UpdateUserRequest](c, database.DB)
}
```

#### Request Example:
```json
PUT /users/1
{
    "name": "John Smith",
    "email": "john.smith@example.com"
}
```

#### Response Example:
```json
{
    "success": true,
    "data": {
        "id": 1,
        "name": "John Smith",
        "email": "john.smith@example.com",
        "created_at": "2025-01-17T10:30:00Z",
        "updated_at": "2025-01-17T14:15:00Z"
    }
}
```

#### Protected Fields
ฟังก์ชันนี้จะป้องกันการอัพเดทฟิลด์ต่อไปนี้โดยอัตโนมัติ:
- `id`, `ID`
- `created_at`, `CreatedAt`
- `updated_at`, `UpdatedAt`
- `deleted_at`, `DeletedAt`

---

### UpdateByIDWithValidation

ฟังก์ชัน `UpdateByIDWithValidation[M any, S any]` ใช้เมื่อต้องการ validation เพิ่มเติมก่อนอัพเดท

#### ตัวอย่างการใช้งาน:

```go
type UpdateProductRequest struct {
    Name       string  `json:"name" validate:"omitempty,min=2,max=200"`
    Price      float64 `json:"price" validate:"omitempty,gt=0"`
    CategoryID uint    `json:"category_id" validate:"omitempty"`
    Stock      int     `json:"stock" validate:"omitempty,min=0"`
}

// Custom validation function สำหรับการอัพเดท
func validateProductUpdate(product Product) error {
    // ตรวจสอบว่า category มีอยู่จริง (ถ้ามีการเปลี่ยน category)
    if product.CategoryID > 0 {
        var category Category
        if err := database.DB.First(&category, product.CategoryID).Error; err != nil {
            return fmt.Errorf("category with ID %d not found", product.CategoryID)
        }
    }
    
    // ตรวจสอบ business rules อื่นๆ
    if product.Stock < 0 {
        return fmt.Errorf("stock cannot be negative")
    }
    
    // ตรวจสอบราคาไม่ต่ำกว่าขั้นต่ำ
    if product.Price > 0 && product.Price < 10.0 {
        return fmt.Errorf("price must be at least 10.00")
    }
    
    return nil
}

// Handler function
func UpdateProduct(c *fiber.Ctx) error {
    return helper.UpdateByIDWithValidation[Product, UpdateProductRequest](c, database.DB, validateProductUpdate)
}
```

#### Request Example:
```json
PUT /products/1
{
    "name": "MacBook Pro M3 Max",
    "price": 99900.00,
    "stock": 5
}
```

---

## Route Registration Examples

### การลงทะเบียน Routes

```go
package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/phonsing-Hub/GoLang/internal/utils/helper"
)

func SetupRoutes(app *fiber.App) {
    api := app.Group("/api")
    
    // User routes
    users := api.Group("/users")
    users.Post("/", CreateUser)
    users.Put("/:id", UpdateUser)
    
    // Product routes
    products := api.Group("/products")
    products.Post("/", CreateProduct)
    products.Put("/:id", UpdateProduct)
    
    // Order routes (with transaction)
    orders := api.Group("/orders")
    orders.Post("/", CreateOrder)
}
```

---

## Error Handling

### Common Error Responses

#### Validation Error:
```json
{
    "success": false,
    "error": "VALIDATION_FAILED",
    "message": "Key: 'CreateUserRequest.Email' Error:Tag 'email' got 'invalid-email'"
}
```

#### Not Found Error:
```json
{
    "success": false,
    "error": "NOT_FOUND",
    "message": "Record with ID 999 not found"
}
```

#### Database Error:
```json
{
    "success": false,
    "error": "DATABASE_ERROR",
    "message": "Failed to create record: UNIQUE constraint failed: users.email"
}
```

---

## Best Practices

### 1. Request Struct Design
```go
// ✅ Good - แยก validation tags ชัดเจน
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// ✅ Good - Update request ใช้ omitempty
type UpdateUserRequest struct {
    Name  string `json:"name" validate:"omitempty,min=2,max=100"`
    Email string `json:"email" validate:"omitempty,email"`
}
```

### 2. Model Struct Design
```go
// ✅ Good - ใส่ JSON tags และ GORM tags
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"` // ซ่อนจาก response
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 3. Custom Validation
```go
// ✅ Good - validation function ที่มีประสิทธิภาพ
func validateUser(user User) error {
    // Check business rules
    if strings.Contains(user.Email, "+") {
        return fmt.Errorf("email with '+' character is not allowed")
    }
    
    // Check uniqueness if needed
    var existingUser User
    if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
        return fmt.Errorf("email already exists")
    }
    
    return nil
}
```

### 4. Transaction Usage
```go
// ✅ ใช้ CreateWithTransaction สำหรับ:
// - การสร้างข้อมูลที่เกี่ยวข้องกับหลายตาราง
// - การดำเนินการทางการเงิน
// - การอัพเดทข้อมูลที่สำคัญ

// ✅ ใช้ Create ธรรมดาสำหรับ:
// - การสร้างข้อมูลง่ายๆ ในตารางเดียว
// - การดำเนินการที่ไม่ซับซ้อน
```

---

## Migration และ Database Setup

```go
// ตัวอย่าง migration
func AutoMigrate() {
    database.DB.AutoMigrate(
        &User{},
        &Product{},
        &Category{},
        &Order{},
        &OrderItem{},
    )
}
```

---

สำหรับข้อมูลเพิ่มเติม สามารถดูได้ที่:
- [Query Examples](query_examples.md)
- [Joins and Preloads](joins_and_preloads.md)
- [Query Parameters](query_parameters.md)
