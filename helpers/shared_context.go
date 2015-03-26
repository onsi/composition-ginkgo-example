package helpers

type SharedContext struct {
	Client *KeyValueStoreClient
	Prefix string
}

func NewSharedContext(keyValueStoreURL string, prefix string) SharedContext {
	return SharedContext{
		Client: &KeyValueStoreClient{URL: keyValueStoreURL},
		Prefix: prefix,
	}
}

func (s SharedContext) PrefixedKey(name string) string {
	return s.Prefix + "-" + name
}
