package registration

type Datasource interface {
	get() map[string]interface{}
	getAll() []map[string]interface{}
	find() map[string]interface{}
	findAll() map[string]interface{}
}
