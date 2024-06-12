package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

type DistributedDriver struct {
	nodes   []string
	mutexes map[string]*sync.Mutex
	dir     string
}

func NewDistributedDriver(nodes []string) (*DistributedDriver, error) {
	return &DistributedDriver{
		nodes:   nodes,
		mutexes: make(map[string]*sync.Mutex),
		dir:     "./dbase",
	}, nil
}

func (d *DistributedDriver) getOrCreateMutex(collection string) *sync.Mutex {
	mutex, ok := d.mutexes[collection]
	if !ok {
		mutex = &sync.Mutex{}
		d.mutexes[collection] = mutex
	}
	return mutex
}

// Shard data across nodes
func (d *DistributedDriver) getShard(resource string) string {
	rand.Seed(time.Now().UnixNano())
	return d.nodes[rand.Intn(len(d.nodes))]
}

func (d *DistributedDriver) Write(collection, resource string, v interface{}) error {
	node := d.getShard(resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection, node)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, resource+".json")
	tempPath := filePath + ".tmp"

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))
	if err := ioutil.WriteFile(tempPath, b, 0644); err != nil {
		return err
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		return err
	}

	// Replication logic
	for _, replica := range d.nodes {
		if replica != node {
			replicaDir := filepath.Join(d.dir, collection, replica)
			if err := os.MkdirAll(replicaDir, 0755); err != nil {
				return err
			}
			replicaFilePath := filepath.Join(replicaDir, resource+".json")
			if err := ioutil.WriteFile(replicaFilePath, b, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *DistributedDriver) Read(collection, resource string, v interface{}) error {
	for _, node := range d.nodes {
		filePath := filepath.Join(d.dir, collection, node, resource+".json")
		if _, err := os.Stat(filePath); err == nil {
			b, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}
			return json.Unmarshal(b, v)
		}
	}
	return errors.New("resource not found")
}

func (d *DistributedDriver) ReadAll(collection string) ([]json.RawMessage, error) {
	var records []json.RawMessage
	for _, node := range d.nodes {
		dir := filepath.Join(d.dir, collection, node)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				continue
			}
			records = append(records, json.RawMessage(b))
		}
	}
	if len(records) == 0 {
		return nil, errors.New("no records found")
	}
	return records, nil
}

func (d *DistributedDriver) Delete(collection, resource string) error {
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	for _, node := range d.nodes {
		filePath := filepath.Join(d.dir, collection, node, resource+".json")
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DistributedDriver) DeleteAll(collection string) error {
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	for _, node := range d.nodes {
		dir := filepath.Join(d.dir, collection, node)
		if _, err := os.Stat(dir); err == nil {
			if err := os.RemoveAll(dir); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DistributedDriver) Search(query map[string]interface{}) (map[string][]string, error) {
	results := make(map[string][]string)
	for _, node := range d.nodes {
		err := filepath.Walk(filepath.Join(d.dir, node), func(path string, info os.FileInfo, err error) error {
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
		if err != nil {
			return nil, err
		}
	}
	return results, nil
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

func (d *DistributedDriver) RegexSearch(collection string, query map[string]string) ([]map[string]interface{}, error) {
	var records []map[string]interface{}
	for _, node := range d.nodes {
		dir := filepath.Join(d.dir, collection, node)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				continue
			}

			var record map[string]interface{}
			if err := json.Unmarshal(content, &record); err != nil {
				continue
			}

			if matchesRegex(record, query) {
				records = append(records, record)
			}
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
