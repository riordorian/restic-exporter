package log

import "go.uber.org/zap"

type ZapAdapter struct {
	Logger *zap.SugaredLogger
}

func (z *ZapAdapter) Info(msg string, args ...interface{}) {
	z.Logger.Infow(msg, args...)
}

func (z *ZapAdapter) Warn(msg string, args ...interface{}) {
	z.Logger.Warnw(msg, args...)
}

func (z *ZapAdapter) Error(msg string, args ...interface{}) {
	z.Logger.Errorw(msg, args...)
}

func (z *ZapAdapter) Sync() error {
	return z.Logger.Sync()
}

func NewZapAdapter() (*ZapAdapter, error) {
	logger, err := zap.NewDevelopment()

	if err != nil {
		return nil, err
	}

	sugared := logger.Sugar()

	return &ZapAdapter{
		Logger: sugared,
	}, nil
}
