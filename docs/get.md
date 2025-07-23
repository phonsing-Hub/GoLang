# Query Parameters สำหรับ API Endpoints

## ภาพรวม
ระบบรองรับ query parameters หลากหลายรูปแบบสำหรับการ filter, search, sort และ pagination

## Basic Query Parameters

### 1. Pagination
```
GET /api/v1/rooms?page=1&limit=10
```
- `page`: หมายเลขหน้า (default: 1)
- `limit`: จำนวนรายการต่อหน้า (default: 10)

### 2. Sorting
```
GET /api/v1/rooms?sort_by=name&sort_order=desc
```
- `sort_by`: ฟิลด์ที่ต้องการ sort (default: id)
- `sort_order`: asc หรือ desc (default: asc)

## Filtering Parameters

### 1. Exact Match Filter
```
GET /api/v1/rooms?is_private=true
GET /api/v1/rooms?created_by=2
GET /api/v1/rooms?name=General
```

### 2. Multiple Values Filter (IN clause)
```
GET /api/v1/rooms?created_by=1,2,3
GET /api/v1/rooms?is_private=true,false
```

### 3. NULL Value Filter
```
GET /api/v1/rooms?description=null
```

### 4. Range Filter
```
GET /api/v1/rooms?filterrange[created_at]=2024-01-01|2024-12-31
GET /api/v1/rooms?filterrange[id]=1|100
```
รูปแบบ: `filterrange[column_name]=min_value|max_value`
- ใช้ `-` สำหรับไม่มีขีดจำกัด: `filterrange[created_at]=2024-01-01|-`

### 5. NOT IN Filter
```
GET /api/v1/rooms?filter_not[created_by]=1,2
GET /api/v1/rooms?filter_not[is_private]=true
```

## Search Parameters

### 1. Single Column Search
```
GET /api/v1/rooms?search[name]=general
GET /api/v1/rooms?search[description]=chat
```
รูปแบบ: `search[column_name]=search_term`
- ใช้ ILIKE (case-insensitive) สำหรับการค้นหา

### 2. Multiple Columns Search
```
GET /api/v1/rooms?search_cols[name|description]=general
```
รูปแบบ: `search_cols[column1|column2|column3]=search_term`
- ค้นหาในหลายคอลัมน์พร้อมกัน (OR condition)

## ตัวอย่างการใช้งานสำหรับ Room Model

### ฟิลด์ที่มีใน Room Model:
- `id` (uint)
- `name` (string)
- `description` (string)
- `is_private` (bool)
- `created_by` (uint)
- `created_at` (time.Time)
- `updated_at` (time.Time)
- `deleted_at` (time.Time)

### ตัวอย่างการใช้งาน:

#### 1. หาห้องที่สร้างโดย user ID 2
```
GET /api/v1/rooms?created_by=2
```

#### 2. หาห้องสาธารณะ (ไม่ใช่ private)
```
GET /api/v1/rooms?is_private=false
```

#### 3. ค้นหาห้องที่มีชื่อหรือคำอธิบายมีคำว่า "general"
```
GET /api/v1/rooms?search_cols[name|description]=general
```

#### 4. หาห้องที่สร้างในช่วงเวลาหนึ่ง
```
GET /api/v1/rooms?filterrange[created_at]=2024-01-01|2024-12-31
```

#### 5. หาห้องที่ไม่ใช่ private และสร้างโดย user 1 หรือ 2
```
GET /api/v1/rooms?is_private=false&created_by=1,2
```

#### 6. เรียงลำดับตามชื่อจาก A-Z และแสดง 20 รายการต่อหน้า
```
GET /api/v1/rooms?sort_by=name&sort_order=asc&limit=20
```

#### 7. หาห้องที่ไม่สร้างโดย user 1, 2, 3
```
GET /api/v1/rooms?filter_not[created_by]=1,2,3
```

#### 8. ค้นหาห้องที่มีชื่อมีคำว่า "chat"
```
GET /api/v1/rooms?search[name]=chat
```

## Response Format

```json
{
  "success": true,
  "data": {
    "total": 50,
    "page": 1,
    "limit": 10,
    "data": [
      {
        "id": 1,
        "name": "General",
        "description": "General chat room",
        "is_private": false,
        "created_by": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

## หมายเหตุ
- Query parameters ทั้งหมดเป็น optional
- การใช้ query parameters หลายตัวพร้อมกันจะใช้ AND condition
- การค้นหาใช้ ILIKE (case-insensitive) สำหรับ string fields
- Pagination จะคำนวณ total count ก่อนแล้วจึง fetch ข้อมูล 