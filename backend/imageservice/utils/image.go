package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nfnt/resize"
	"go.uber.org/zap"
)

type ImageUtils struct {
	logger *zap.Logger
}

func NewImageUtils(logger *zap.Logger) *ImageUtils {
	return &ImageUtils{
		logger: logger,
	}
}

var mu sync.Mutex

func (util *ImageUtils) Resize(filename string, ext string) {
	var wg sync.WaitGroup
	sizes := []uint{100, 200, 300}
	wg.Add(len(sizes))
	for _, size := range sizes {
		// Runnig go routine to resize image to different sizes and save.
		go func(s uint) {
			defer wg.Done()
			util.processImage(filename, ext, s)
		}(size)
	}

	wg.Wait()
	util.logger.Info("All images processed.")
}

func (util *ImageUtils) processImage(filename string, ext string, size uint) {

	destinationOriginal := filepath.Join("./uploads", filename)
	src, err := os.Open(destinationOriginal)
	if err != nil {
		util.logger.Error("Error opening file: " + err.Error())
		return
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		fmt.Printf("Error decoding %s: %v\n", filename, err)
		return
	}

	resizedImg := resize.Resize(size, 0, img, resize.Lanczos3)
	newFilename := util.getNewFilename(filename, size)
	destination := filepath.Join("./uploads", newFilename)

	dst, err := os.Create(destination)
	if err != nil {
		util.logger.Error("Error creating file: " + err.Error())
		return
	}
	defer dst.Close()

	if strings.ToLower(ext) == ".jpg" || ext == ".jpeg" {
		if err := jpeg.Encode(dst, resizedImg, nil); err != nil {
			util.logger.Error("Error resizing image: " + err.Error())
		}
	} else if strings.ToLower(ext) == ".png" {
		if err := png.Encode(dst, resizedImg); err != nil {
			util.logger.Error("Error resizing image: " + err.Error())
		}
	}
}

func (util *ImageUtils) getNewFilename(filename string, size uint) string {
	mu.Lock()
	defer mu.Unlock()

	filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	newFilename := fmt.Sprintf("%s_%d%s", filenameWithoutExt, size, filepath.Ext(filename))
	return newFilename
}

func (util *ImageUtils) RemoveImageFile(imageId string) {
	var wgr sync.WaitGroup
	sizes := []uint{0, 100, 200, 300}
	wgr.Add(len(sizes))
	for _, size := range sizes {
		go func(id string, s uint) {
			defer wgr.Done()
			util.remove(id, s)
		}(imageId, size)
	}
	wgr.Wait()
	util.logger.Info("All images removed.")
}

func (util *ImageUtils) remove(imageId string, size uint) {
	if size > 0 {
		imageId = util.getNewFilename(imageId, size)
	}
	filePath := filepath.Join("./uploads", imageId)
	_, err := os.Stat(filePath)
	if err == nil {
		if err := os.Remove(filePath); err != nil {

		}
	}
}
