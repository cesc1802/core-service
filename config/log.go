package config

type LogConfig struct {
	Level                 string `yaml:"level"`
	ConsoleLoggingEnabled bool   `yaml:"consoleLoggingEnabled"`
	EncodeLogsAsJson      bool   `yaml:"encodeLogsAsJson"`
	FileLoggingEnabled    bool   `yaml:"fileLoggingEnabled"`
	Directory             string `yaml:"directory"`
	Filename              string `yaml:"filename"`
	MaxSize               int    `yaml:"maxSize"`
	MaxBackups            int    `yaml:"maxBackups"`
	MaxAge                int    `yaml:"maxAge"`
}
