package log

import "testing"

func Test_log(t *testing.T) {
	//SetOutputLogFile("./app.log")
	Debug("ttt%s", "ghhj")
	Info("ttt%s", "ghhj")
	Warn("ttt%s", "ghhj")
	Error("ttt%s", "ghhj")
	Fatal("ttt%s", "ghhj")
	Panic("ttt%s", "ghhj")
}
