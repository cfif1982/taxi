package logger

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// делаю логирование отдельным пакетом, чтобы иметь возможность сменить логгер и изменять настройки
// также можно в разных файлах хранить разные логи
// пока до конца не понял на сколько это нужно, но пусть будет

type Logger struct {
	logger *zap.Logger
}

// получаем логгер
func GetLogger() (*Logger, error) {
	// настраиваем цветной вывод в лог
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

	// это нужно добавить, если логер буферизован
	// в данном случае не буферизован, но привычка хорошая
	defer logger.Sync()

	return &Logger{
		logger: logger,
	}, nil
}

// выводим сообщение уровня INFO
func (l *Logger) Info(messages ...string) {

	i := 0
	var msg, zapKey, zapVal string
	zapcore := make([]zapcore.Field, 0) // срез для вызова zap logger.Info

	// пробегаемся по переданным параметрам
	for _, message := range messages {
		// первый параметр - заголовок сообщнеия
		// остальные - это пары значений для zap.String("key", "val"), которые используются в logger.Info
		if i == 0 {
			msg = message
		} else {
			// формируем zap.String("key", "val")
			// если i нечетное, то это key, если i четное, то это val
			if i%2 == 0 {
				zapVal = message

				// добавляем в слайс
				zapcore = append(zapcore, zap.String(zapKey, zapVal))

			} else {
				zapKey = message
			}

		}
		i++
	}

	// выводим сообщение INFO
	l.logger.Info(msg, zapcore...)

}

// QUESTION: в каких случаях вызвать fatal?
// выводим сообщение уровня FATAL
func (l *Logger) Fatal(message string) {
	l.logger.Fatal(message)
}
