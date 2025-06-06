package data

type StoreInterface interface {
	AddTask(t *Task) error
	ListAllTasks() []*Task
	DeleteTask(id int) error
}

var GlobalStore StoreInterface

func SetStore(s StoreInterface) {
	GlobalStore = s
}
