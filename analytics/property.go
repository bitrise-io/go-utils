package analytics

// Property ...
type Property interface {
	GetKey() string
	GetValue() interface{}
}

type stringProperty struct {
	key   string
	value string
}

// StringProperty ...
func StringProperty(key string, value string) Property {
	return stringProperty{key: key, value: value}
}

// GetKey ...
func (s stringProperty) GetKey() string {
	return s.key
}

// GetValue ...
func (s stringProperty) GetValue() interface{} {
	return s.value
}

type intProperty struct {
	key   string
	value int
}

// IntProperty ...
func IntProperty(key string, value int) Property {
	return intProperty{key: key, value: value}
}

// GetKey ...
func (s intProperty) GetKey() string {
	return s.key
}

// GetValue ...
func (s intProperty) GetValue() interface{} {
	return s.value
}

type longProperty struct {
	key   string
	value int64
}

// LongProperty ...
func LongProperty(key string, value int64) Property {
	return longProperty{key: key, value: value}
}

// GetKey ...
func (s longProperty) GetKey() string {
	return s.key
}

// GetValue ...
func (s longProperty) GetValue() interface{} {
	return s.value
}

type floatProperty struct {
	key   string
	value float64
}

// FloatProperty ...
func FloatProperty(key string, value float64) Property {
	return floatProperty{key: key, value: value}
}

// GetKey ...
func (s floatProperty) GetKey() string {
	return s.key
}

// GetValue ...
func (s floatProperty) GetValue() interface{} {
	return s.value
}

type boolProperty struct {
	key   string
	value bool
}

// BoolProperty ...
func BoolProperty(key string, value bool) Property {
	return boolProperty{key: key, value: value}
}

// GetKey ...
func (s boolProperty) GetKey() string {
	return s.key
}

// GetValue ...
func (s boolProperty) GetValue() interface{} {
	return s.value
}

type nestedProperty struct {
	key   string
	value []Property
}

// NestedProperty ...
func NestedProperty(key string, value ...Property) Property {
	return nestedProperty{key: key, value: value}
}

// GetKey ...
func (s nestedProperty) GetKey() string {
	return s.key
}

// GetValue ...
func (s nestedProperty) GetValue() interface{} {
	result := map[string]interface{}{}
	for _, property := range s.value {
		result[property.GetKey()] = property.GetValue()
	}
	return result
}
