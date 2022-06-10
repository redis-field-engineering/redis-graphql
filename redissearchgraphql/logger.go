package redissearchgraphql

import (
	"encoding/json"

	"go.uber.org/zap"
)

func InitLogger(filepath string) *zap.SugaredLogger {
	rawJSON := []byte(`{
		"level": "info",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoding": "json",
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	cfg.OutputPaths = []string{filepath}
	cfg.ErrorOutputPaths = []string{filepath}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	return sugar
}
