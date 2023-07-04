package id

import "os"

type User struct {
	TOTSUMARU string
}

func UserID() User {
	if os.Getenv("ENV") == "dev" {
		return User{
			TOTSUMARU: "960104306151948328",
		}
	} else {
		return User{
			TOTSUMARU: "960104306151948328",
		}
	}
}
