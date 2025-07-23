
# UpdateWithValidation

`UpdateWithValidation` คือ generic helper สำหรับอัปเดตข้อมูลในฐานข้อมูล โดยสามารถกำหนดฟังก์ชันตรวจสอบสิทธิ์หรือเงื่อนไขก่อนอัปเดตได้ เหมาะกับกรณีที่ต้องเช็ค role, owner, หรือสถานะ resource ก่อนอัปเดต

## ฟังก์ชัน

```go
func UpdateWithValidation[M any, S any](
    c *fiber.Ctx,
    db *gorm.DB,
    validateFn func(c *fiber.Ctx, db *gorm.DB, args ...any) error,
    where string,
    whereArgs []any,
    preload ...string,
) error
```

### Parameters
- `M` : ประเภท model (เช่น `models.Product`)
- `S` : ประเภท schema สำหรับรับ input (เช่น `schema.UpdateProduct`)
- `validateFn` : ฟังก์ชันสำหรับตรวจสอบสิทธิ์/สถานะก่อนอัปเดต
- `where` : SQL where clause (เช่น `"id = ?"`)
- `whereArgs` : argument สำหรับ where clause (เช่น `[]any{id}`)
- `preload` : รายชื่อ relation ที่ต้องการ preload (optional)

---

## ตัวอย่างการใช้งาน

### Use Case 1: อัปเดตสินค้า (Product)

สถานการณ์: ต้องการให้แก้ไขสินค้าได้เฉพาะ admin หรือเจ้าของสินค้า และสินค้าต้องยัง active อยู่

```go
// ฟังก์ชันตรวจสอบสิทธิ์
func canEditProduct(c *fiber.Ctx, db *gorm.DB, args ...any) error {
    user := c.Locals("user").(*jwt.Claims)
    if user == nil {
        return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
    }
    productID := args[0]

    // ตรวจสอบสถานะสินค้า
    var product models.Product
    if err := db.First(&product, productID).Error; err != nil {
        return response.Fail(c, "NOT_FOUND", "Product not found", fiber.StatusNotFound)
    }
    if !product.IsActive {
        return response.Fail(c, "FORBIDDEN", "Product is not active", fiber.StatusForbidden)
    }

    // ตรวจสอบสิทธิ์ (admin หรือเจ้าของ)
    if user.Role != "admin" && product.OwnerID != user.UserID {
        return response.Fail(c, "FORBIDDEN", "Permission denied", fiber.StatusForbidden)
    }
    return nil
}

// ใช้ใน handler
func updateProductHandler(c *fiber.Ctx) error {
    idParam := c.Params("id")
    return helper.UpdateWithValidation[models.Product, schema.UpdateProduct](
        c,
        database.DB,
        canEditProduct,
        "id = ?",
        []any{idParam},
        "Category", "Owner", // preload relations
    )
}
```

### Use Case 2: อัปเดต Ticket

สถานการณ์: อัปเดต ticket ได้เฉพาะผู้รายงาน, ผู้ที่ได้รับมอบหมาย, หรือ admin และ ticket ต้องไม่อยู่ในสถานะ "closed"

```go
func canEditTicket(c *fiber.Ctx, db *gorm.DB, args ...any) error {
    user := c.Locals("user").(*jwt.Claims)
    if user == nil {
        return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
    }
    ticketID := args[0]

    var ticket models.Ticket
    if err := db.Preload("Status").First(&ticket, ticketID).Error; err != nil {
        return response.Fail(c, "NOT_FOUND", "Ticket not found", fiber.StatusNotFound)
    }

    // ตรวจสอบสถานะ ticket
    if ticket.Status.Name == "closed" {
        return response.Fail(c, "FORBIDDEN", "Cannot edit closed ticket", fiber.StatusForbidden)
    }

    // ตรวจสอบสิทธิ์
    if user.Role != "admin" && 
       ticket.ReporterID != user.UserID && 
       ticket.AssigneeID != user.UserID {
        return response.Fail(c, "FORBIDDEN", "Permission denied", fiber.StatusForbidden)
    }
    return nil
}

func updateTicketHandler(c *fiber.Ctx) error {
    idParam := c.Params("id")
    return helper.UpdateWithValidation[models.Ticket, schema.UpdateTicket](
        c,
        database.DB,
        canEditTicket,
        "id = ?",
        []any{idParam},
        "Status", "Reporter", "Assignee", "Project",
    )
}
```

### Use Case 3: อัปเดต Project Member

สถานการณ์: แก้ไขสมาชิกใน project ได้เฉพาะ project owner หรือ admin

```go
func canEditProjectMember(c *fiber.Ctx, db *gorm.DB, args ...any) error {
    user := c.Locals("user").(*jwt.Claims)
    if user == nil {
        return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
    }
    memberID := args[0]

    var member models.ProjectMember
    if err := db.Preload("Project").First(&member, memberID).Error; err != nil {
        return response.Fail(c, "NOT_FOUND", "Project member not found", fiber.StatusNotFound)
    }

    // ตรวจสอบสิทธิ์
    if user.Role != "admin" && member.Project.OwnerID != user.UserID {
        return response.Fail(c, "FORBIDDEN", "Only project owner or admin can edit members", fiber.StatusForbidden)
    }
    return nil
}

func updateProjectMemberHandler(c *fiber.Ctx) error {
    idParam := c.Params("id")
    return helper.UpdateWithValidation[models.ProjectMember, schema.UpdateProjectMember](
        c,
        database.DB,
        canEditProjectMember,
        "id = ?",
        []any{idParam},
        "Project", "User", "Role",
    )
}
```

---

## ข้อดี

- **รวม logic ไว้ที่เดียว**: update และ validation อยู่ในฟังก์ชันเดียว
- **ยืดหยุ่น**: ใช้กับ resource อะไรก็ได้ ตั้งแต่ user, product, ticket ฯลฯ
- **ลด code ซ้ำ**: ไม่ต้องเขียน validation logic ซ้ำในแต่ละ handler
- **รองรับ preload**: สามารถ preload relation ที่ต้องการได้หลายตัว
- **Where clause ยืดหยุ่น**: ไม่จำกัดแค่ `id = ?` สามารถใช้ condition อื่นได้

---

## Best Practices

1. **ตั้งชื่อ validate function ให้ชัดเจน** เช่น `canEditProduct`, `canUpdateTicket`
2. **ตรวจสอบสิทธิ์ก่อน** แล้วค่อยตรวจสอบสถานะ resource
3. **Return response.Fail ที่เหมาะสม** เพื่อให้ client เข้าใจปัญหา
4. **Preload เฉพาะที่จำเป็น** เพื่อประสิทธิภาพ

---

**Tip:** สามารถนำ pattern นี้ไปใช้กับทุก resource ที่ต้องตรวจสอบสิทธิ์หรือสถานะก่อนอัปเดต
