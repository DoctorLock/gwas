package handlers

import "fmt"

type DataType string

const (
	String DataType = "string"
	Int    DataType = "int"
	Bool   DataType = "bool"
)

type RequestVar struct {
	Required  bool
	IsArray   bool
	Type      DataType
	MaxLength int
	Value     interface{} // final parsed value
}

func (rv *RequestVar) Get() (interface{}, error) {
	switch rv.Type {
	case String:
		v, ok := rv.Value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", rv.Value)
		}
		return v, nil

	case Int:
		v, ok := rv.Value.(int)
		if !ok {
			return nil, fmt.Errorf("expected int, got %T", rv.Value)
		}
		return v, nil

	case Bool:
		v, ok := rv.Value.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool, got %T", rv.Value)
		}
		return v, nil

	default:
		return nil, fmt.Errorf("unsupported datatype: %s", rv.Type)
	}
}

func (rv *RequestVar) GetString() (string, error) {
	v, err := rv.Get()
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func (rv *RequestVar) GetInt() (int, error) {
	v, err := rv.Get()
	if err != nil {
		return 0, err
	}
	return v.(int), nil
}

func (rv *RequestVar) GetBool() (bool, error) {
	v, err := rv.Get()
	if err != nil {
		return false, err
	}
	return v.(bool), nil
}
