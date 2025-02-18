package api

import (
	"net/http"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/stretchr/testify/assert"
)

func TestGetSubjects(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		app, router, _ := NewApiTest()
		GetSubjects(router)
		r := PerformRequest(app, "GET", "/api/v1/subjects?count=10")
		count := gjson.Get(r.Body.String(), "#")
		assert.LessOrEqual(t, int64(3), count.Int())
		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("invalid request", func(t *testing.T) {
		app, router, _ := NewApiTest()
		GetSubjects(router)
		r := PerformRequest(app, "GET", "/api/v1/subjects?xxx=10")
		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestGetSubject(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		app, router, _ := NewApiTest()
		GetSubject(router)
		r := PerformRequest(app, "GET", "/api/v1/subjects/jqy1y111h1njaaaa")
		val := gjson.Get(r.Body.String(), "Slug")
		assert.Equal(t, "dangling-subject", val.String())
		assert.Equal(t, http.StatusOK, r.Code)
	})
	t.Run("invalid request", func(t *testing.T) {
		app, router, _ := NewApiTest()
		GetSubject(router)
		r := PerformRequest(app, "GET", "/api/v1/subjects/xxx1y111h1njaaaa")
		val := gjson.Get(r.Body.String(), "error")
		assert.Equal(t, "Subject not found", val.String())
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
}

func TestLikeSubject(t *testing.T) {
	t.Run("like not existing subject", func(t *testing.T) {
		app, router, _ := NewApiTest()
		LikeSubject(router)
		r := PerformRequest(app, "POST", "/api/v1/subjects/8775789/like")
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
	t.Run("like existing subject", func(t *testing.T) {
		app, router, _ := NewApiTest()

		// Register routes.
		GetSubject(router)
		LikeSubject(router)

		r := PerformRequest(app, "GET", "/api/v1/subjects/jqy3y652h8njw0sx")
		t.Log(r.Body.String())
		val := gjson.Get(r.Body.String(), "Slug")
		assert.Equal(t, "joe-biden", val.String())
		val2 := gjson.Get(r.Body.String(), "Favorite")
		assert.Equal(t, "false", val2.String())

		r2 := PerformRequest(app, "POST", "/api/v1/subjects/jqy3y652h8njw0sx/like")
		t.Log(r2.Body.String())
		assert.Equal(t, http.StatusOK, r2.Code)

		r3 := PerformRequest(app, "GET", "/api/v1/subjects/jqy3y652h8njw0sx")
		t.Log(r3.Body.String())
		val3 := gjson.Get(r3.Body.String(), "Slug")
		assert.Equal(t, "joe-biden", val3.String())
		val4 := gjson.Get(r3.Body.String(), "Favorite")
		assert.Equal(t, "true", val4.String())
	})
}
func TestDislikeSubject(t *testing.T) {
	t.Run("dislike not existing subject", func(t *testing.T) {
		app, router, _ := NewApiTest()
		DislikeSubject(router)
		r := PerformRequest(app, "DELETE", "/api/v1/subjects/8775789/like")
		assert.Equal(t, http.StatusNotFound, r.Code)
	})
	t.Run("dislike existing subject", func(t *testing.T) {
		app, router, _ := NewApiTest()

		// Register routes.
		GetSubject(router)
		DislikeSubject(router)

		r := PerformRequest(app, "GET", "/api/v1/subjects/jqu0xs11qekk9jx8")
		t.Log(r.Body.String())
		val := gjson.Get(r.Body.String(), "Slug")
		assert.Equal(t, "john-doe", val.String())
		val2 := gjson.Get(r.Body.String(), "Favorite")
		assert.Equal(t, "true", val2.String())

		r2 := PerformRequest(app, "DELETE", "/api/v1/subjects/jqu0xs11qekk9jx8/like")
		t.Log(r2.Body.String())
		assert.Equal(t, http.StatusOK, r2.Code)

		r3 := PerformRequest(app, "GET", "/api/v1/subjects/jqu0xs11qekk9jx8")
		t.Log(r3.Body.String())
		val3 := gjson.Get(r3.Body.String(), "Slug")
		assert.Equal(t, "john-doe", val3.String())
		val4 := gjson.Get(r3.Body.String(), "Favorite")
		assert.Equal(t, "false", val4.String())
	})
}
