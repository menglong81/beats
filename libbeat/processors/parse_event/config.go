package parse_event

// Config for parse_event processor.
type Config struct {
	//Name           string          `config:"name"`
	Index             int    `config:"host"`
	fieldPath         string `config:"field_path"`
	separator         string `config:"separator"`
	ignoreError       bool   `config:"ignorer"`
	mode              string `config:"mode"`
	keyName           string `config:"keyframe"`
	enableEnvType     bool   `config:"enable_env_type"`
	enableTime        bool   `config:"enable_time"`
	deleteUnuseFields bool   `config:"delete_unused_fields"`
}

func defaultConfig() Config {
	return Config{
		Index: -1,
		mode: "auto",
		fieldPath: "source_path",
		separator: "/",
		ignoreError: false,
		keyName: "default_key",
		enableEnvType: false,
		enableTime: true,
		deleteUnuseFields: true,

	}
}

