package id

import "os"

type User struct {
	TOTSUMARU string
	THIS_BOT  string
}

func UserID() User {
	if os.Getenv("ENV") == "dev" {
		return User{
			TOTSUMARU: "960104306151948328",
			THIS_BOT:  "1125712152968314962",
		}
	} else {
		return User{
			TOTSUMARU: "960104306151948328",
			THIS_BOT:  "1125712152968314962",
		}
	}
}
