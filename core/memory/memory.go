package memory

import (
	"sync"

	"go.mau.fi/whatsmeow"
)

var ServerAddr string

type User struct {
	WhatsClient *whatsmeow.Client
	sync.Mutex
}

var m *mUser

type mUser struct {
	sync.RWMutex
	User map[int]*User
}

var once sync.Once

func Init() {
	once.Do(func() {
		m = &mUser{}
		m.User = make(map[int]*User)
	})

}

func LockUser(socketId int) {
	if socketId == 0 {
		return
	}
	m.User[socketId].Lock()
}

func UnlockUser(socketId int) {
	if socketId == 0 {
		return
	}
	m.User[socketId].Unlock()
}

func SetUser(socketId int, user User) {
	defer func() {
		m.Unlock()
	}()
	m.Lock()
	m.User[socketId] = &user
}

func GetUser(socketId int) (user *User, ok bool) {
	user, ok = m.User[socketId]
	if !ok {
		return nil, ok
	}
	return user, ok
}

func DelUser(socketId int) {
	delete(m.User, socketId)
}

func GetAllUser() map[int]*User {
	return m.User
}
