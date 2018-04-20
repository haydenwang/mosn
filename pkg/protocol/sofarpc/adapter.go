package sofarpc

import (
	"context"
	"reflect"
	"strconv"
	"sync"

	"gitlab.alipay-inc.com/afe/mosn/pkg/types"
)

func SofaPropertyHeader(name string) string {
	return name
}

func GetPropertyValue(properHeaders map[string]reflect.Kind, headers map[string]string, name string) interface{} {

	propertyHeaderName := SofaPropertyHeader(name)

	if value, ok := headers[propertyHeaderName]; ok {
		delete(headers, propertyHeaderName)

		return ConvertPropertyValue(value, properHeaders[name])
	} else {
		if value, ok := headers[name]; ok {

			return ConvertPropertyValue(value, properHeaders[name])
		}
	}

	return nil
}

func ConvertPropertyValue(strValue string, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Uint8:
		value, _ := strconv.ParseUint(strValue, 10, 8)
		return byte(value)
	case reflect.Uint16:
		value, _ := strconv.ParseUint(strValue, 10, 8)
		return uint16(value)
	case reflect.Uint32:
		value, _ := strconv.ParseUint(strValue, 10, 32)
		return uint32(value)
	case reflect.Uint64:
		value, _ := strconv.ParseUint(strValue, 10, 64)
		return uint64(value)
	case reflect.Int8:
		value, _ := strconv.ParseInt(strValue, 10, 8)
		return int8(value)
	case reflect.Int16:
		value, _ := strconv.ParseInt(strValue, 10, 16)
		return int16(value)
	case reflect.Int:
		value, _ := strconv.ParseInt(strValue, 10, 32)
		return int(value)
	case reflect.Int64:
		value, _ := strconv.ParseInt(strValue, 10, 64)
		return int64(value)
	default:
		return strValue
	}
}

func IsSofaRequest(headers map[string]string) bool {
	procode := ConvertPropertyValue(headers[SofaPropertyHeader(HeaderProtocolCode)], reflect.Uint8)

	if procode == PROTOCOL_CODE_V1 || procode == PROTOCOL_CODE_V2 {
		cmdtype := ConvertPropertyValue(headers[SofaPropertyHeader(HeaderCmdType)], reflect.Uint8)

		if cmdtype == REQUEST {
			return true
		}
	} else if procode == TR_PROTOCOL_CODE {
		requestFlage := ConvertPropertyValue(headers[SofaPropertyHeader(HeaderReqFlag)], reflect.Uint8)

		if requestFlage == HEADER_REQUEST {
			return true
		}
	}

	return false
}

func HasCodecException(headers map[string]string) bool {
	if v, ok := headers[types.HeaderException]; ok && v == types.MosnExceptionCodeC {
		return true
	}

	return false
}

func GetMap(context context.Context, defaultSize int) map[string]string {
	var amap map[string]string

	if context != nil && context.Value(types.ContextKeyConnectionCodecBufferPool) != nil {
		pool := context.Value(types.ContextKeyConnectionCodecBufferPool).(sync.Pool)
		amap = pool.Get().(map[string]string)
	}

	if amap == nil {
		amap = make(map[string]string, defaultSize)
	}

	return amap
}