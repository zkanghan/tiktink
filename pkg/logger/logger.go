package logger

import (
	"os"
	"tiktink/pkg/tracer"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger

func InitLogger() {
	var coreArr []zapcore.Core
	//zapcore.Core需要三个配置——Encoder，WriteSyncer，LogLevel

	//获取编码器encoder
	encoderConfig := zap.NewProductionEncoderConfig()            //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	encoder := zapcore.NewConsoleEncoder(encoderConfig)          //日志同时输出到控制台和文件

	//获取日志级别LogLevel
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	//info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info.log", //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    2,                //文件大小限制,单位MB
		MaxBackups: 100,              //最大保留日志文件数量
		MaxAge:     30,               //日志文件保留天数
		Compress:   false,            //是否压缩处理
	})
	//error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log", //日志文件存放目录
		MaxSize:    1,                 //文件大小限制,单位MB
		MaxBackups: 5,                 //最大保留日志文件数量
		MaxAge:     30,                //日志文件保留天数
		Compress:   false,             //是否压缩处理
	})
	//用获取到的Encoder，WriteSyncer，LogLevel组装Core
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)    //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	l = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //logger.AddCaller()为显示文件名和行号，可省略

}

func PrintLogWithCTX(errorMsg string, err error, ctx *tracer.TraceCtx) {
	errorMsg += "  "
	l.Error(errorMsg, zap.Error(err), zap.String("TracePath", ctx.ToString()))
}

func PrintLog(errorMsg string, err error) {
	errorMsg += "  "
	l.Error(errorMsg, zap.Error(err))
}
