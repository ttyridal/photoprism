package api

import (
	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/query"
)

type EntityEvent string

const (
	EntityUpdated EntityEvent = "updated"
	EntityCreated EntityEvent = "created"
	EntityDeleted EntityEvent = "deleted"
)

func PublishPhotoEvent(e EntityEvent, uid string, c *gin.Context) {
	f := form.PhotoSearch{ID: uid, Merged: true}
	result, _, err := query.PhotoSearch(f)

	if err != nil {
		log.Error(err)
		AbortUnexpected(c)
		return
	}

	event.PublishEntities("photos", string(e), result)
}

func PublishAlbumEvent(e EntityEvent, uid string, c *gin.Context) {
	f := form.AlbumSearch{ID: uid}
	result, err := query.AlbumSearch(f)

	if err != nil {
		log.Error(err)
		AbortUnexpected(c)
		return
	}

	event.PublishEntities("albums", string(e), result)
}

func PublishLabelEvent(e EntityEvent, uid string, c *gin.Context) {
	f := form.LabelSearch{ID: uid}
	result, err := query.Labels(f)

	if err != nil {
		log.Error(err)
		AbortUnexpected(c)
		return
	}

	event.PublishEntities("labels", string(e), result)
}

func PublishSubjectEvent(e EntityEvent, uid string, c *gin.Context) {
	f := form.SubjectSearch{ID: uid}
	result, err := query.SubjectSearch(f)

	if err != nil {
		log.Error(err)
		AbortUnexpected(c)
		return
	}

	event.PublishEntities("subjects", string(e), result)
}
