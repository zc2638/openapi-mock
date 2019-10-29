package db

import (
	"github.com/boltdb/bolt"
	"github.com/zctod/go-tool/common/utils"
	"log"
)

/**
 * Created by zc on 2019-10-29.
 */

var DB *bolt.DB

func init() {

	if DB != nil {
		return
	}
	var err error
	if err := utils.PathCreate("db"); err != nil {
		log.Fatal(err)
	}

	DB, err = bolt.Open("db/mock.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Update(bucket, key, value string) error {

	return DB.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
}

func View(bucket, key string) ([]byte, error) {

	var buffer []byte
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		buffer = b.Get([]byte(key))
		return nil
	})
	return buffer, err
}

func Delete(bucket, key string) error {

	return DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Delete([]byte(key))
	})
}
