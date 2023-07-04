package id

import "os"

type Role struct {
}

func RoleID() Role {
	if os.Getenv("ENV") == "dev" {
		return Role{}
	} else {
		return Role{}
	}
}
