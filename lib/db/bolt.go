package db

import (
	"github.com/boltdb/bolt"
	"github.com/zctod/go-tool/common/utils"
	"log"
)

/**
 * Created by zc on 2019-10-29.
 */

const (
	CUBA = "cuba"
	UUC = "uuc"
	Store = "store"
)

var DB *bolt.DB

func init() {
	if DB != nil {
		return
	}

	if err := utils.PathCreate("db"); err != nil {
		log.Fatal("bolt文件夹创建失败：", err)
	}

	var err error
	DB, err = open()
	if err != nil {
		log.Fatal("bolt连接失败：", err)
	}

	err = DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(Store))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(UUC))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(CUBA))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal("bolt初始化bucket失败：", err)
	}
}

func open() (*bolt.DB, error) {
	return bolt.Open("db/mock.db", 0755, nil)
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
