package main

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	//testExampleSugar()
	//testExample()
	//testDevelopment()
	//testDevelopmentSugar()
	//testConfig()
	//testJsonConfig()
	testDevelopProductConfig()
}

func testExample() {
	// example logger only output json format strings
	logger := zap.NewExample()
	logger = logger.Named("logger1").
		With(zap.String("lige", "shanghai"), zap.Int("year", 2019)).
		WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer func() { _ = logger.Sync() }()
	logger.Debug("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Info("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Warn("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Error("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	func() {
		defer func() { recover() }()
		logger.DPanic("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	}()
	logger.Fatal("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	fmt.Println("unreachable")
}

func testExampleSugar() {
	// SugaredLogger can output formatted log, like fmt.Sprintf("%s %d", ...)
	logger := zap.NewExample()
	logger = logger.Named("logger1").
		With(zap.String("lige", "shanghai"), zap.Int("year", 2019)).
		WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar := logger.Sugar()
	defer func() { _ = sugar.Sync() }()
	sugar.Debugf("my name is: %s %d", "van", 2019)
	sugar.Infof("my name is: %s %d", "van", 2019)
	sugar.Warnf("my name is: %s %d", "van", 2019)
	sugar.Errorf("my name is: %s %d", "van", 2019)
	func() {
		defer func() { recover() }()
		sugar.Panicf("my name is: %s %d", "van", 2019)
	}()
	sugar.Fatalf("my name is: %s %d", "van", 2019)
	fmt.Println("unreachable")
}

func testDevelopment() {
	// development logger can output structured log, like:
	// 2019-07-25T16:58:15.424+0800	DEBUG logger1 main/testZap.go:66 my name is: {"lige": "shanghai", "year": 2019, "name": "van", "fs": [2.5, 3.6]}
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("get new development err:", err)
		return
	}
	logger = logger.Named("logger1").
		With(zap.String("lige", "shanghai"), zap.Int("year", 2019)).
		WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)) // this take effect
	defer func() { _ = logger.Sync() }()
	logger.Debug("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Info("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Warn("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Error("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	func() {
		defer func() { recover() }()
		logger.DPanic("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	}()
	logger.Fatal("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	fmt.Println("unreachable")
}

func testDevelopmentSugar() {
	// development sugared logger can output structured log, like:
	// 2019-07-25T17:01:13.099+0800 DEBUG logger1 main/testZap.go:90 my name is: van 2019 {"lige": "shanghai", "year": 2019}
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("get new development err:", err)
		return
	}
	logger = logger.Named("logger1").
		With(zap.String("lige", "shanghai"), zap.Int("year", 2019)).
		WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)) // this take effect
	defer func() { _ = logger.Sync() }()
	sugar := logger.Sugar()
	sugar.Debugf("my name is: %s %d", "van", 2019)
	sugar.Infof("my name is: %s %d", "van", 2019)
	sugar.Warnf("my name is: %s %d", "van", 2019)
	sugar.Errorf("my name is: %s %d", "van", 2019)
	func() {
		defer func() { recover() }()
		sugar.Panicf("my name is: %s %d", "van", 2019)
	}()
	sugar.Fatalf("my name is: %s %d", "van", 2019)
	fmt.Println("unreachable")
}

func testConfigBuild() {
	dep := zap.NewDevelopmentConfig()
	dep.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // suppress debug
	dep.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	dep.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	dep.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	dep.OutputPaths = []string{"stdout", "/tmp/tmplog1"} // add output files here
	dep.ErrorOutputPaths = []string{"stderr"}
	printConfig(&dep)
	fmt.Println()

	logger, err := dep.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		fmt.Println("get new development err:", err)
		return
	}
	logger = logger.Named("logger1").
		With(zap.String("lige", "shanghai"), zap.Int("year", 2019))
	defer func() { _ = logger.Sync() }()

	logger.Debug("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Info("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Warn("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Error("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))

	sugar := logger.Sugar()
	sugar.Debugf("my name is: %s %d", "van", 2019)
	sugar.Infof("my name is: %s %d", "van", 2019)
	sugar.Warnf("my name is: %s %d", "van", 2019)
	sugar.Errorf("my name is: %s %d", "van", 2019)
}

func testJsonConfig() {
	rawJSON := []byte(`{` +
		`  "level": "info",                    ` + // log level
		`  "outputPaths": ["stdout"],          ` + // url to output logs
		`  "encoding": "console",                 ` + // output as json format or console line
		`  "encoderConfig": {                  ` +
		` 	  "messageKey": "M",               ` + // json keys
		` 	  "levelKey": "L",                 ` +
		` 	  "timeKey": "T",                  ` +
		` 	  "nameKey": "N",                  ` +
		` 	  "callerKey": "C",                ` +
		` 	  "levelEncoder": "capitalColor",  ` + // output format
		` 	  "timeEncoder": "iso8601",        ` +
		` 	  "callerEncoder": "short",        ` +
		` 	  "nameEncoder": "full"            ` +
		`}
	}`)
	// "encoding" is "json", output:
	// {"L":"\u001b[31mERROR\u001b[0m","T":"2019-07-26T15:46:54.795+0800","N":"logger1","C":"main/testZap.go:179","M":"my name is:","lige":"shanghai","year":2019,"name":"van","fs":[2.5,3.6]}
	// "encoding" is "console", output:
	// 2019-07-26T15:48:14.859+0800 INFO logger1 main/testZap.go:179 my name is: {"lige": "shanghai", "year": 2019, "name": "van", "fs": [2.5, 3.6]}
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		fmt.Println("unmarshal err:", err)
		return
	}

	logger, err := cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		fmt.Println("config build err:", err)
		return
	}
	logger = logger.Named("logger1").
		With(zap.String("lige", "shanghai"), zap.Int("year", 2019))
	defer func() { _ = logger.Sync() }()

	logger.Debug("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Info("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Warn("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
	logger.Error("my name is:", zap.String("name", "van"), zap.Float64s("fs", []float64{2.5, 3.6}))
}

func testDevelopProductConfig() {
	// difference between default production and development config
	//            dep         pro
	// format:    console     json
	// sampling:  no          yes
	// time:      iso8601     epoch
	// production encoder config
	_ = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// production config
	_ = zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// development encoder config
	_ = zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// development config
	_ = zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func printConfig(config *zap.Config) {
	fmt.Println("Level            :", config.Level)
	fmt.Println("Development      :", config.Development)
	fmt.Println("DisableCaller    :", config.DisableCaller)
	fmt.Println("DisableStacktrace:", config.DisableStacktrace)
	fmt.Println("Sampling         :", config.Sampling)
	fmt.Println("Encoding         :", config.Encoding)
	fmt.Println("EncoderConfig    :", config.EncoderConfig)
	fmt.Println("OutputPaths      :", config.OutputPaths)
	fmt.Println("ErrorOutputPaths :", config.ErrorOutputPaths)
	fmt.Println("InitialFields    :", config.InitialFields)
}
