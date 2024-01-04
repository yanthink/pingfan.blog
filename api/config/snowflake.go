package config

type snowflake struct {
	Epoch int64
	Node  int64
}

var Snowflake snowflake

func loadSnowflakeConfig() {
	Snowflake = snowflake{
		Epoch: GetInt64("snowflake.epoch", 1689292800000),
		Node:  GetInt64("snowflake.node", 0),
	}
}
