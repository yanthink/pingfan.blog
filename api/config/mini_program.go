package config

type miniProgramWechat struct {
	Appid  string
	Secret string
}

type miniProgram struct {
	Wechat miniProgramWechat
}

var MiniProgram miniProgram

func loadMiniProgramConfig() {
	MiniProgram = miniProgram{
		Wechat: miniProgramWechat{
			Appid:  GetString("mini_program.wechat.appid"),
			Secret: GetString("mini_program.wechat.secret"),
		},
	}
}
