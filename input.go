package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func decodeJSON(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	return nil
}

func readInputFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := decodeJSON(data, v); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

type flagSet struct {
	values map[string]string
	rest   []string
}

func parseFlags(args []string) (*flagSet, error) {
	fs := &flagSet{values: make(map[string]string)}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			key := strings.TrimLeft(arg, "-")
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag -%s needs a value", key)
			}
			fs.values[key] = args[i+1]
			i++
			continue
		}
		fs.rest = append(fs.rest, arg)
	}
	return fs, nil
}

func (fs *flagSet) Float(key string) (float64, error) {
	raw, ok := fs.values[key]
	if !ok {
		return 0, fmt.Errorf("missing -%s", key)
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("-%s value %q is not a number", key, raw)
	}
	return v, nil
}

func (fs *flagSet) Has(key string) bool {
	_, ok := fs.values[key]
	return ok
}

func (fs *flagSet) String(key string) string {
	return fs.values[key]
}

func requireOneInput(fs *flagSet, what string) (string, error) {
	if len(fs.rest) != 1 {
		return "", errors.New(what + " needs exactly one input path")
	}
	return fs.rest[0], nil
}
