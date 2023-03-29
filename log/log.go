package log

type logCallback func(Level LevelLog, log interface{}) error

var Logger logCallback
var Level LevelLog

func Init(l LevelLog, f logCallback) {
	Level = l
	Logger = f
}

func TestLog(logText interface{}) {
	if InfoLevel >= Level {
		//设置日志等级，满足条件才打印
		Logger(InfoLevel, logText)
	}
}
