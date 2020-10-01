package go_restful_routes

import (
	"strconv"
	"strings"
)

type pathParams map[string]interface{}

func newPathParams(patternBlocks []string, pathBlocks []string) (pathParams, bool) {
	params := make(pathParams)
	for i, patternBlock := range patternBlocks {
		if ok := params.match(patternBlock, pathBlocks[i]); !ok {
			return nil, false
		}
	}
	return params, true
}

// users <=> users
// "{int:user_id}"  <=>  "123"
func (params pathParams) match(patternBlock string, pathBlock string) bool {
	if strings.HasPrefix(patternBlock, "{") &&
		strings.Contains(patternBlock, ":") &&
		strings.HasSuffix(patternBlock, "}") {
		return params.matchPattern(patternBlock, pathBlock)
	}
	return patternBlock == pathBlock
}

func (params pathParams) matchPattern(patternBlock string, pathBlock string) bool {
	splits := strings.Split(patternBlock[1:len(patternBlock)-1], ":") // trim `{` and `}`, split by `:`
	if len(splits) != 2 {
		return false
	}
	t, k := strings.ToLower(splits[0]), splits[1] // {t:k}
	switch t {
	case "string":
		params[k] = pathBlock
		return true
	case "int":
		if n, err := strconv.Atoi(pathBlock); err == nil {
			params[k] = int(n)
			return true
		}
	case "int8":
		if n, err := strconv.ParseInt(pathBlock, 10, 8); err == nil {
			params[k] = int8(n)
			return true
		}
	case "int16":
		if n, err := strconv.ParseInt(pathBlock, 10, 16); err == nil {
			params[k] = int16(n)
			return true
		}
	case "int32":
		if n, err := strconv.ParseInt(pathBlock, 10, 32); err == nil {
			params[k] = int32(n)
			return true
		}
	case "int64":
		if n, err := strconv.ParseInt(pathBlock, 10, 64); err == nil {
			params[k] = n
			return true
		}
	case "uint":
		if n, err := strconv.ParseUint(pathBlock, 10, 32); err == nil {
			params[k] = uint(n)
			return true
		}
	case "uint8":
		if n, err := strconv.ParseUint(pathBlock, 10, 8); err == nil {
			params[k] = uint8(n)
			return true
		}
	case "uint16":
		if n, err := strconv.ParseUint(pathBlock, 10, 16); err == nil {
			params[k] = uint16(n)
			return true
		}

	case "uint32":
		if n, err := strconv.ParseUint(pathBlock, 10, 32); err == nil {
			params[k] = uint32(n)
			return true
		}
	case "uint64":
		if n, err := strconv.ParseUint(pathBlock, 10, 64); err == nil {
			params[k] = n
			return true
		}
	case "float", "float32":
		if n, err := strconv.ParseFloat(pathBlock, 32); err == nil {
			params[k] = float32(n)
			return true
		}
	case "float64":
		if n, err := strconv.ParseFloat(pathBlock, 64); err == nil {
			params[k] = n
			return true
		}
	default:
		params[k] = pathBlock
		return true
	}
	return false
}
