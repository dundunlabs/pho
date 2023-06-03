package tra

import "strconv"

type Params map[string]string

func (p Params) Param(key string) string {
	return p[key]
}

func (p Params) ParamInt(key string) (int, error) {
	return strconv.Atoi(p.Param(key))
}
