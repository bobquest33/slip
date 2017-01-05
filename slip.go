package slip

import (
	"errors"
	"github.com/rameshvk/cache"
	"github.com/rameshvk/linkinfo"
	"io/ioutil"
	"log"
	"time"
)

type Slip struct {
	ImageCache      cache.ICache
	InfoCache       cache.ICache
	ImageExpiration *time.Duration
	InfoExpiration  *time.Duration
	ImageUrlPrefix  string
	InfoUrlPrefix   string
}

func (slip *Slip) cacheImageInfo(url string, filepath string) (string, error) {
	key := hash256(url)

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	var expires *time.Time
	if slip.ImageExpiration != nil {
		t := time.Now().Add(*slip.ImageExpiration)
		expires = &t
	}
	if err := slip.ImageCache.Add(key, bytes, expires); err != nil {
		return "", err
	}

	return slip.ImageUrlPrefix + key, nil
}

func (slip *Slip) CacheInfo(info *linkinfo.Info) (*linkinfo.Info, error) {
	key := hash256(info.CanonicalUrl)

	// cache check
	if result, _, err := slip.InfoCache.Get(key); err != nil {
		if info, ok := result.(linkinfo.Info); !ok {
			log.Panic("Unexpected info type result", result)
			return nil, errors.New("Sorry")
		} else {
			return &info, nil
		}
	}

	var result linkinfo.Info
	result = *info

	// not there in cache, rename all the image urls
	if info.Images != nil && len(info.Images) > 0 {
		result.Images = []linkinfo.ImageInfo{}
		for _, image := range info.Images {
			if newUrl, err := slip.cacheImageInfo(image.ImageUrl, image.LocalTempFilePath); err != nil {
				return nil, err
			} else {
				image.ImageUrl = newUrl
				image.LocalTempFilePath = ""
				result.Images = append(result.Images, image)
			}
		}
	}

	// result is good to be saved, write it to the cache
	var expires *time.Time
	if slip.InfoExpiration != nil {
		t := time.Now().Add(*slip.InfoExpiration)
		expires = &t
	}
	if err := slip.InfoCache.Add(key, result, expires); err != nil {
		return nil, err
	}

	return &result, nil
}
