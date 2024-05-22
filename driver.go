package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string
}

func NewDriver(dir string) (*Driver, error) {
	dir = filepath.Clean(dir)

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
	}

	switch _, err := os.Stat(dir); {
	case err == nil:
		return &driver, nil
	case os.IsNotExist(err):
		fmt.Println("Creating the database at :", dir)
		return &driver, os.MkdirAll(dir, 0777)
	default:
		return nil, err
	}
}

func (d *Driver) Write(dataset string, name string, s Student) error {
	if err := d.handleError(dataset, name); err != nil {
		return err
	}

	mutex := d.mutexSingleton(dataset)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, dataset)
	dst := filepath.Join(dir, strings.ToLower(name)+".json")

	switch _, err := os.Stat(dir); {
	case os.IsNotExist(err):
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	case err != nil:
		return err
	default:
	}

	s_json, err := json.MarshalIndent(s, "", "    ")

	if err != nil {
		return err
	}
	if err := os.WriteFile(dst, s_json, 0644); err != nil {
		return err
	}

	return nil
}

func (d *Driver) Read(dataset string, f string, s_obj Student) error {
	if err := d.handleError(dataset, f); err != nil {
		return err
	}

	s := filepath.Join(d.dir, dataset, f)

	if _, err := os.Stat(s + ".json"); err != nil {
		return err
	}

	b, err := os.ReadFile(s + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &s_obj)
}

func (d *Driver) ReadAll(dataset string) ([]string, error) {
	if dataset == "" {
		return nil, errors.New("invalid dataset")
	}

	dir := filepath.Join(d.dir, dataset)

	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}

	fs, _ := os.ReadDir(dir)

	var rprts []string
	for _, f := range fs {
		b, err := os.ReadFile(filepath.Join(dir, f.Name()))

		if err != nil {
			return nil, err
		}

		rprts = append(rprts, string(b))
	}

	return rprts, nil
}

func (d *Driver) Delete(dataset string, f string) error {
	p := filepath.Join(dataset, f)

	mutex := d.mutexSingleton(dataset)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, p)

	switch info, err := os.Stat(dir + ".json"); {
	case info == nil, err != nil:
		return fmt.Errorf("unable to find the file or dir named %v", p)
	case info.Mode().IsDir():
		return os.RemoveAll(dir)
	case info.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

func (d *Driver) mutexSingleton(dataset string) *sync.Mutex {
	m, ok := d.mutexes[dataset]
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[dataset] = m
	}

	return m
}

func (d *Driver) handleError(dataset string, f string) error {
	if dataset == "" {
		return fmt.Errorf("missing collection - no place to save record")
	}

	if f == "" {
		return fmt.Errorf("missing f - unable to save record")
	}
	return nil
}
