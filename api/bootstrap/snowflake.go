package bootstrap

import (
	"blog/config"
	"github.com/bwmarrin/snowflake"
)

func SetupSnowflake() {
	snowflake.Epoch = config.Snowflake.Epoch
}
