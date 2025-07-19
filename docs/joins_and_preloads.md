# การใช้งาน FindAll กับ Joins และ Preloads

## ภาพรวม
เมื่อต้องการดึงข้อมูลที่มีการ join หรือ preload relationships เราไม่สามารถใช้ `FindAll` ปกติได้ ต้องใช้ฟังก์ชันพิเศษที่รองรับ preloading

## ฟังก์ชันที่รองรับ Preloading

### 1. FindAllWithPreload
```go
func FindAllWithPreload[T any](c *fiber.Ctx, db *gorm.DB, preloads ...string) error
```

**การใช้งาน:**
```go
// โหลด rooms พร้อมกับ members และ created_by_user
return helper.FindAllWithPreload[models.Room](c, database.DB, 
    "Members", 
    "Members.User", 
    "CreatedByUser")
```

### 2. FindAllWithJoins
```go
func FindAllWithJoins[T any](c *fiber.Ctx, db *gorm.DB, joins []string, preloads []string) error
```

**การใช้งาน:**
```go
// ใช้ joins และ preloads พร้อมกัน
joins := []string{
    "JOIN room_members ON rooms.id = room_members.room_id",
    "JOIN users ON room_members.user_id = users.id",
}
preloads := []string{"Members", "Messages"}

return helper.FindAllWithJoins[models.Room](c, database.DB, joins, preloads)
```

### 3. FindByIDWithPreload
```go
func FindByIDWithPreload[T any](c *fiber.Ctx, db *gorm.DB, preloads ...string) error
```

**การใช้งาน:**
```go
// โหลด room เดียวพร้อมข้อมูลทั้งหมด
return helper.FindByIDWithPreload[models.Room](c, database.DB, 
    "Members", 
    "Members.User", 
    "Messages", 
    "Messages.Sender", 
    "CreatedByUser")
```

## ตัวอย่างการใช้งานใน Routes

### 1. ห้องที่มีสมาชิก (Preload Members)
```go
func get_rooms_with_members(c *fiber.Ctx) error {
    return helper.FindAllWithPreload[models.Room](c, database.DB, 
        "Members", 
        "Members.User", 
        "CreatedByUser")
}
```

**API Endpoint:** `GET /api/v1/rooms/with-members`

**Query Parameters ที่ใช้ได้:**
```bash
# หาห้องที่สร้างโดย user 2 พร้อมสมาชิก
GET /api/v1/rooms/with-members?created_by=2

# ค้นหาห้องที่มีคำว่า "general" พร้อมสมาชิก
GET /api/v1/rooms/with-members?search[name]=general

# หาห้องสาธารณะพร้อมสมาชิก เรียงตามจำนวนสมาชิก
GET /api/v1/rooms/with-members?is_private=false&sort_by=id&sort_order=desc
```

### 2. ห้องที่มีข้อความ (Preload Messages)
```go
func get_rooms_with_messages(c *fiber.Ctx) error {
    return helper.FindAllWithPreload[models.Room](c, database.DB, 
        "Messages", 
        "Messages.Sender")
}
```

**API Endpoint:** `GET /api/v1/rooms/with-messages`

**Query Parameters ที่ใช้ได้:**
```bash
# หาห้องที่มีข้อความจาก user 1
GET /api/v1/rooms/with-messages?messages.sender_id=1

# ค้นหาห้องที่มีข้อความมีคำว่า "hello"
GET /api/v1/rooms/with-messages?search[messages.content]=hello
```

### 3. ห้องเดียวพร้อมรายละเอียดทั้งหมด
```go
func get_room_with_details(c *fiber.Ctx) error {
    return helper.FindByIDWithPreload[models.Room](c, database.DB, 
        "Members", 
        "Members.User", 
        "Messages", 
        "Messages.Sender", 
        "CreatedByUser")
}
```

**API Endpoint:** `GET /api/v1/rooms/{id}/with-details`

## ตัวอย่าง Response

### 1. ห้องพร้อมสมาชิก
```json
{
  "success": true,
  "data": {
    "total": 5,
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
          },
          {
            "id": 2,
            "room_id": 1,
            "user_id": 2,
            "role": "member",
            "user": {
              "id": 2,
              "email": "user2@example.com",
              "first_name": "Jane",
              "last_name": "Smith"
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
    ]
  }
}
```

### 2. ห้องพร้อมข้อความ
```json
{
  "success": true,
  "data": {
    "total": 3,
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
        "updated_at": "2024-01-01T00:00:00Z",
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
          },
          {
            "id": 2,
            "room_id": 1,
            "sender_id": 2,
            "content": "Hi John!",
            "message_type": "text",
            "created_at": "2024-01-01T10:05:00Z",
            "sender": {
              "id": 2,
              "email": "user2@example.com",
              "first_name": "Jane",
              "last_name": "Smith"
            }
          }
        ]
      }
    ]
  }
}
```

## การใช้งานกับ Models อื่นๆ

### Users with Credentials
```go
func get_users_with_credentials(c *fiber.Ctx) error {
    return helper.FindAllWithPreload[models.Users](c, database.DB, 
        "Credentials", 
        "CurrentStatus")
}
```

### Messages with Sender and Room
```go
func get_messages_with_details(c *fiber.Ctx) error {
    return helper.FindAllWithPreload[models.Message](c, database.DB, 
        "Sender", 
        "Room", 
        "ReplyTo")
}
```

### Private Rooms with Users
```go
func get_private_rooms_with_users(c *fiber.Ctx) error {
    return helper.FindAllWithPreload[models.PrivateRoom](c, database.DB, 
        "User1", 
        "User2", 
        "Messages")
}
```

## ข้อควรระวัง

1. **Performance**: การใช้ preload มากเกินไปอาจทำให้ query ช้า
2. **N+1 Problem**: ใช้ preload เพื่อหลีกเลี่ยง N+1 query problem
3. **Memory Usage**: การโหลดข้อมูลมากเกินไปอาจใช้ memory มาก
4. **Query Parameters**: Query parameters ยังคงทำงานได้เหมือนเดิม

## Best Practices

1. **เลือกเฉพาะข้อมูลที่จำเป็น**: ไม่ต้อง preload ทุก relationship
2. **ใช้ pagination**: เพื่อจำกัดจำนวนข้อมูลที่โหลด
3. **ใช้ query parameters**: เพื่อ filter ข้อมูลก่อน preload
4. **Monitor Performance**: ตรวจสอบ query performance เมื่อใช้ preload

## ตัวอย่างการใช้งานจริง

```bash
# หาห้องสาธารณะที่มีสมาชิกมากกว่า 5 คน
GET /api/v1/rooms/with-members?is_private=false&limit=20

# ค้นหาห้องที่มีข้อความจาก user 1 ในเดือนมกราคม
GET /api/v1/rooms/with-messages?filterrange[created_at]=2024-01-01|2024-01-31

# หาห้องที่สร้างโดย user 2 พร้อมสมาชิก เรียงตามวันที่สร้างใหม่สุด
GET /api/v1/rooms/with-members?created_by=2&sort_by=created_at&sort_order=desc
``` 