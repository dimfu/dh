package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.etcd.io/bbolt"
)

const (
	BUCKET = "dirs"
)

var (
	args    = make([]string, 0)
	command string
)

func init() {
	flag.Parse()
	for i := len(os.Args) - len(flag.Args()) + 1; i < len(os.Args); {
		if i > 1 && os.Args[i-2] == "--" {
			break
		}
		args = append(args, flag.Arg(0))
		if err := flag.CommandLine.Parse(os.Args[i:]); err != nil {
			log.Fatal("error while parsing arguments")
		}

		i += 1 + len(os.Args[i:]) - len(flag.Args())
	}
	args = append(args, flag.Args()...)

	if len(args) < 1 {
		flag.Usage()
		os.Exit(0)
	}
	command = args[0]
}

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dir := fmt.Sprintf("%s/%s", homedir, "dh.db")
	db, err := bbolt.Open(dir, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BUCKET))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	switch command {
	case "add":
		if len(args) <= 1 {
			log.Fatal("dir key is required when using command `add`\n")
		}

		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("error getting working dir: %v\n", err)
		}

		db.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(BUCKET))
			err := b.Put([]byte(args[1]), []byte(dir))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		break
	case "list":
		db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(BUCKET))
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("key=%s, value=%s\n", k, v)
			}
			return nil
		})
	case "goto":
		if len(args) <= 1 {
			log.Fatal("dir key is required when using command `goto`\n")
		}

		db.View(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(BUCKET))
			v := b.Get([]byte(args[1]))

			if v == nil {
				log.Fatalf("key %s dir not found\n", args[1])
			}

			strval := string(v)
			fmt.Println(strval)
			return nil
		})
	case "delete":
		if len(args) <= 1 {
			log.Fatal("dir key is required when using command `goto`\n")
		}

		db.Update(func(tx *bbolt.Tx) error {
			key := []byte(args[1])
			b := tx.Bucket([]byte(BUCKET))
			v := b.Get(key)

			if v == nil {
				log.Fatalf("key %s dir not found\n", args[1])
			}

			err := b.Delete(key)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Successfully deleted %v from the bookmark\n", string(v))

			return nil
		})

	default:
		log.Fatal("command not found")
	}
}
