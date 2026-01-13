package data_splitting

type config struct {
	Field         string `config:"field"`
	Target        string `config:"target"`
	MaxByteSize   int64  `config:"max_byte_size"`
	IgnoreFailure bool   `config:"ignore_failure"`
	//是否保留发送原始字段 默认不保留
	RetainOriginalField bool `config:"retain_original_field"`
}

var defaultConfig = config{
	Field:               "message",
	Target:              "messageChunks",
	RetainOriginalField: false,
	MaxByteSize:         1024 * 5, // 5KB
}
