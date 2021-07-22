package log

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strings"
	"time"
)

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

func InsertLog(fileName string, line int, message interface{}) {
	ctx, cannel := context.WithTimeout(context.Background(), time.Minute)
	defer cannel()
	v := GetLogConnectString()
	path := strings.Join([]string{"mongodb://", v.UserName, ":", v.Password, "@", v.Ip, ":", v.Port}, "")
	// 建立mongodb连接
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(path))
	if err != nil {
		logError(Error, err)
		return
	}

	// 2, 选择数据库my_db
	database := client.Database(v.DbName)

	// 3, 选择表my_collection
	collection := database.Collection(cronLog)
	// 4, 插入记录(bson)
	record := &Record{
		fileName:  fileName,
		line:      line,
		Err:       message,
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	result, err := collection.InsertOne(context.TODO(), record)
	if err != nil {
		logError(Error, err)
		return
	}

	if result.InsertedID == nil {
		fmt.Println("Failed to insert log")
	}
}

func GetLog(record *RecordList) error {
	ctx, cannel := context.WithTimeout(context.Background(), time.Minute)
	defer cannel()
	v := GetLogConnectString()
	path := strings.Join([]string{"mongodb://", v.UserName, ":", v.Password, "@", v.Ip, ":", v.Port}, "")
	// 建立mongodb连接
	clientOptions := options.Client().ApplyURI(path)
	client, err := mongo.Connect(
		context.TODO(), clientOptions)
	if err != nil {
		logError(Error, err)
		return err
	}

	// 2, 选择数据库my_db
	database := client.Database(v.DbName)

	// 3, 选择表my_collection
	collection := database.Collection(cronLog)

	var data Record
	var findoptions = new(options.FindOptions)
	findoptions.SetLimit(record.Limit)
	findoptions.SetSkip(record.Skip)
	findoptions.SetSort(bsonx.Doc{{"StartTime", bsonx.Int32(-1)}})
	d, err := collection.Find(ctx, bson.M{}, findoptions)
	if err != nil {
		logError(Error, err)
		return err
	}
	for d.Next(context.Background()) {
		if err = d.Decode(&data); err != nil {
			logError(Error, err)
			return err
		}
		record.Data = append(record.Data, data)
	}
	return nil
}
