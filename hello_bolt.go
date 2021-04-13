package main

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/rand"
	"time"
)

func i64tob(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint64(0); i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}
	return r
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := bolt.Open("hello_bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	writeFlag := flag.Bool("write", false, "write 10 random entries to db")

	readFlag := flag.Bool("read", false, "read all entries in db")

	flag.Parse()

	if *writeFlag {
		fmt.Println("writing 10 random entries")
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("YeOldeBucket"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			for i := 0; i < 10; i++ {
				id, _ := b.NextSequence()
				putError := b.Put(i64tob(id), []byte(RandStringRunes(15)))
				if putError != nil {
					return fmt.Errorf("write error: %s", putError)
				}

			}

			return nil

		})
	}

	if *readFlag {
		fmt.Println("reading")
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("YeOldeBucket"))

			if b == nil {
				fmt.Printf("database has 0 entries")
				return nil
			}

			c := b.Cursor()

			count := 0

			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("key=%d, value=%s\n", k, v)
				count++
			}
			fmt.Printf("database has %d entries", count)
			return nil
		})
	}

}