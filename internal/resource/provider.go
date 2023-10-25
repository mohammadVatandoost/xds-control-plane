package resource

type Resource struct {
	Name    string
	Version string
	Path    string
	Type    string
	Key		string
	Watchers map[string]struct{}
}

func NewResource(name, version, path, resourceType, key string) *Resource {
	return &Resource{
		Name:    name,
		Version: version,
		Path:    path,
		Type:    resourceType,
		Key:	key,
		Watchers: make(map[string]struct{}),
	}
}