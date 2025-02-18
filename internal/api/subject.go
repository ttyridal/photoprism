package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/photoprism/photoprism/internal/acl"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/event"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/i18n"
	"github.com/photoprism/photoprism/internal/query"
	"github.com/photoprism/photoprism/pkg/txt"
)

// GetSubjects finds and returns subjects as JSON.
//
// GET /api/v1/subjects
func GetSubjects(router *gin.RouterGroup) {
	router.GET("/subjects", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourceSubjects, acl.ActionSearch)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		var f form.SubjectSearch

		err := c.MustBindWith(&f, binding.Form)

		if err != nil {
			AbortBadRequest(c)
			return
		}

		result, err := query.SubjectSearch(f)

		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		AddCountHeader(c, len(result))
		AddLimitHeader(c, f.Count)
		AddOffsetHeader(c, f.Offset)
		AddTokenHeaders(c)

		c.JSON(http.StatusOK, result)
	})
}

// GetSubject returns a subject as JSON.
//
// GET /api/v1/subjects/:uid
func GetSubject(router *gin.RouterGroup) {
	router.GET("/subjects/:uid", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourceSubjects, acl.ActionRead)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		if subj := entity.FindSubject(c.Param("uid")); subj == nil {
			Abort(c, http.StatusNotFound, i18n.ErrSubjectNotFound)
			return
		} else {
			c.JSON(http.StatusOK, subj)
		}
	})
}

// UpdateSubject updates subject properties.
//
// PUT /api/v1/subjects/:uid
func UpdateSubject(router *gin.RouterGroup) {
	router.PUT("/subjects/:uid", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourceSubjects, acl.ActionUpdate)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		var f form.Subject

		if err := c.BindJSON(&f); err != nil {
			AbortBadRequest(c)
			return
		}

		uid := c.Param("uid")
		m := entity.FindSubject(uid)

		if m == nil {
			Abort(c, http.StatusNotFound, i18n.ErrSubjectNotFound)
			return
		}

		if _, err := m.UpdateName(f.SubjName); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		event.SuccessMsg(i18n.MsgSubjectSaved)

		c.JSON(http.StatusOK, m)
	})
}

// LikeSubject flags a subject as favorite.
//
// POST /api/v1/subjects/:uid/like
//
// Parameters:
//   uid: string Subject UID
func LikeSubject(router *gin.RouterGroup) {
	router.POST("/subjects/:uid/like", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourceSubjects, acl.ActionUpdate)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		uid := c.Param("uid")
		subj := entity.FindSubject(uid)

		if subj == nil {
			Abort(c, http.StatusNotFound, i18n.ErrSubjectNotFound)
			return
		}

		if err := subj.Update("SubjFavorite", true); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		PublishSubjectEvent(EntityUpdated, uid, c)

		c.JSON(http.StatusOK, http.Response{})
	})
}

// DislikeSubject removes the favorite flag from a subject.
//
// DELETE /api/v1/subjects/:uid/like
//
// Parameters:
//   uid: string Subject UID
func DislikeSubject(router *gin.RouterGroup) {
	router.DELETE("/subjects/:uid/like", func(c *gin.Context) {
		s := Auth(SessionID(c), acl.ResourceSubjects, acl.ActionUpdate)

		if s.Invalid() {
			AbortUnauthorized(c)
			return
		}

		uid := c.Param("uid")
		subj := entity.FindSubject(uid)

		if subj == nil {
			Abort(c, http.StatusNotFound, i18n.ErrSubjectNotFound)
			return
		}

		if err := subj.Update("SubjFavorite", false); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": txt.UcFirst(err.Error())})
			return
		}

		PublishSubjectEvent(EntityUpdated, uid, c)

		c.JSON(http.StatusOK, http.Response{})
	})
}
