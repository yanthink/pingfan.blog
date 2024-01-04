package config

import "fmt"

type storage struct {
	Disk              string
	ResourceStoreDisk string

	Disks map[string]diskConfig
}

type diskConfig struct {
	Driver  string
	Options map[string]any
}

var Storage storage

func loadStorageConfig() {
	Storage = storage{
		Disk: GetString("storage.disk", "public"),
		ResourceStoreDisk: GetString(
			"storage.resource_store_disk",
			GetString("storage.disk", "public"),
		),

		Disks: map[string]diskConfig{
			"public": {
				Driver: "local",
				Options: map[string]any{
					"root":    "storage/app/public",
					"baseUrl": fmt.Sprintf("%s/storage", GetString("app.url")),
				},
			},
			"oss": {
				Driver: "oss",
				Options: map[string]any{
					"accessKeyId":     GetString("storage.oss.access_key_id"),
					"accessKeySecret": GetString("storage.oss.access_key_secret"),
					"endpoint":        GetString("storage.oss.endpoint"),
					"bucket":          GetString("storage.oss.bucket"),
					"domain":          GetString("storage.oss.domain"),
					"ssl":             GetBool("storage.oss.ssl"),
				},
			},
		},
	}
}
