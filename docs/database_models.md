# Database Models Overview

เอกสารนี้อธิบายโครงสร้างของ database models ในระบบ Project Management และ Ticket System

## Table of Contents

1. [Status & Lookup Tables](#status--lookup-tables)
2. [User Management](#user-management)
3. [Organization Management](#organization-management)
4. [Project Management](#project-management)
5. [Ticket System](#ticket-system)
6. [Relationships Diagram](#relationships-diagram)
7. [Status Flow](#status-flow)

## Status & Lookup Tables

ระบบได้แยก status ต่างๆ เป็น lookup tables แยกต่างหากเพื่อความยืดหยุ่นและการจัดการที่ดีขึ้น

### UserStatus
สถานะของผู้ใช้ในระบบ

| Field | Type | Description | Example Values |
|-------|------|-------------|----------------|
| id | uint | Primary Key | 1, 2, 3, 4 |
| name | string | Status code | "active", "inactive", "suspended", "pending_verification" |
| display_name | string | Human readable name | "Active", "Inactive", "Suspended", "Pending Verification" |
| description | string | Status description | "User is active and can access the system" |
| color | string | Hex color for UI | "#10B981", "#6B7280", "#EF4444", "#F59E0B" |
| is_active | bool | Is this status active | true, false |
| position | int | Display order | 1, 2, 3, 4 |

**Default Data:**
- Active (id: 1)
- Inactive (id: 2)
- Suspended (id: 3)
- Pending Verification (id: 4)

### ProjectStatus
สถานะของโปรเจค

| Field | Type | Description | Example Values |
|-------|------|-------------|----------------|
| id | uint | Primary Key | 1, 2, 3, 4 |
| name | string | Status code | "active", "archived", "on_hold", "cancelled" |
| display_name | string | Human readable name | "Active", "Archived", "On Hold", "Cancelled" |

**Default Data:**
- Active (id: 1)
- Archived (id: 2)
- On Hold (id: 3)
- Cancelled (id: 4)

### EpicStatus
สถานะของ Epic

| Field | Type | Description | Example Values |
|-------|------|-------------|----------------|
| id | uint | Primary Key | 1, 2, 3, 4 |
| name | string | Status code | "planning", "in_progress", "completed", "cancelled" |
| display_name | string | Human readable name | "Planning", "In Progress", "Completed", "Cancelled" |

**Default Data:**
- Planning (id: 1)
- In Progress (id: 2)
- Completed (id: 3)
- Cancelled (id: 4)

### SprintStatus
สถานะของ Sprint

| Field | Type | Description | Example Values |
|-------|------|-------------|----------------|
| id | uint | Primary Key | 1, 2, 3, 4 |
| name | string | Status code | "planning", "active", "completed", "cancelled" |
| display_name | string | Human readable name | "Planning", "Active", "Completed", "Cancelled" |

**Default Data:**
- Planning (id: 1)
- Active (id: 2)
- Completed (id: 3)
- Cancelled (id: 4)

### OrganizationStatus
สถานะของ Organization

| Field | Type | Description | Example Values |
|-------|------|-------------|----------------|
| id | uint | Primary Key | 1, 2, 3, 4 |
| name | string | Status code | "active", "trial", "suspended", "expired" |
| display_name | string | Human readable name | "Active", "Trial", "Suspended", "Expired" |

**Default Data:**
- Active (id: 1)
- Trial (id: 2)
- Suspended (id: 3)
- Expired (id: 4)

### MemberStatus
สถานะของสมาชิกใน Organization

| Field | Type | Description | Example Values |
|-------|------|-------------|----------------|
| id | uint | Primary Key | 1, 2, 3, 4 |
| name | string | Status code | "active", "invited", "suspended", "inactive" |
| display_name | string | Human readable name | "Active", "Invited", "Suspended", "Inactive" |

**Default Data:**
- Active (id: 1)
- Invited (id: 2)
- Suspended (id: 3)
- Inactive (id: 4)

### Priority
ระดับความสำคัญ (ใช้สำหรับ Tickets และ Epics)

| Field | Type | Description | Example Values |
|-------|------|-------------|----------------|
| id | uint | Primary Key | 1, 2, 3, 4, 5 |
| name | string | Priority code | "low", "medium", "high", "critical", "blocker" |
| display_name | string | Human readable name | "Low", "Medium", "High", "Critical", "Blocker" |
| level | int | Priority level (1-5) | 1, 2, 3, 4, 5 |
| color | string | Hex color for UI | "#10B981", "#3B82F6", "#F59E0B", "#EF4444", "#7C2D12" |

**Default Data:**
- Low (id: 1, level: 1)
- Medium (id: 2, level: 2)
- High (id: 3, level: 3)
- Critical (id: 4, level: 4)
- Blocker (id: 5, level: 5)

### TicketType
ประเภทของ Ticket

| Field | Type | Description | Example Values |
|-------|------|-------------|----------------|
| id | uint | Primary Key | 1, 2, 3, 4, 5 |
| name | string | Type code | "task", "bug", "feature", "improvement", "story" |
| display_name | string | Human readable name | "Task", "Bug", "Feature", "Improvement", "User Story" |
| icon | string | Icon name for UI | "task", "bug", "feature", "improvement", "story" |
| color | string | Hex color for UI | "#3B82F6", "#EF4444", "#10B981", "#F59E0B", "#8B5CF6" |

**Default Data:**
- Task (id: 1)
- Bug (id: 2)
- Feature (id: 3)
- Improvement (id: 4)
- User Story (id: 5)

## User Management

### User
ผู้ใช้หลักในระบบ

| Field | Type | Description | FK/Relationship |
|-------|------|-------------|-----------------|
| id | uint | Primary Key | - |
| email | string | Email address (unique) | - |
| username | *string | Username (unique, nullable) | - |
| status_id | uint | User status | FK → user_statuses.id |
| is_email_verified | bool | Email verification status | - |
| last_login_at | *time.Time | Last login timestamp | - |

**Relationships:**
- `Status` → UserStatus
- `Profile` → UserProfile (1:1)
- `AuthMethods` → []UserAuthMethod (1:many)
- `OrganizationMembers` → []OrganizationMember (1:many)
- `OwnedProjects` → []Project (1:many)
- `AssignedTickets` → []Ticket (1:many)

### UserProfile
ข้อมูลโปรไฟล์ของผู้ใช้

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary Key |
| user_id | uint | FK → users.id |
| first_name | string | First name |
| last_name | string | Last name |
| display_name | string | Display name |
| bio | string | Biography |
| avatar | string | Avatar URL |
| phone_number | string | Phone number |
| language_preference | string | UI language |
| time_zone | string | User timezone |

### UserAuthMethod
วิธีการ authentication ของผู้ใช้

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary Key |
| user_id | uint | FK → users.id |
| auth_type | string | "password", "oauth", "sso" |
| auth_provider | string | "google", "github", "microsoft" |
| provider_id | string | ID from OAuth provider |
| is_primary | bool | Is primary auth method |

## Organization Management

### Organization
องค์กรหรือบริษัท (Multi-tenancy)

| Field | Type | Description | FK/Relationship |
|-------|------|-------------|-----------------|
| id | uint | Primary Key | - |
| name | string | Organization name | - |
| slug | string | URL slug (unique) | - |
| description | string | Description | - |
| logo_url | string | Logo URL | - |
| plan_type | string | "free", "pro", "enterprise" | - |
| status_id | uint | Organization status | FK → organization_statuses.id |
| settings | string | JSON settings | - |

**Relationships:**
- `Status` → OrganizationStatus
- `Members` → []OrganizationMember (1:many)
- `Projects` → []Project (1:many)

### OrganizationMember
สมาชิกในองค์กร

| Field | Type | Description | FK/Relationship |
|-------|------|-------------|-----------------|
| id | uint | Primary Key | - |
| organization_id | uint | FK → organizations.id | - |
| user_id | uint | FK → users.id | - |
| role | string | "owner", "admin", "member", "guest" | - |
| status_id | uint | Member status | FK → member_statuses.id |
| invited_at | *time.Time | Invitation timestamp | - |
| joined_at | *time.Time | Join timestamp | - |
| invited_by | *uint | FK → users.id | - |

## Project Management

### Project
โปรเจคในระบบ

| Field | Type | Description | FK/Relationship |
|-------|------|-------------|-----------------|
| id | uint | Primary Key | - |
| organization_id | *uint | FK → organizations.id (nullable for personal projects) | - |
| name | string | Project name | - |
| description | string | Project description | - |
| key | string | Project key (unique) | e.g., "PROJ" |
| owner_id | uint | FK → users.id | - |
| status_id | uint | Project status | FK → project_statuses.id |

**Relationships:**
- `Status` → ProjectStatus
- `Organization` → Organization (nullable)
- `Owner` → User
- `Members` → []ProjectMember (1:many)
- `Tickets` → []Ticket (1:many)
- `Statuses` → []TicketStatus (1:many)
- `Labels` → []Label (1:many)
- `Sprints` → []Sprint (1:many)

### ProjectMember
สมาชิกในโปรเจค

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary Key |
| project_id | uint | FK → projects.id |
| user_id | uint | FK → users.id |
| role | string | "admin", "developer", "viewer" |
| joined_at | time.Time | Join timestamp |

## Ticket System

### Epic
Epic หรือ feature ขนาดใหญ่

| Field | Type | Description | FK/Relationship |
|-------|------|-------------|-----------------|
| id | uint | Primary Key | - |
| project_id | uint | FK → projects.id | - |
| title | string | Epic title | - |
| description | string | Epic description | - |
| epic_key | string | Epic key (unique) | e.g., "PROJ-E1" |
| status_id | uint | Epic status | FK → epic_statuses.id |
| priority_id | uint | Epic priority | FK → priorities.id |
| owner_id | uint | FK → users.id | - |
| start_date | *time.Time | Start date | - |
| target_date | *time.Time | Target completion date | - |

**Relationships:**
- `Status` → EpicStatus
- `Priority` → Priority
- `Project` → Project
- `Owner` → User
- `Tickets` → []Ticket (1:many)
- `Labels` → []Label (many:many)

### Ticket
Ticket หรือ Issue ในระบบ

| Field | Type | Description | FK/Relationship |
|-------|------|-------------|-----------------|
| id | uint | Primary Key | - |
| project_id | uint | FK → projects.id | - |
| epic_id | *uint | FK → epics.id (nullable) | - |
| title | string | Ticket title | - |
| description | string | Ticket description | - |
| ticket_key | string | Ticket key (unique) | e.g., "PROJ-123" |
| type_id | uint | Ticket type | FK → ticket_types.id |
| status_id | uint | Ticket status | FK → ticket_statuses.id |
| priority_id | uint | Ticket priority | FK → priorities.id |
| assignee_id | *uint | FK → users.id (nullable) | - |
| reporter_id | uint | FK → users.id | - |
| parent_id | *uint | FK → tickets.id (for sub-tasks) | - |
| estimated_hours | *float64 | Estimated hours | - |
| actual_hours | *float64 | Actual hours worked | - |
| due_date | *time.Time | Due date | - |

**Relationships:**
- `Project` → Project
- `Epic` → Epic (nullable)
- `Type` → TicketType
- `Status` → TicketStatus
- `Priority` → Priority
- `Assignee` → User (nullable)
- `Reporter` → User
- `Parent` → Ticket (nullable)
- `Children` → []Ticket (1:many)
- `Comments` → []TicketComment (1:many)
- `Attachments` → []TicketAttachment (1:many)
- `Labels` → []Label (many:many)
- `Watchers` → []User (many:many)
- `TimeLogs` → []TimeLog (1:many)

### TicketStatus
สถานะของ Ticket แต่ละโปรเจค (customizable per project)

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary Key |
| project_id | uint | FK → projects.id |
| name | string | Status name |
| position | int | Display order |
| is_default | bool | Is default status for new tickets |
| color | string | Hex color for UI |

**Default Statuses per Project:**
- To Do (position: 1, is_default: true)
- In Progress (position: 2)
- In Review (position: 3)
- Done (position: 4)

### TicketComment
ความเห็นใน Ticket

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary Key |
| ticket_id | uint | FK → tickets.id |
| user_id | uint | FK → users.id |
| content | string | Comment content |

### TicketAttachment
ไฟล์แนบใน Ticket

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary Key |
| ticket_id | uint | FK → tickets.id |
| filename | string | Original filename |
| file_path | string | Stored file path |
| file_size | int64 | File size in bytes |
| mime_type | string | MIME type |
| uploaded_by | uint | FK → users.id |

### Label
Label หรือ Tag ของ Ticket

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary Key |
| project_id | uint | FK → projects.id |
| name | string | Label name |
| color | string | Hex color |
| description | string | Label description |

### TimeLog
การบันทึกเวลาทำงานใน Ticket

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary Key |
| ticket_id | uint | FK → tickets.id |
| user_id | uint | FK → users.id |
| hours | float64 | Hours worked |
| description | string | Work description |
| logged_date | time.Time | Date of work |

### Sprint
Sprint สำหรับ Agile Development

| Field | Type | Description | FK/Relationship |
|-------|------|-------------|-----------------|
| id | uint | Primary Key | - |
| project_id | uint | FK → projects.id | - |
| name | string | Sprint name | - |
| start_date | *time.Time | Sprint start date | - |
| end_date | *time.Time | Sprint end date | - |
| status_id | uint | Sprint status | FK → sprint_statuses.id |

**Relationships:**
- `Status` → SprintStatus
- `Project` → Project
- `Tickets` → []Ticket (many:many through sprint_tickets)

## Migration จาก String Status เป็น FK

ระบบได้ปรับปรุงจาก string-based status เป็น FK-based status เพื่อ:

1. **ความยืดหยุ่น**: สามารถเพิ่ม/แก้ไข status ได้โดยไม่ต้องแก้ไข code
2. **Consistency**: ป้องกันการพิมพ์ผิดและ typo
3. **Performance**: การ query และ index ที่ดีขึ้น
4. **UI/UX**: สามารถกำหนด color, icon, description สำหรับแต่ละ status
5. **Ordering**: สามารถกำหนดลำดับการแสดงผล

### การ Migration

เมื่อรันการ migration ระบบจะ:

1. สร้าง status tables ใหม่
2. Seed ข้อมูล default statuses
3. เพิ่ม `status_id` columns ให้ tables ที่เกี่ยวข้อง
4. Migrate ข้อมูลจาก string status ไป FK status
5. ลบ string status columns เก่า

### การใช้งาน

```go
// หา tickets ที่มี status "In Progress"  
tickets := []Ticket{}
db.Joins("Status").Where("statuses.name = ?", "in_progress").Find(&tickets)

// หา users ที่ active
users := []User{}
db.Joins("Status").Where("user_statuses.name = ?", "active").Find(&users)

// เปลี่ยน status ของ ticket
ticket.StatusID = 2 // In Progress
db.Save(&ticket)
```
