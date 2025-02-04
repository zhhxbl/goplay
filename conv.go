package goplay

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	REINT, _   = regexp.Compile("^[0-9,]+$")
	REFLOAT, _ = regexp.Compile("^[0-9]+[.][0-9]+$")
)

func string2Val(str string) interface{} {
	tmp := []byte(str)
	switch {
	case str == "false":
		return false
	case str == "true":
		return true
	case REINT.Match(tmp):
		_int, _ := strconv.ParseInt(str, 10, 64)
		return _int
	case REFLOAT.Match(tmp):
		_float, _ := strconv.ParseFloat(str, 64)
		return _float
	}
	return str
}

func ParseBool(v interface{}) (bool, error) {
	if v != nil {
		switch v := v.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%d", v)
			if strV == "1" {
				return true, nil
			} else if strV == "0" {
				return false, nil
			}
		case float64:
			if v == 1.0 {
				return true, nil
			} else if v == 0.0 {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", v)
	}
	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}

func ParseInt(v interface{}) (int, error) {
	switch v.(type) {
	case int:
		return v.(int), nil
	case int64:
		return int(v.(int64)), nil
	case float32:
		return int(v.(float32)), nil
	case float64:
		return int(v.(float64)), nil
	case string:
		return strconv.Atoi(v.(string))
	}
	return 0, errors.New("can not convert " + reflect.TypeOf(v).String() + " to int")
}

func ParseInt64(v interface{}) (int64, error) {
	switch v.(type) {
	case int:
		return int64(v.(int)), nil
	case int8:
		return int64(v.(int8)), nil
	case int16:
		return int64(v.(int16)), nil
	case int32:
		return int64(v.(int32)), nil
	case int64:
		return v.(int64), nil
	case float32:
		return int64(v.(float32)), nil
	case float64:
		return int64(v.(float64)), nil
	case string:
		return strconv.ParseInt(v.(string), 10, 0)
	}
	return 0, errors.New("can not convert " + reflect.TypeOf(v).String() + " to int")
}

func ParseInt32(v interface{}) (int32, error) {
	switch v.(type) {
	case int:
		return int32(v.(int)), nil
	case int8:
		return int32(v.(int8)), nil
	case int16:
		return int32(v.(int16)), nil
	case int32:
		return v.(int32), nil
	case int64:
		return int32(v.(int64)), nil
	case float32:
		return int32(v.(float32)), nil
	case float64:
		return int32(v.(float64)), nil
	case string:
		val, _ := strconv.ParseInt(v.(string), 10, 0)
		return int32(val), nil
	}
	return 0, errors.New("can not convert " + reflect.TypeOf(v).String() + " to int")
}

func ParseInt8(v interface{}) (int8, error) {
	switch v.(type) {
	case int:
		return int8(v.(int)), nil
	case int8:
		return v.(int8), nil
	case int16:
		return int8(v.(int16)), nil
	case int32:
		return int8(v.(int32)), nil
	case int64:
		return int8(v.(int64)), nil
	case float32:
		return int8(v.(float32)), nil
	case float64:
		return int8(v.(float64)), nil
	case string:
		if tmp, err := strconv.ParseInt(v.(string), 10, 0); err != nil {
			return 0, err
		} else {
			return int8(tmp), nil
		}
	}
	return 0, errors.New("can not convert " + reflect.TypeOf(v).String() + " to int")
}

func ParseFloat64(v interface{}) (float64, error) {
	switch v.(type) {
	case int:
		return float64(v.(int)), nil
	case int8:
		return float64(v.(int8)), nil
	case int16:
		return float64(v.(int16)), nil
	case int32:
		return float64(v.(int32)), nil
	case int64:
		return float64(v.(int64)), nil
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return v.(float64), nil
	case string:
		return strconv.ParseFloat(v.(string), 64)
	}
	return 0, errors.New("can not convert " + reflect.TypeOf(v).String() + " to int")
}

func ParseFloat32(v interface{}) (float32, error) {
	switch v.(type) {
	case int:
		return float32(v.(int)), nil
	case int8:
		return float32(v.(int8)), nil
	case int16:
		return float32(v.(int16)), nil
	case int32:
		return float32(v.(int32)), nil
	case int64:
		return float32(v.(int64)), nil
	case float32:
		return v.(float32), nil
	case float64:
		return float32(v.(float64)), nil
	case string:
		if tmp, err := strconv.ParseFloat(v.(string), 32); err != nil {
			return 0, err
		} else {
			return float32(tmp), nil
		}
	}
	return 0, errors.New("can not convert " + reflect.TypeOf(v).String() + " to int")
}

func ParseString(v interface{}) (string, error) {
	switch v.(type) {
	case string:
		return v.(string), nil
	case int:
		return strconv.FormatInt(int64(v.(int)), 10), nil
	case int64:
		return strconv.FormatInt(v.(int64), 10), nil
	case float32:
		return strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64), nil
	}
	return "", errors.New("can not convert " + reflect.TypeOf(v).String() + " to string")
}

func ParseSliceString(v interface{}) (list []string, err error) {
	var str string
	switch v.(type) {
	case []string:
		return v.([]string), nil
	case []interface{}:
		for _, i := range v.([]interface{}) {
			if str, err = ParseString(i); err != nil {
				return nil, errors.New("can not conver " + reflect.TypeOf(i).String() + " to slice string")
			}
			list = append(list, str)
		}
	case []int:
		for _, i := range v.([]int) {
			list = append(list, strconv.Itoa(i))
		}
	case []int64:
		for _, i := range v.([]int64) {
			list = append(list, strconv.FormatInt(i, 10))
		}
	case []float32:
		for _, i := range v.([]float32) {
			list = append(list, strconv.FormatFloat(float64(i), 'f', -1, 32))
		}
	case []float64:
		for _, i := range v.([]float64) {
			list = append(list, strconv.FormatFloat(i, 'f', -1, 32))
		}
	}

	return list, errors.New("can not convert to slice string")
}

func ParseSliceInterface(v interface{}) (list []interface{}, err error) {
	switch v.(type) {
	case []interface{}:
		return v.([]interface{}), nil
	case []string:
		for _, i := range v.([]string) {
			list = append(list, i)
		}
	case []int:
		for _, i := range v.([]int) {
			list = append(list, i)
		}
	case []int32:
		for _, i := range v.([]int32) {
			list = append(list, i)
		}
	case []int64:
		for _, i := range v.([]int64) {
			list = append(list, i)
		}
	case []float64:
		for _, i := range v.([]float64) {
			list = append(list, i)
		}
	case []float32:
		for _, i := range v.([]float32) {
			list = append(list, i)
		}
	}
	return nil, errors.New("can not convert " + reflect.TypeOf(v).String() + " to slice interface{}")
}

func ParseMapInterface(v interface{}) (list map[string]interface{}, err error) {
	var ok bool
	if list, ok = v.(map[string]interface{}); !ok {
		err = errors.New("can not convert to map[string]interface {}")
	}
	return
}

func DecodeHost(driver, dest string) (username, password, host, database string) {
	uIndex := strings.Index(dest, ":")
	pIndex := strings.Index(dest, fmt.Sprintf("@%s(", driver))
	dIndex := strings.Index(dest, ")/")

	username = dest[:uIndex]
	password = dest[uIndex+1 : pIndex]
	host = dest[pIndex+2+len(driver) : dIndex]
	database = dest[dIndex+2:]
	return
}
