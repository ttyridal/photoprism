package api

import (
	"net/http"

	"github.com/photoprism/photoprism/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/photoprism/photoprism/internal/photoprism"
	"github.com/photoprism/photoprism/internal/query"
	"github.com/photoprism/photoprism/internal/video"
	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/txt"
)

// GET /api/v1/videos/:hash/:token/:type
//
// Parameters:
//   hash: string The photo or video file hash as returned by the search API
//   type: string Video type
func GetVideo(router *gin.RouterGroup) {
	router.GET("/videos/:hash/:token/:type", func(c *gin.Context) {
		if InvalidPreviewToken(c) {
			c.Data(http.StatusForbidden, "image/svg+xml", brokenIconSvg)
			return
		}

		fileHash := c.Param("hash")
		typeName := c.Param("type")

		videoType, ok := video.Types[typeName]

		if !ok {
			log.Errorf("video: invalid type %s", txt.Quote(typeName))
			c.Data(http.StatusOK, "image/svg+xml", videoIconSvg)
			return
		}

		f, err := query.FileByHash(fileHash)

		if err != nil {
			log.Errorf("video: %s", err.Error())
			c.Data(http.StatusOK, "image/svg+xml", videoIconSvg)
			return
		}

		if !f.FileVideo {
			f, err = query.VideoByPhotoUID(f.PhotoUID)

			if err != nil {
				log.Errorf("video: %s", err.Error())
				c.Data(http.StatusOK, "image/svg+xml", videoIconSvg)
				return
			}
		}

		if f.FileError != "" {
			log.Errorf("video: file error %s", f.FileError)
			c.Data(http.StatusOK, "image/svg+xml", videoIconSvg)
			return
		}

		fileName := photoprism.FileName(f.FileRoot, f.FileName)

		mf, err := photoprism.NewMediaFile(fileName)
		if err != nil {
			log.Errorf("video: file %s is missing", txt.Quote(f.FileName))
			c.Data(http.StatusOK, "image/svg+xml", videoIconSvg)

			// Set missing flag so that the file doesn't show up in search results anymore.
			logError("video", f.Update("FileMissing", true))

			return
		}

		conf := service.Config()
		avcName := videoType.Format.FindFirst(mf.FileName(), []string{conf.SidecarPath(), fs.HiddenPath}, conf.OriginalsPath(), false)
		mediaFile, err := photoprism.NewMediaFile(avcName)
		if err == nil && mediaFile.IsVideo() {
			fileName = mediaFile.FileName()
		} else if f.FileCodec != string(videoType.Codec) {
			conv := service.Convert()

			if avcFile, err := conv.ToAvc(mf, service.Config().FFmpegEncoder()); err != nil {
				log.Errorf("video: transcoding %s failed", txt.Quote(f.FileName))
				c.Data(http.StatusOK, "image/svg+xml", videoIconSvg)
				return
			} else {
				fileName = avcFile.FileName()
			}
		}

		AddContentTypeHeader(c, ContentTypeAvc)

		if c.Query("download") != "" {
			c.FileAttachment(fileName, f.DownloadName(DownloadName(c), 0))
		} else {
			c.File(fileName)
		}

		return
	})
}
