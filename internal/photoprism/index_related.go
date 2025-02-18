package photoprism

import (
	"fmt"
	"path/filepath"

	"github.com/photoprism/photoprism/internal/query"
	"github.com/photoprism/photoprism/pkg/txt"
)

// IndexMain indexes the main file from a group of related files and returns the result.
func IndexMain(related *RelatedFiles, ind *Index, opt IndexOptions) (result IndexResult) {
	// Skip sidecar files without related media file.
	if related.Main == nil {
		result.Err = fmt.Errorf("index: no main file found for %s", txt.Quote(related.String()))
		result.Status = IndexFailed
		return result
	}

	f := related.Main
	sizeLimit := ind.conf.OriginalsLimit()

	// Enforce file size limit for originals.
	if sizeLimit > 0 && f.FileSize() > sizeLimit {
		result.Err = fmt.Errorf("index: %s exceeds file size limit (%d / %d MB)", txt.Quote(f.BaseName()), f.FileSize()/(1024*1024), sizeLimit/(1024*1024))
		result.Status = IndexFailed
		return result
	}

	if f.NeedsExifToolJson() {
		if jsonName, err := ind.convert.ToJson(f); err != nil {
			log.Debugf("index: %s in %s (extract metadata)", txt.Quote(err.Error()), txt.Quote(f.BaseName()))
		} else {
			log.Debugf("index: %s created", filepath.Base(jsonName))
		}
	}

	if opt.Convert && f.IsMedia() && !f.HasJpeg() {
		if jpegFile, err := ind.convert.ToJpeg(f); err != nil {
			result.Err = fmt.Errorf("index: failed converting %s to jpeg (%s)", txt.Quote(f.BaseName()), err.Error())
			result.Status = IndexFailed

			return result
		} else {
			log.Debugf("index: %s created", txt.Quote(jpegFile.BaseName()))

			if err := jpegFile.ResampleDefault(ind.thumbPath(), false); err != nil {
				result.Err = fmt.Errorf("index: failed creating thumbs for %s (%s)", txt.Quote(f.BaseName()), err.Error())
				result.Status = IndexFailed

				return result
			}

			related.Files = append(related.Files, jpegFile)
		}
	}

	result = ind.MediaFile(f, opt, "")

	if result.Indexed() && f.IsJpeg() {
		if err := f.ResampleDefault(ind.thumbPath(), false); err != nil {
			log.Errorf("index: failed creating thumbs for %s (%s)", txt.Quote(f.BaseName()), err.Error())
			query.SetFileError(result.FileUID, err.Error())
		}
	}

	log.Infof("index: %s main %s file %s", result, f.FileType(), txt.Quote(f.RelName(ind.originalsPath())))

	return result
}

// IndexMain indexes a group of related files and returns the result.
func IndexRelated(related RelatedFiles, ind *Index, opt IndexOptions) (result IndexResult) {
	done := make(map[string]bool)
	sizeLimit := ind.conf.OriginalsLimit()

	result = IndexMain(&related, ind, opt)

	if result.Failed() {
		log.Warn(result.Err)
		return result
	} else if !result.Success() || result.Stacked() {
		// Skip related files if main file was stacked or indexing was not completely successful.
		return result
	}

	done[related.Main.FileName()] = true

	i := 0

	for i < len(related.Files) {
		f := related.Files[i]
		i++

		if f == nil {
			continue
		}

		if done[f.FileName()] {
			continue
		}

		done[f.FileName()] = true

		// Enforce file size limit for originals.
		if sizeLimit > 0 && f.FileSize() > sizeLimit {
			log.Warnf("index: %s exceeds file size limit (%d / %d MB)", txt.Quote(f.BaseName()), f.FileSize()/(1024*1024), sizeLimit/(1024*1024))
			continue
		}

		if f.NeedsExifToolJson() {
			if jsonName, err := ind.convert.ToJson(f); err != nil {
				log.Debugf("index: %s in %s (extract metadata)", txt.Quote(err.Error()), txt.Quote(f.BaseName()))
			} else {
				log.Debugf("index: %s created", filepath.Base(jsonName))
			}
		}

		if opt.Convert && f.IsMedia() && !f.HasJpeg() {
			if jpegFile, err := ind.convert.ToJpeg(f); err != nil {
				result.Err = fmt.Errorf("index: failed converting %s to jpeg (%s)", txt.Quote(f.BaseName()), err.Error())
				result.Status = IndexFailed

				return result
			} else {
				log.Debugf("index: %s created", txt.Quote(jpegFile.BaseName()))

				if err := jpegFile.ResampleDefault(ind.thumbPath(), false); err != nil {
					result.Err = fmt.Errorf("index: failed creating thumbs for %s (%s)", txt.Quote(f.BaseName()), err.Error())
					result.Status = IndexFailed

					return result
				}

				related.Files = append(related.Files, jpegFile)
			}
		}

		res := ind.MediaFile(f, opt, "")

		if res.Indexed() && f.IsJpeg() {
			if err := f.ResampleDefault(ind.thumbPath(), false); err != nil {
				log.Errorf("index: failed creating thumbs for %s (%s)", txt.Quote(f.BaseName()), err.Error())
				query.SetFileError(res.FileUID, err.Error())
			}
		}

		log.Infof("index: %s related %s file %s", res, f.FileType(), txt.Quote(f.BaseName()))
	}

	return result
}
