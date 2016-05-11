package domain

import (
	"opensource/chaos/background/server/domain/mongo"
)

func Close() {
	mongo.MongoClose()
}
