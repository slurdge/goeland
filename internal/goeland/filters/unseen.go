package filters

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/log"
	bolt "go.etcd.io/bbolt"
)

func openDatabase(config config.Provider) (*bolt.DB, error) {
	databaseName := config.GetString("database")
	if databaseName == "" {
		databaseName = "goeland.db"
		if ex, err := os.Executable(); err == nil {
			exPath := filepath.Dir(ex)
			databaseName = filepath.Join(exPath, databaseName)
		}
	}
	fmt.Println(databaseName)
	database, err := bolt.Open(databaseName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("error opening seen database %s: %v", databaseName, err)
	}
	return database, err
}

//PurgeUnseen will remove all the entries for this source
func PurgeUnseen(config config.Provider, sourceName string, numOfDays int) error {
	database, err := openDatabase(config)
	if err != nil {
		return fmt.Errorf("cannot open database: %v", err)
	}
	defer database.Close()
	err = database.Update(func(tx *bolt.Tx) error {
		log.Infof("Purging source: %s...", sourceName)
		bucket, err := tx.CreateBucketIfNotExists([]byte(sourceName))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		prefix := []byte(sourceName + "/")
		cursor := bucket.Cursor()
		numPurged := 0
		for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
			date := new(time.Time)
			date.UnmarshalText(v)
			if date.Before(time.Now().AddDate(0, 0, -numOfDays)) {
				numPurged++
				cursor.Delete()
			}
		}
		log.Infof("Purged %d entries", numPurged)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error in updating the database: %v", err)
	}
	return nil
}

func filterUnSeen(source *goeland.Source, params *filterParams) {
	database, err := openDatabase(params.config)
	if err != nil {
		log.Errorf("cannot open database: %v", err)
		return
	}
	defer database.Close()
	err = database.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(source.Name))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		var current int
		for _, entry := range source.Entries {
			key := []byte(source.Name + "/" + entry.UID)
			value := bucket.Get(key)
			if now, err := time.Now().MarshalText(); err == nil {
				if err := bucket.Put(key, now); err != nil {
					log.Debugf("error recording seen status for key: %s", string(key))
				}
			}
			if value != nil {
				log.Debugf("already seen entry with key: %s", string(key))
				continue
			}
			source.Entries[current] = entry
			current++
		}
		source.Entries = source.Entries[:current]
		return nil
	})

	if err != nil {
		log.Debugf("error in updating the database: %v", err)
	}
}
