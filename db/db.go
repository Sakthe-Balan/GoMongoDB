package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/jcelliott/lumber"
)

type Logger interface {
	Fatal(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}

type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string
	log     Logger
}

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}
	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}
	opts.Logger.Debug("Creating the database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - no place to save record!")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save record (no name)!")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - unable to read!")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource - unable to read record (no name)!")
	}

	record := filepath.Join(d.dir, collection, resource+".json")

	if _, err := stat(record); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(record)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (d *Driver) ReadAll(collection string) ([]json.RawMessage, error) {
	if collection == "" {
		return nil, fmt.Errorf("Missing collection - unable to read")
	}
	dir := filepath.Join(d.dir, collection)

	d.log.Debug("Checking directory: %s", dir)
	if _, err := stat(dir); err != nil {
		d.log.Error("Directory check error: %s", err)
		return nil, err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		d.log.Error("Read directory error: %s", err)
		return nil, err
	}

	var records []json.RawMessage
	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			d.log.Error("Read file error: %s", err)
			return nil, err
		}
		records = append(records, json.RawMessage(b))
	}
	return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("Unable to find file or directory named %v\n", path)

	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func stat(path string) (os.FileInfo, error) {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.Stat(path + ".json")
	}
	return fi, err
}

func (d *Driver) DeleteAll(collection string) error {
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	d.log.Debug("Checking directory: %s", dir)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		d.log.Error("Unable to find directory: %s", err)
		return fmt.Errorf("Unable to find directory named %v\n", dir)
	case fi.Mode().IsDir():
		d.log.Debug("Deleting directory: %s", dir)
		return os.RemoveAll(dir)
	default:
		d.log.Error("Invalid file mode: %s", dir)
		return fmt.Errorf("Invalid file mode for %v\n", dir)
	}
}

func (d *Driver) Search(query map[string]interface{}) (map[string][]string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	results := make(map[string][]string)
	err := filepath.Walk(d.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var record map[string]interface{}
			if err := json.Unmarshal(content, &record); err != nil {
				return err
			}

			if matchesQuery(record, query) {
				collection := filepath.Base(filepath.Dir(path))
				resource := filepath.Base(path)
				resource = resource[:len(resource)-len(filepath.Ext(resource))]
				results[collection] = append(results[collection], resource)
			}
		}
		return nil
	})

	return results, err
}

func matchesQuery(record, query map[string]interface{}) bool {
	for key, value := range query {
		if recordValue, ok := record[key]; ok {
			switch value := value.(type) {
			case map[string]interface{}: // Handling nested queries like { "age": { "$gt": 25 } }
				for op, v := range value {
					switch op {
					case "$gt":
						if recordValue.(float64) <= v.(float64) {
							return false
						}
					case "$lt":
						if recordValue.(float64) >= v.(float64) {
							return false
						}
					case "$gte":
						if recordValue.(float64) < v.(float64) {
							return false
						}
					case "$lte":
						if recordValue.(float64) > v.(float64) {
							return false
						}
					case "$ne":
						if recordValue == v {
							return false
						}
					case "$in":
						found := false
						for _, item := range v.([]interface{}) {
							if recordValue == item {
								found = true
								break
							}
						}
						if !found {
							return false
						}
					}
				}
			default:
				if recordValue != value {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

func containsKeyword(content, keyword string) bool {
	return strings.Contains(content, keyword)
}

func (d *Driver) RegexSearch(collection string, query map[string]string) ([]map[string]interface{}, error) {
	if collection == "" {
		return nil, fmt.Errorf("Missing collection - unable to search")
	}

	dir := filepath.Join(d.dir, collection)
	d.log.Debug("Checking directory: %s", dir)
	if _, err := stat(dir); err != nil {
		d.log.Error("Directory check error: %s", err)
		return nil, err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		d.log.Error("Read directory error: %s", err)
		return nil, err
	}

	var records []map[string]interface{}
	for _, file := range files {
		content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			d.log.Error("Read file error: %s", err)
			continue
		}

		var record map[string]interface{}
		if err := json.Unmarshal(content, &record); err != nil {
			d.log.Error("Unmarshal error: %s", err)
			continue
		}

		if matchesRegex(record, query) {
			records = append(records, record)
		}
	}
	return records, nil
}

func matchesRegex(record map[string]interface{}, query map[string]string) bool {
	for key, pattern := range query {
		if value, ok := record[key]; ok {
			strValue := fmt.Sprintf("%v", value)
			matched, err := regexp.MatchString(pattern, strValue)
			if err != nil || !matched {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
