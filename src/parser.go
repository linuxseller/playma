package main
// TODO: write a parser

import (
    "strings"
    "errors"
    // "fmt"
)

func ParseKeyVal(data string) (map[string]string, error) {
    result := make(map[string]string)
    lines := strings.Split(data, "\n")
    for _, i := range lines {
        keyval := strings.Split(i, "=")
        if len(keyval)==1 {continue}
        if i[0] == ';' {continue}
        if len(keyval) != 2 {
            return nil, errors.New("Could not parse provided string")
        }
        key, val := strings.TrimSpace(keyval[0]), strings.TrimSpace(keyval[1])
        result[key] = val
    }
    return result, nil
}
