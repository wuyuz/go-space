package db

import (
	"context"
	"database/sql"
	"testing"
)

func TestCreateUser(t *testing.T) {

	arg := CreateUserParams{
		Username: "wang",
		Pwd:      "xdsfw2312!@#",
		Email: sql.NullString{
			String: "sss",
			Valid:  true,
		},
	}

	u, err := testQuery.CreateUser(context.Background(), arg)
	if err != nil {
		t.Log("error: ", err)
	}
	t.Log("successful: ", u)

}
