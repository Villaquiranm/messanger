package database

import (
	"encoding/json"
	"fmt"
	"log"
	"messager/config"
	"messager/model"
	"strconv"

	"github.com/boltdb/bolt"
)

// DBManager manages the chat database
type DBManager struct {
	db *bolt.DB
}

// NewDBManager creates a manager of the bolt database
func NewDBManager() *DBManager {
	db, err := bolt.Open(config.BoltDirectory, 0600, nil)
	if err != nil {
		log.Fatal("Error while openign boldDB: %s", err.Error())
	}
	return &DBManager{db: db}
}

// UserExist returns if an user with that name is already in the chat
func (m *DBManager) UserExist(name string) bool {
	var user []byte
	err := m.db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.UserBucket))
		if b == nil {
			return nil
		}
		// bucket.Get returs nil if key does not exist
		// Users exist if the returned bytes are different than nil
		user = b.Get([]byte(name))
		return nil
	})
	if err != nil {
		fmt.Printf("Error user exist: %s", err.Error())
	}
	return user != nil
}

// CreateUser Create a new user in the chat
func (m *DBManager) CreateUser(name string) error {
	// getCurrentMessageIndex user will not have access to previous messages
	currentIndex, err := m.getCurrentMessageIndex()
	if err != nil {
		return nil
	}
	return m.db.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucketIfNotExists([]byte(config.UserBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(name), []byte(fmt.Sprint(currentIndex+1)))
	})
}

// DeleteUser Deletes the user from the chat
func (m *DBManager) DeleteUser(name string) error {
	return m.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.UserBucket))
		return b.Delete([]byte(name))
	})
}

// StoreMessage Stores one message on the database
func (m *DBManager) StoreMessage(msg model.Message) error {
	m.db.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucketIfNotExists([]byte(config.MessagesBucket))
		if err != nil {
			return err
		}
		next, err := b.NextSequence()
		if err != nil {
			return err
		}
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		return b.Put([]byte(fmt.Sprint(next)), data)
	})
	return nil
}

// MessagesForUser retrieves all not read messages for an specific user
func (m *DBManager) MessagesForUser(name string) ([]model.Message, error) {
	// Find last message read by the user
	lastRead, err := m.getLastRead(name)
	if err != nil {
		return nil, err
	}
	// Find all messages for that user
	messages, lastRead, err := m.getUnreadMessages(lastRead)
	// Update last read for user

	err = m.db.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucketIfNotExists([]byte(config.UserBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(name), []byte(fmt.Sprint(lastRead)))
	})
	return messages, err
}

// Find last read message by user
func (m *DBManager) getLastRead(name string) (int, error) {
	var lastUpdate int
	var err error
	err = m.db.View(func(t *bolt.Tx) error {
		lastUpdate = 1
		b := t.Bucket([]byte(config.UserBucket))
		if b == nil {
			return nil
		}

		// Get the last index of the message read by that user
		lastRead := b.Get([]byte(name))
		if lastRead != nil {
			lastUpdate, err = strconv.Atoi(string(lastRead))
		}
		return err
	})
	return lastUpdate, err
}

// Gets all unreaded messages for that user
func (m *DBManager) getUnreadMessages(lastRead int) ([]model.Message, int, error) {
	messages := make([]model.Message, 0)
	err := m.db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.MessagesBucket))
		if b == nil {
			return nil
		}
		data := b.Get([]byte(strconv.Itoa(lastRead)))
		for data != nil {
			var message model.Message
			err := json.Unmarshal(data, &message)
			if err != nil {
				return err
			}
			messages = append(messages, message)
			lastRead++
			data = b.Get([]byte(strconv.Itoa(lastRead)))
		}
		return nil
	})

	return messages, lastRead, err
}

func (m *DBManager) getCurrentMessageIndex() (uint64, error) {
	var currentIndex uint64
	err := m.db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(config.MessagesBucket))
		if b == nil {
			return nil
		}
		currentIndex = b.Sequence()
		return nil
	})
	return currentIndex, err
}
