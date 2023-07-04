package id

import "os"

type Channel struct {
	ERR_LOG string
	TEST    string
}

func ChannelID() Channel {
	if os.Getenv("ENV") == "dev" {
		return Channel{
			ERR_LOG: "1125719815458406410",
			TEST:    "1125719801558478848",
		}
	} else {
		return Channel{
			ERR_LOG: "1125719815458406410",
			TEST:    "1125719801558478848",
		}
	}
}
