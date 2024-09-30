package cache

import (
	"encoding/json"
	"sync"
	"time"
)

type safeStorage struct {
	mu     *sync.RWMutex
	chache map[string]chacheStruct
}

var storage = safeStorage{mu: &sync.RWMutex{}, chache: make(map[string]chacheStruct)}

type chacheStruct struct {
	lastRead time.Time
	encoded  []byte
}

func Add[T any](data T, key string) error {
	// fmt.Println("saving to storage ", key)

	// Encode the data with Marshal and lock for mutex
	storage.mu.Lock()
	defer storage.mu.Unlock()
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Add to storage
	cache := chacheStruct{
		lastRead: time.Now(),
		encoded:  encoded,
	}

	// Assuming storage is a map[string]cacheStruct
	storage.chache[key] = cache

	return nil
}

func StartInterval() chan struct{} {
	//time.netTicker
	// fmt.Println("start time ")
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-quit:
				ticker.Stop()
				return
			case <-ticker.C:
				// do something
				Remove()
			}
		}
	}()

	return quit
}

func Remove() {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	currentTime := time.Now()

	for key, value := range storage.chache {
		lastModified := value.lastRead

		if currentTime.Sub(lastModified) > 10*time.Second {
			delete(storage.chache, key)
		}
	}
}

func GetChache(key string) ([]byte, bool) {
	storage.mu.RLock()
	defer storage.mu.RUnlock()
	value, exists := storage.chache[key]
	if !exists {
		// fmt.Println("DOESNT EXIST IN CHACHE")

		return value.encoded, exists
	}

	// fmt.Println("EXISTs IN CHACHE")
	return value.encoded, exists
}
