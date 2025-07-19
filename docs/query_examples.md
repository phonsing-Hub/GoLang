# ตัวอย่างการใช้งาน Query Parameters ด้วย cURL

## การตั้งค่า
```bash
# ตั้งค่า base URL และ token
BASE_URL="http://localhost:3000/api/v1"
TOKEN="your_jwt_token_here"
```

## ตัวอย่างการใช้งาน

### 1. Pagination
```bash
# หน้าแรก 10 รายการ
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?page=1&limit=10"

# หน้า 2 20 รายการ
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?page=2&limit=20"
```

### 2. Sorting
```bash
# เรียงตามชื่อ A-Z
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?sort_by=name&sort_order=asc"

# เรียงตามวันที่สร้างใหม่สุด
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?sort_by=created_at&sort_order=desc"
```

### 3. Exact Match Filter
```bash
# หาห้องที่สร้างโดย user ID 2
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?created_by=2"

# หาห้องสาธารณะ
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?is_private=false"

# หาห้องที่มีชื่อ "General"
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?name=General"
```

### 4. Multiple Values Filter
```bash
# หาห้องที่สร้างโดย user 1, 2, หรือ 3
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?created_by=1,2,3"

# หาห้องที่เป็นทั้ง private และ public
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?is_private=true,false"
```

### 5. Search
```bash
# ค้นหาห้องที่มีชื่อมีคำว่า "chat"
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?search[name]=chat"

# ค้นหาห้องที่มีคำว่า "general" ในชื่อหรือคำอธิบาย
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?search_cols[name|description]=general"
```

### 6. Range Filter
```bash
# หาห้องที่สร้างระหว่างวันที่ 1 มกราคม 2024 ถึง 31 ธันวาคม 2024
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?filterrange[created_at]=2024-01-01|2024-12-31"

# หาห้องที่มี ID ระหว่าง 1-100
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?filterrange[id]=1|100"

# หาห้องที่สร้างหลังจากวันที่ 1 มกราคม 2024 (ไม่มีขีดจำกัดบน)
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?filterrange[created_at]=2024-01-01|-"
```

### 7. NOT IN Filter
```bash
# หาห้องที่ไม่สร้างโดย user 1 และ 2
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?filter_not[created_by]=1,2"

# หาห้องที่ไม่ใช่ private
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?filter_not[is_private]=true"
```

### 8. NULL Value Filter
```bash
# หาห้องที่ไม่มีคำอธิบาย
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?description=null"
```

### 9. Combined Filters
```bash
# หาห้องสาธารณะที่สร้างโดย user 2 และเรียงตามชื่อ
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?is_private=false&created_by=2&sort_by=name&sort_order=asc"

# ค้นหาห้องที่มีคำว่า "general" ในชื่อหรือคำอธิบาย และไม่ใช่ private
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?search_cols[name|description]=general&is_private=false"

# หาห้องที่สร้างโดย user 1 หรือ 2 ในเดือนมกราคม 2024 และแสดง 5 รายการต่อหน้า
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?created_by=1,2&filterrange[created_at]=2024-01-01|2024-01-31&limit=5"
```

### 10. Complex Search Examples
```bash
# หาห้องสาธารณะที่มีคำว่า "chat" ในชื่อ และสร้างโดย user 1-5
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?is_private=false&search[name]=chat&created_by=1,2,3,4,5"

# หาห้องที่ไม่สร้างโดย user 1 และมีคำว่า "general" ในชื่อหรือคำอธิบาย
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?filter_not[created_by]=1&search_cols[name|description]=general"

# หาห้องที่สร้างในปี 2024 และเรียงตามวันที่สร้างใหม่สุด
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/rooms?filterrange[created_at]=2024-01-01|2024-12-31&sort_by=created_at&sort_order=desc"
```

## การใช้งานกับ Users API
```bash
# หา users ที่มี email ถูกยืนยันแล้ว
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/users?is_email_verified=true"

# ค้นหา users ที่มีชื่อหรือนามสกุลมีคำว่า "john"
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/users?search_cols[first_name|last_name]=john"

# หา users ที่สร้างในเดือนมกราคม 2024
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/users?filterrange[created_at]=2024-01-01|2024-01-31"
```

## การใช้งานกับ Messages API
```bash
# หาข้อความในห้อง ID 1
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/messages?room_id=1"

# ค้นหาข้อความที่มีคำว่า "hello"
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/messages?search[content]=hello"

# หาข้อความที่ส่งโดย user 1 และ 2
curl -H "Authorization: Bearer $TOKEN" "$BASE_URL/messages?sender_id=1,2"
```

## หมายเหตุ
- เปลี่ยน `your_jwt_token_here` เป็น JWT token จริง
- เปลี่ยน `localhost:3000` เป็น URL ของเซิร์ฟเวอร์จริง
- Query parameters ทั้งหมดเป็น optional
- สามารถใช้ query parameters หลายตัวพร้อมกันได้ 