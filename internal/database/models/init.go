package models

type ModelList []interface{}

func All() ModelList {
	return ModelList{
		// Status and lookup tables (ต้องสร้างก่อน)
		&UserStatus{},
		&ProjectStatus{},
		&EpicStatus{},
		&SprintStatus{},
		&OrganizationStatus{},
		&MemberStatus{},
		&Priority{},
		&TicketType{},
		&TicketStatus{}, // สำหรับ per-project ticket statuses

		// User models
		&User{},
		&UserAuthMethod{},
		&UserPreference{},
		&UserLocation{},
		// Organization models
		&Organization{},
		&OrganizationMember{},

		// Project and Ticket models
		&Project{},
		&ProjectMember{},
		&Epic{},
		&Ticket{},
		&TicketComment{},
		&TicketAttachment{},
		&Label{},
		&TimeLog{},
		&Sprint{},
	}
}
