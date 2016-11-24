package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"gopkg.in/redis.v5"
	"io"
	"os"
	"path/filepath"
	"fmt"
	"github.com/charoit/fileupdater/settings"
)

type Cache struct {
	client   *redis.Client
	settings *settings.Settings
}

type file struct {
	Path string `json: "path"`
	Hash string `json: "hash"`
	Size string `json: "size"`
}

func GetCache(set *settings.Settings) (*Cache, error) {

	var cache Cache
	cache.settings = set

	cache.client = redis.NewClient(&redis.Options{
		Addr:     set.Redis.Addr,
		Password: set.Redis.Pass, // no password set
		DB:       set.Redis.DB,   // use default DB
	})

	_, err := cache.client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

func (c *Cache) Update() error {

	for _, v := range c.settings.Paths {

		var list []file

		filepath.Walk(c.settings.Root+v.Path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				hash, err := hash(path)
				if err != nil {
					return err
				}

				m := file{path, hash, fmt.Sprint(info.Size())}
				list = append(list, m)
			}
			return nil
		})

		j, err := json.Marshal(list)
		if err != nil {
			return err
		}

		c.client.Set(v.Path, j, 0)
	}

	return nil
}

func (c *Cache) Get(key string) ([]file, error) {

	val, err := c.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	var files []file
	err = json.Unmarshal([]byte(val), &files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (c *Cache) Raw(key string) ([]byte, error) {
	val, err := c.client.Get(key).Bytes()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func hash(path string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	var hashMD5 = hex.EncodeToString(hashInBytes)

	return hashMD5, nil
}
