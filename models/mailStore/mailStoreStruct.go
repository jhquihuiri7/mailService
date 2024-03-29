package mailStore

import (
	"github.com/google/uuid"
)

type MailStore struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
	Mail string `bson:"mail"`
}

func (m *MailStore) AddMail(req map[string]string) {
	m.Id = uuid.New().String()
	for _, toName := range []string{"name", "nombre"} {
		_, ok := req[toName]
		if ok {
			m.Name = toName
			break
		}
	}
}
