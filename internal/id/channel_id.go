package id

import "os"

type Channel struct {
	ERR_LOG string
}

func ChannelID() Channel {
	if os.Getenv("ENV") == "dev" {
		return Channel{
			ERR_LOG: "1075023820395663440",
		}
	} else {
		return Channel{
			ERR_LOG: "1075023820395663440",
		}
	}
}
