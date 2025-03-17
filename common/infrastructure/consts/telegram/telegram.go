package telegram

type TokenChatId struct {
	Token  string
	ChatId int64
}

// "6309738216:AAEHBiMgW1sSiEOd-8hj_JH30td4a9xbjGg" "-1001990358319" 水雞小姐姐:雷速伊諾同學會

var EnvTokenChatIdMap = map[string]TokenChatId{
	"cqa":  {Token: "6989971521:AAEMKgUYhiI3dh", ChatId: -4093476937},
	"uat":  {Token: "6920110912:AAERHMRZ16wrP", ChatId: -4010650094},
	"prod": {Token: "6489537481:AAHesjYgb", ChatId: -4049998543},
}
