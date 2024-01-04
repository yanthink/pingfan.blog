package storage

import (
	"blog/config"
	"fmt"
	"sync"
)

type Manager struct {
	mutex sync.RWMutex
	disks map[string]*Storage
}

func NewManager() *Manager {
	return &Manager{
		disks: map[string]*Storage{},
	}
}

func (m *Manager) Disk(name string) (*Storage, error) {
	m.mutex.RLock()

	storage, ok := m.disks[name]
	m.mutex.RUnlock()

	if ok {
		return storage, nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if diskConfig, ok := config.Storage.Disks[name]; ok {
		switch diskConfig.Driver {
		case "local":
			return m.createLocalDriver(diskConfig.Options), nil
		case "oss":
			return m.createOssDriver(diskConfig.Options), nil
		}
	}

	return nil, fmt.Errorf("driver [%s] is not supported", name)
}

func (*Manager) createLocalDriver(options map[string]any) *Storage {
	adapter := NewLocalAdapter(
		options["root"].(string),
		options["baseUrl"].(string),
	)

	return &Storage{Adapter: adapter}
}

func (*Manager) createOssDriver(options map[string]any) *Storage {
	adapter := NewOssAdapter(
		options["accessKeyId"].(string),
		options["accessKeySecret"].(string),
		options["endpoint"].(string),
		options["bucket"].(string),
		options["domain"].(string),
		options["ssl"].(bool),
	)

	return &Storage{Adapter: adapter}
}
