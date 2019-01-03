package log

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"sort"
)

/**

@author jingsong.zhu
@version 2018/11/5 下午2:15
*/
const (
	StdErrLogOutput = "stderr"
	StdOutLogOutput = "stdout"
)

type Logger struct {
	// logger logs server-side operations. The default is nil,
	// and "setupLogging" must be called before starting server.
	// Do not set logger directly.
	logger       *zap.Logger
	loggerConfig *zap.Config

	Conf Config `json:"log"`
}

type Config struct {
	// LogOutputs is either:
	//  - "default" as os.Stderr,
	//  - "stderr" as os.Stderr,
	//  - "stdout" as os.Stdout,
	//  - file path to append server logs to.
	// It can be multiple when "Logger" is zap.
	Outputs []string `json:"outputs"`
	// Debug is true, to enable debug level logging.
	Debug bool `json:"debug"`
}

var lg = NewDefault()

// NewDefault creates default logger.
func NewDefault() *Logger {
	var l = new(Logger)
	l.Conf = Config{
		Outputs: []string{"stdout"},
		Debug:   false,
	}

	// setupLogging
	if err := l.setupLogging(); err != nil {
		log.Fatalf("failed to setup logging: %v", err)
	}
	return l
}

// GetLogger returns the logger.
func GetLogger() *zap.Logger {
	return lg.logger
}

// setupLogging initializes the logging.
// Must be called after finishing configuring Config.
func (l *Logger) setupLogging() error {
	conf := l.Conf
	if len(conf.Outputs) == 0 {
		conf.Outputs = []string{StdOutLogOutput}
	}

	// use zapcore to support more features?
	lcfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:      "json",
		EncoderConfig: zap.NewProductionEncoderConfig(),

		OutputPaths:      make([]string, 0),
		ErrorOutputPaths: make([]string, 0),
	}

	outputPaths, errOutputPaths := make(map[string]struct{}), make(map[string]struct{})
	for _, v := range conf.Outputs {
		switch v {
		case StdErrLogOutput:
			outputPaths[StdErrLogOutput] = struct{}{}
			errOutputPaths[StdErrLogOutput] = struct{}{}

		case StdOutLogOutput:
			outputPaths[StdOutLogOutput] = struct{}{}
			errOutputPaths[StdOutLogOutput] = struct{}{}

		default:
			outputPaths[v] = struct{}{}
			errOutputPaths[v] = struct{}{}
		}
	}

	for v := range outputPaths {
		lcfg.OutputPaths = append(lcfg.OutputPaths, v)
	}
	for v := range errOutputPaths {
		lcfg.ErrorOutputPaths = append(lcfg.ErrorOutputPaths, v)
	}
	sort.Strings(lcfg.OutputPaths)
	sort.Strings(lcfg.ErrorOutputPaths)

	if conf.Debug {
		lcfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		grpc.EnableTracing = true
	}

	var err error
	l.logger, err = lcfg.Build()
	if err != nil {
		return err
	}

	l.loggerConfig = &lcfg

	if err != nil {
		return err
	}

	l.logger.Info("success to init logger of zap")
	return nil
}
