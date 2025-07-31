package domain

import (
	"github.com/google/uuid"
)

var (
	RealUserId = "8e94e3f7-8944-454b-ab6a-5ef208337e2c"
	FakeUserId = "121df86a-d02d-4b69-b6aa-6463df162831"
	MockToken  = "mockedToken"
	TestUser   = &User{
		Id:              uuid.MustParse(RealUserId),
		FullName:        "Test User",
		Email:           "user@user.com",
		Role:            "USER",
		IsEmailVerified: true,
	}
)
