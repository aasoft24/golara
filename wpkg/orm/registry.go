package orm

var modelRegistry []interface{}

func RegisterModel(m interface{}) {
	modelRegistry = append(modelRegistry, m)
}

func GetRegisteredModels() []interface{} {
	return modelRegistry
}
