package Helheim

//SQLServerRecord 一条日志
type SQLServerRecord struct {
	Record
}

func (data *SQLServerRecord) insert(fileName string, line int, message interface{}) {

}

func (data *SQLServerRecord) GetLog(record *RecordList) error {
	return nil
}
