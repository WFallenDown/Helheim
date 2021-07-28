package Helheim

//MySQLRecord 一条日志
type MySQLRecord struct {
	Record
}

func (data *MySQLRecord) insert(fileName string, line int, message interface{}) {

}

func (data *MySQLRecord) GetLog(record *RecordList) error {
	return nil
}
