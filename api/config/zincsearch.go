package config

type zincsearch struct {
	Url      string
	Index    string
	ShardNum int64
	UserID   string
	Password string
}

var Zincsearch zincsearch

func loadZincsearchConfig() {
	Zincsearch = zincsearch{
		Url:      GetString("zincsearch.url"),
		Index:    GetString("zincsearch.index"),
		ShardNum: GetInt64("zincsearch.shard_num"),
		UserID:   GetString("zincsearch.user_id"),
		Password: GetString("zincsearch.password"),
	}
}
