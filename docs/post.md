# Create Functions Usage Guide

คู่มือการใช้งานฟังก์ชัน `Create` และ `CreateWithValidation` สำหรับสร้างข้อมูลใหม่ในฐานข้อมูล พร้อมรองรับการ validate และ preload relations

## ฟังก์ชันที่มีให้ใช้

### 1. Create - การสร้างข้อมูลแบบพื้นฐาน

```go
func Create[M any, S any](c *fiber.Ctx, db *gorm.DB, preload ...string) error
```

### 2. CreateWithValidation - การสร้างข้อมูลพร้อม custom validation

```go
func CreateWithValidation[M any, S any](
    c *fiber.Ctx,
    db *gorm.DB,
    validateFn func(c *fiber.Ctx, db *gorm.DB, schema S) error,
    preload ...string,
) error
```

---

## ตัวอย่างการใช้งาน

### Use Case 1: สร้าง Project ใหม่ (แบบพื้นฐาน)

```go
func createProjectHandler(c *fiber.Ctx) error {
    return helper.Create[models.Project, schema.CreateProject](
        c, 
        database.DB, 
        "Status", "Owner", // preload relations
    )
}
```

**Request Body:**
```json
{
    "name": "My New Project",
    "description": "Project description",
    "key": "MNP",
    "owner_id": 1
}
```

### Use Case 2: สร้าง Project พร้อมตรวจสอบสิทธิ์

สถานการณ์: เฉพาะ admin เท่านั้นที่สามารถสร้าง project ได้ และต้องตรวจสอบว่าชื่อ project ไม่ซ้ำ

```go
func canCreateProject(c *fiber.Ctx, db *gorm.DB, schema schema.CreateProject) error {
    user := c.Locals("user").(*jwt.Claims)
    if user == nil {
        return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
    }

    // ตรวจสอบ role
    if user.Role != "admin" {
        return response.Fail(c, "FORBIDDEN", "Only admin can create project", fiber.StatusForbidden)
    }

    // ตรวจสอบชื่อ project ห้ามซ้ำ
    var count int64
    db.Model(&models.Project{}).Where("name = ?", schema.Name).Count(&count)
    if count > 0 {
        return response.Fail(c, "DUPLICATE", "Project name already exists", fiber.StatusBadRequest)
    }

    return nil
}

func createProjectWithValidationHandler(c *fiber.Ctx) error {
    return helper.CreateWithValidation[models.Project, schema.CreateProject](
        c,
        database.DB,
        canCreateProject,
        "Status", "Owner",
    )
}
```

### Use Case 3: สร้าง Ticket โดยตรวจสอบ Project Status

สถานการณ์: สร้างได้เฉพาะใน project ที่ active และ user ต้องเป็นสมาชิกของ project

```go
func canCreateTicket(c *fiber.Ctx, db *gorm.DB, schema schema.CreateTicket) error {
    user := c.Locals("user").(*jwt.Claims)
    if user == nil {
        return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
    }

    // ตรวจสอบว่า project มีอยู่และ active
    var project models.Project
    if err := db.Preload("Status").First(&project, schema.ProjectID).Error; err != nil {
        return response.Fail(c, "NOT_FOUND", "Project not found", fiber.StatusNotFound)
    }
    if project.Status.Name != "active" {
        return response.Fail(c, "FORBIDDEN", "Project is not active", fiber.StatusForbidden)
    }

    // ตรวจสอบว่า user เป็นสมาชิกของ project
    var memberCount int64
    db.Model(&models.ProjectMember{}).Where("project_id = ? AND user_id = ?", schema.ProjectID, user.UserID).Count(&memberCount)
    if memberCount == 0 {
        return response.Fail(c, "FORBIDDEN", "You are not a member of this project", fiber.StatusForbidden)
    }

    return nil
}

func createTicketHandler(c *fiber.Ctx) error {
    return helper.CreateWithValidation[models.Ticket, schema.CreateTicket](
        c,
        database.DB,
        canCreateTicket,
        "Project", "Reporter", "Status",
    )
}
```

### Use Case 4: เพิ่มสมาชิกใน Project

สถานการณ์: เฉพาะ project owner หรือ admin เท่านั้นที่สามารถเพิ่มสมาชิกได้

```go
func canAddProjectMember(c *fiber.Ctx, db *gorm.DB, schema schema.CreateProjectMember) error {
    user := c.Locals("user").(*jwt.Claims)
    if user == nil {
        return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
    }

    // ตรวจสอบว่า project มีอยู่
    var project models.Project
    if err := db.First(&project, schema.ProjectID).Error; err != nil {
        return response.Fail(c, "NOT_FOUND", "Project not found", fiber.StatusNotFound)
    }

    // ตรวจสอบสิทธิ์ (owner หรือ admin)
    if user.Role != "admin" && project.OwnerID != user.UserID {
        return response.Fail(c, "FORBIDDEN", "Only project owner or admin can add members", fiber.StatusForbidden)
    }

    // ตรวจสอบว่า user ยังไม่เป็นสมาชิกอยู่แล้ว
    var memberCount int64
    db.Model(&models.ProjectMember{}).Where("project_id = ? AND user_id = ?", schema.ProjectID, schema.UserID).Count(&memberCount)
    if memberCount > 0 {
        return response.Fail(c, "DUPLICATE", "User is already a member", fiber.StatusBadRequest)
    }

    return nil
}

func addProjectMemberHandler(c *fiber.Ctx) error {
    return helper.CreateWithValidation[models.ProjectMember, schema.CreateProjectMember](
        c,
        database.DB,
        canAddProjectMember,
        "Project", "User",
    )
}
```

### Use Case 5: สร้าง Organization

สถานการณ์: ตรวจสอบชื่อ organization ห้ามซ้ำ และ user ต้องมี role เป็น admin หรือ manager

```go
func canCreateOrganization(c *fiber.Ctx, db *gorm.DB, schema schema.CreateOrganization) error {
    user := c.Locals("user").(*jwt.Claims)
    if user == nil {
        return response.Fail(c, "UNAUTHORIZED", "Unauthorized", fiber.StatusUnauthorized)
    }

    // ตรวจสอบ role
    if user.Role != "admin" && user.Role != "manager" {
        return response.Fail(c, "FORBIDDEN", "Insufficient permissions", fiber.StatusForbidden)
    }

    // ตรวจสอบชื่อ organization ห้ามซ้ำ
    var count int64
    db.Model(&models.Organization{}).Where("name = ?", schema.Name).Count(&count)
    if count > 0 {
        return response.Fail(c, "DUPLICATE", "Organization name already exists", fiber.StatusBadRequest)
    }

    return nil
}

func createOrganizationHandler(c *fiber.Ctx) error {
    return helper.CreateWithValidation[models.Organization, schema.CreateOrganization](
        c,
        database.DB,
        canCreateOrganization,
        "Owner",
    )
}
```

---

## ข้อดี

- **Generic**: ใช้ได้กับ model และ schema อะไรก็ได้
- **Flexible Validation**: กำหนด logic ตรวจสอบเองได้
- **Auto Preload**: preload relation ที่ต้องการได้หลายตัว
- **Consistent Response**: response format เหมือนกันทุกที่
- **Error Handling**: จัดการ error แบบมาตรฐาน

---

## Best Practices

1. **ตั้งชื่อ validate function ให้ชัดเจน** เช่น `canCreateProject`, `canAddMember`
2. **ตรวจสอบ authorization ก่อน** แล้วค่อยตรวจสอบ business logic
3. **Return response.Fail ที่เหมาะสม** พร้อมข้อความที่เข้าใจง่าย
4. **Preload เฉพาะที่จำเป็น** เพื่อประสิทธิภาพ
5. **ใช้ Create สำหรับกรณีธรรมดา** และ CreateWithValidation เมื่อต้องตรวจสอบ logic ซับซ้อน

---

**Tip:** GORM จะจัดการ `CreatedAt`, `UpdatedAt` ให้อัตโนมัติ ไม่ต้องส่งมาใน request body