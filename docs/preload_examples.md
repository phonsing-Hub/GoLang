# ตัวอย่างการใช้งาน Preloading Endpoints

## การตั้งค่า
```bash
# ตั้งค่า base URL และ token
BASE_URL="http://localhost:3000/api/v1"
TOKEN="your_jwt_token_here"
```

## 1. ห้องพร้อมสมาชิก (GET /api/v1/rooms/with-members)

### Basic Usage
```bash
# ดึงห้องทั้งหมดพร้อมสมาชิก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members"

# ดึงห้องพร้อมสมาชิก 5 รายการแรก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?limit=5"
```

### Filtering
```bash
# หาห้องสาธารณะพร้อมสมาชิก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?is_private=false"

# หาห้องที่สร้างโดย user 2 พร้อมสมาชิก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?created_by=2"

# หาห้องที่มีชื่อมีคำว่า "general" พร้อมสมาชิก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?search[name]=general"
```

### Sorting
```bash
# เรียงห้องตามชื่อ A-Z พร้อมสมาชิก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?sort_by=name&sort_order=asc"

# เรียงห้องตามวันที่สร้างใหม่สุด พร้อมสมาชิก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?sort_by=created_at&sort_order=desc"
```

### Combined Filters
```bash
# หาห้องสาธารณะที่สร้างโดย user 1 หรือ 2 พร้อมสมาชิก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?is_private=false&created_by=1,2"

# ค้นหาห้องที่มีคำว่า "chat" ในชื่อหรือคำอธิบาย พร้อมสมาชิก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?search_cols[name|description]=chat"

# หาห้องที่สร้างในเดือนมกราคม 2024 พร้อมสมาชิก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?filterrange[created_at]=2024-01-01|2024-01-31"
```

## 2. ห้องพร้อมข้อความ (GET /api/v1/rooms/with-messages)

### Basic Usage
```bash
# ดึงห้องทั้งหมดพร้อมข้อความ
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-messages"

# ดึงห้องพร้อมข้อความ 3 รายการแรก
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-messages?limit=3"
```

### Filtering
```bash
# หาห้องสาธารณะพร้อมข้อความ
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-messages?is_private=false"

# หาห้องที่สร้างโดย user 1 พร้อมข้อความ
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-messages?created_by=1"
```

### Complex Filters
```bash
# หาห้องที่มีข้อความจาก user 1 หรือ 2
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-messages?messages.sender_id=1,2"

# หาห้องที่สร้างในปี 2024 พร้อมข้อความ
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-messages?filterrange[created_at]=2024-01-01|2024-12-31"
```

## 3. ห้องเดียวพร้อมรายละเอียด (GET /api/v1/rooms/{id}/with-details)

### Basic Usage
```bash
# ดึงห้อง ID 1 พร้อมรายละเอียดทั้งหมด
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/1/with-details"

# ดึงห้อง ID 2 พร้อมรายละเอียดทั้งหมด
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/2/with-details"
```

## 4. เปรียบเทียบ Response

### ห้องปกติ (ไม่มี preload)
```bash
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/1"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "General",
    "description": "General chat room",
    "is_private": false,
    "created_by": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### ห้องพร้อมสมาชิก (มี preload)
```bash
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/1/with-details"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "General",
    "description": "General chat room",
    "is_private": false,
    "created_by": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "members": [
      {
        "id": 1,
        "room_id": 1,
        "user_id": 1,
        "role": "admin",
        "user": {
          "id": 1,
          "email": "user1@example.com",
          "first_name": "John",
          "last_name": "Doe"
        }
      }
    ],
    "messages": [
      {
        "id": 1,
        "room_id": 1,
        "sender_id": 1,
        "content": "Hello everyone!",
        "message_type": "text",
        "created_at": "2024-01-01T10:00:00Z",
        "sender": {
          "id": 1,
          "email": "user1@example.com",
          "first_name": "John",
          "last_name": "Doe"
        }
      }
    ],
    "created_by_user": {
      "id": 1,
      "email": "user1@example.com",
      "first_name": "John",
      "last_name": "Doe"
    }
  }
}
```

## 5. ตัวอย่างการใช้งานจริง

### สำหรับ Chat App Frontend

#### 1. โหลดรายการห้องพร้อมสมาชิก
```bash
# สำหรับหน้าแสดงรายการห้อง
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?limit=20&sort_by=updated_at&sort_order=desc"
```

#### 2. โหลดห้องเดียวพร้อมข้อความล่าสุด
```bash
# สำหรับหน้าแชท
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/1/with-details"
```

#### 3. ค้นหาห้องที่มีสมาชิก
```bash
# สำหรับการค้นหาห้อง
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?search[name]=general&is_private=false"
```

#### 4. โหลดห้องที่สร้างโดยตัวเอง
```bash
# สำหรับห้องที่สร้างเอง
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?created_by=1&sort_by=created_at&sort_order=desc"
```

## 6. Performance Tips

### 1. ใช้ Pagination
```bash
# จำกัดจำนวนรายการเพื่อเพิ่ม performance
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?limit=10&page=1"
```

### 2. ใช้ Filtering ก่อน Preload
```bash
# Filter ก่อนเพื่อลดข้อมูลที่ต้อง preload
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/with-members?is_private=false&limit=5"
```

### 3. ใช้เฉพาะข้อมูลที่จำเป็น
```bash
# ใช้ /with-members เมื่อต้องการดูสมาชิก
# ใช้ /with-messages เมื่อต้องการดูข้อความ
# ใช้ /with-details เมื่อต้องการข้อมูลทั้งหมด
```

## 7. Error Handling

### 404 - ห้องไม่พบ
```bash
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms/999/with-details"
```

**Response:**
```json
{
  "success": false,
  "error": "NOT_FOUND",
  "message": "Record with ID 999 not found"
}
```

### 401 - ไม่มี Authorization
```bash
curl "$BASE_URL/rooms/with-members"
```

**Response:**
```json
{
  "success": false,
  "error": "UNAUTHORIZED",
  "message": "Missing or invalid token"
}
```

## หมายเหตุ
- เปลี่ยน `your_jwt_token_here` เป็น JWT token จริง
- เปลี่ยน `localhost:3000` เป็น URL ของเซิร์ฟเวอร์จริง
- การใช้ preload จะทำให้ response ใหญ่ขึ้นและ query ช้าลง
- ใช้เฉพาะเมื่อจำเป็นต้องมีข้อมูล relationships 