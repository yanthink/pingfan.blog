package config

type database struct {
	Host string
	Port int
	User string
	Pass string
	Name string
	// 最大连接数
	// 决定最大连接数之前，需要考虑以下几个因素：
	// 1.系统硬件资源：系统的CPU、内存和磁盘等硬件资源越高，可以支持的最大连接数就越多。
	// 2.应用程序负载：如果应用程序不是很耗费资源，那么较低的最大连接数就足够了；而如果应用程序需要处理大量并发请求，那么需要增加最大连接数。
	// 3.MySQL服务器配置：MySQL服务器的配置也会影响最大连接数设置。例如，如果您使用InnoDB存储引擎，则可能需要更多的连接来处理事务。
	// 通常来说，建议将最大连接数设置为可用内存除以每个连接所需内存的结果。默认情况下，每个连接使用大约1MB的内存，因此如果您有8GB内存，最大连接数应该设置为8000。
	// 不过这只是一个基本指导，实际上还需要根据具体情况进行调整。当最大连接数设置得太高时，可能会导致系统出现性能问题，例如响应时间变慢、应用程序崩溃等。
	MaxOpenConnections int
	// 设置最大空闲连接数
	// 空闲连接是已经建立但未被使用的连接。在某些情况下，通过增加最大空闲连接数可以提高服务器性能，例如：
	// 对于具有高并发访问量的 Web 应用程序，如果希望减少客户端重新连接到服务器的开销，则可以适当增加最大空闲连接数。
	// 如果应用程序有较长的查询等待时间，并且可能存在同时运行多个查询的情况，则可以增加最大空闲连接数以提高响应速度。
	// 一般来说，最大空闲连接数的建议值为总连接数的25%至50%之间，但也要根据具体情况进行调整。如果服务器资源较为紧张，可以适当降低最大空闲连接数以避免资源浪费。
	MaxIdleConnections int
	// 链接的过期时间
	// 连接过期是指连接处于空闲状态（未被查询或更新）的时间超过了设定的最大值。如果一个连接过期，则 MySQL 会自动关闭该连接，并释放相应的资源，以便其他连接可以使用。
	// 一般来说，连接过期时间越长，系统的负担就越小，但同时也可能会占用更多的资源。如果您的应用程序有长时间的查询等待时间，则建议将连接过期时间设置得比较长，
	// 例如1小时、2小时，甚至更长时间。但如果您的应用程序通常只需要短时间内的数据访问，则可以将连接过期时间设为较短的时间，例如10分钟或30分钟。
	// 另外，如果您的应用程序需要频繁地连接数据库，则建议将连接过期时间设为较短的时间，这样可以减少服务器上持续存在的连接数，并提高系统资源利用率。
	MaxLifeSeconds int
	AutoMigrate    bool
}

var Database database

func loadDatabaseConfig() {
	Database = database{
		GetString("mysql.host", "127.0.0.1"),
		GetInt("mysql.port", 3306),
		GetString("mysql.user"),
		GetString("mysql.pass"),
		GetString("mysql.name"),
		GetInt("mysql.max_open_connections", 1000),
		GetInt("mysql.max_idle_connections", 250),
		GetInt("mysql.max_life_seconds", 1800),
		GetBool("mysql.auto_migrate", false),
	}
}
