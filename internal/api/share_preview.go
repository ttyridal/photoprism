package api

import (
	"fmt"
	"image"
	"image/color"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/entity"
	"github.com/photoprism/photoprism/internal/form"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/photoprism/photoprism/internal/query"
	"github.com/photoprism/photoprism/internal/service"
	"github.com/photoprism/photoprism/internal/thumb"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/txt"
)

// SharePreview returns a link share preview image.
//
// GET /s/:token/:uid/preview
// TODO: Proof of concept, needs refactoring.
func SharePreview(router *gin.RouterGroup) {
	router.GET("/:token/:share/preview", func(c *gin.Context) {
		conf := service.Config()

		token := c.Param("token")
		share := c.Param("share")
		links := entity.FindLinks(token, share)

		if len(links) != 1 {
			log.Warn("share: invalid token (preview)")
			c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
			return
		}

		thumbPath := path.Join(conf.ThumbPath(), "share")

		if err := os.MkdirAll(thumbPath, os.ModePerm); err != nil {
			log.Error(err)
			c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
			return
		}

		previewFilename := fmt.Sprintf("%s/%s.jpg", thumbPath, share)
		yesterday := time.Now().Add(-24 * time.Hour)

		if info, err := os.Stat(previewFilename); err != nil {
			log.Debugf("share: creating new preview for %s", share)
		} else if info.ModTime().After(yesterday) {
			log.Debugf("share: using cached preview for %s", share)
			c.File(previewFilename)
			return
		} else if err := os.Remove(previewFilename); err != nil {
			log.Errorf("share: could not remove old preview of %s", share)
			c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
			return
		}

		var f form.PhotoSearch

		// Previews may only contain public content in shared albums.
		f.Album = share
		f.Public = true
		f.Private = false
		f.Hidden = false
		f.Archived = false
		f.Review = false
		f.Primary = true

		// Get first 12 album entries.
		f.Count = 12
		f.Order = "relevance"

		p, count, err := query.PhotoSearch(f)

		if err != nil {
			log.Error(err)
			c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
			return
		}

		if count == 0 {
			c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
			return
		} else if count < 12 {
			f := p[0]
			size, _ := thumb.Sizes[thumb.Fit720]

			fileName := photoprism.FileName(f.FileRoot, f.FileName)

			if !fs.FileExists(fileName) {
				log.Errorf("share: file %s is missing (preview)", txt.Quote(f.FileName))
				c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
				return
			}

			thumbnail, err := thumb.FromFile(fileName, f.FileHash, conf.ThumbPath(), size.Width, size.Height, f.FileOrientation, size.Options...)

			if err != nil {
				log.Error(err)
				c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
				return
			}

			c.File(thumbnail)

			return
		}

		width := 908
		height := 680
		x := 0
		y := 0

		preview := imaging.New(width, height, color.NRGBA{255, 255, 255, 255})
		size, _ := thumb.Sizes[thumb.Tile224]

		for _, f := range p {
			fileName := photoprism.FileName(f.FileRoot, f.FileName)

			if !fs.FileExists(fileName) {
				log.Errorf("share: file %s is missing (preview)", txt.Quote(f.FileName))
				c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
				return
			}

			thumbnail, err := thumb.FromFile(fileName, f.FileHash, conf.ThumbPath(), size.Width, size.Height, f.FileOrientation, size.Options...)

			if err != nil {
				log.Error(err)
				c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
				return
			}

			src, err := imaging.Open(thumbnail)

			if err != nil {
				log.Error(err)
				c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
				return
			}

			preview = imaging.Paste(preview, src, image.Pt(x, y))

			x += 228

			if x > width {
				x = 0
				y += 228
			}
		}

		// Save the resulting image as JPEG.
		err = imaging.Save(preview, previewFilename)

		if err != nil {
			log.Error(err)
			c.Redirect(http.StatusTemporaryRedirect, conf.SitePreview())
			return
		}

		c.File(previewFilename)
	})
}
