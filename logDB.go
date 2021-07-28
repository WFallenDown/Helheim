package Helheim

const cronLog = "cron_log"

//TimePoint 任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

//Record 一条日志
type Record struct {
	fileName  string      `bson:"fileName"`
	line      int         `bson:"line"`
	Err       interface{} `bson:"err"`
	TimePoint TimePoint   `bson:"timePoint"` // 执行时间点
}

type RecordList struct {
	Data  []Record
	Skip  int64
	Limit int64
}

type IRecord interface {
	insert(fileName string, line int, message interface{})
	GetLog(record *RecordList) error
}

func insertLog(fileName string, line int, message interface{}) {
	var record IRecord
	switch dbName {
	case Mongo:
		record = new(MongoRecord)
		record.insert(fileName, line, message)
		break

	}
}
