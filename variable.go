package tst

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

type VariableType int32

const (
	String VariableType = iota
	Int
	Datetime
	Timestamp
	Uuid
	DatetimeNow
	DatetimeYesterday
	DatetimeTomorrow
	FakeVehicleId
	Customer
)

func (t VariableType) String() string {
	switch t {
	case String:
		return "string"
	case Uuid:
		return "uuid"
	case DatetimeNow:
		return "datetime-now"
	case DatetimeYesterday:
		return "datetime-yesterday"
	case DatetimeTomorrow:
		return "datetime-tomorrow"
	case FakeVehicleId:
		return "fake-vehicle-id"
	case Customer:
		return "customer"
	}

	panic("unknown variable type")

}

type VariableAction int32

const (
	Static VariableAction = iota
	Dynamic
	Data
)

func (t VariableAction) String() string {
	switch t {
	case Static:
		return "static"
	case Dynamic:
		return "dynamic"
	case Data:
		return "data"
	}

	panic("unknown variable action")
}

type VariableScope int32

const (
	Global VariableScope = iota
	Local
)

func (t VariableScope) String() string {
	switch t {
	case Global:
		return "global"
	case Local:
		return "local"
	}

	panic("unknown variable scope")
}

type Variable struct {
	Type   VariableType
	Action VariableAction
	Scope  VariableScope
}

var placeholders map[VariableAction]string = map[VariableAction]string{
	Dynamic: "$",
	Static:  "#",
	Data:    "@",
}

var container map[string]string = make(map[string]string)

func ParseByte(input []byte) []byte {
	return []byte(Parse(string(input)))
}

func Parse(input string) string {
	for _, placeholder := range placeholders {
		input = Replace(input, placeholder)
	}
	return input
}

func Replace(input string, placeholder string) string {
	regex := `\` + placeholder + `\w*\` + placeholder
	pattern := regexp.MustCompile(regex)
	matches := pattern.FindAllString(input, -1)
	for _, match := range matches {
		key := strings.Replace(match, placeholder, "", -1)
		if match != "" {
			variable := Variables[key]
			if variable.Action == Dynamic {
				value, alreadyGenerated := container[key]
				if !alreadyGenerated {
					value = Generate(variable.Type)
					container[key] = value
				}

				input = strings.ReplaceAll(input, match, value)
			} else {
				contextValue, existInContext := container[key]
				if existInContext {
					input = strings.ReplaceAll(input, match, contextValue)
				} else {
					panic(key + " does not exist in context")
				}
			}
		}
	}

	return input
}

func Generate(t VariableType) string {
	switch t.String() {
	case "fake-vehicle-id":
		return "1333337"
	case "uuid":
		id, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		return id.String()
	case "customer":
		return "62af3baa-a3d8-400a-8e94-cff18913b43b"
	case "datetime-now":
		return time.Now().Format("2006-01-02T15:04:05.000Z")
	case "datetime-yesterday":
		return time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05.000Z")
	case "datetime-tomorrow":
		return time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05.000Z")
	case "timestamp-now":
		return time.Now().Format("20060102150405")
	case "timestamp-yesterday":
		return time.Now().AddDate(0, 0, -1).Format("20060102150405")
	}

	panic("unsupported variable type")
}

func EmptyContainer() {
	for k := range container {
		if Variables[k].Scope == Local {
			delete(container, k)
		}
	}
}

func SetParameter(key string, value string) {
	log.Printf("Set parameter %v:%v", key, value)
	container[key] = value
}

func SetIntParameter(key string, value int) {
	container[key] = strconv.Itoa(value)
}

func GetParameter(key string) string {
	return container[key]
}

func GetIntParameter(key string) int {
	value, err := strconv.Atoi(container[key])
	if err != nil {
		panic(err)
	}
	return value
}
