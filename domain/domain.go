package domain

import (
	"opensource/chaos/domain/dao/mongo"
)

func Close() {
	mongo.MongoClose()
}
