package domain

import (
	"github.com/google/uuid"
)

var (
	RealUserId     = "8e94e3f7-8944-454b-ab6a-5ef208337e2c"
	FakeUserId     = "121df86a-d02d-4b69-b6aa-6463df162831"
	RealTodoId     = "b1c8f0d2-3c4e-4f5a-9b6d-7e8f9a0b1c2d"
	FakeTodoId     = "e687f0ab-6965-4631-9e89-1ce86986fdec"
	FakeTodoIdUuid = uuid.MustParse(FakeTodoId)
	MockToken      = "mockedToken"
	TestUser       = &User{
		Id:              uuid.MustParse(RealUserId),
		FullName:        "Test User",
		Email:           "user@user.com",
		Role:            "USER",
		Password:        "testpassword123",
		IsEmailVerified: false,
	}
	TestTodo = &Todo{
		Id:        uuid.MustParse(RealTodoId),
		UserId:    TestUser.Id,
		Title:     "Test Todo",
		Completed: false,
	}

	MockJWTTestKey = "d16fb74a2c2d3ab64d6247fdcb703d08f9f4dd86625420a3e23ea08c1deaad19"
)
