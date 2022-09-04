package server

import (
	"fmt"
	"strings"
)

type Path []string

func NewPath(path string) Path {
	return strings.Split(path, "/")
}

func (reqPath Path) RoutePath(routerPath Path) (map[string]string, error) {
	pathParamKey := ""
	pathParam := map[string]string{}

	if len(reqPath) != len(routerPath) {
		return pathParam, ErrInvaildURLPath
	}

	for i := range reqPath {
		if strings.Compare(reqPath[i], routerPath[i]) == 0 {
			continue
		}

		if _, err := fmt.Sscanf(routerPath[i], ":%s", &pathParamKey); err != nil {
			return pathParam, ErrInvaildURLPath
		}

		pathParam[pathParamKey] = reqPath[i]
	}
	return pathParam, nil
}
