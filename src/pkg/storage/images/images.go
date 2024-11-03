package images

import (
	"errors"
	"io"

	"github.com/davidbyttow/govips/v2/vips"
)

type FitType string

const (
	FitTypeCover   FitType = "cover"
	FitTypeContain FitType = "contain"
)

func (f FitType) String() string {
	return string(f)
}

func (f FitType) Set(value string) error {
	switch value {
	case "cover":
		f = FitTypeCover
	case "contain":
		f = FitTypeContain
	default:
		return errors.New("invalid fit type")
	}
	return nil
}

type ImageExportParams struct {
	Format        string
	Quality       int
	StripMetadata bool
}

type ImageMetadata struct {
	Width  int
	Height int
	Format string
}

type ThumbnailParams struct {
	Width  int
	Height int
	Fit    FitType // cover (the image will be resized to cover the width and height specified, cropping the image if necessary), contain (the image will be resized to fit within the width and height specified, maintaining the aspect ratio)
}

func init() {
	vips.LoggingSettings(nil, vips.LogLevelError)
}

type Image interface {
	// close the image
	Close() error
	// Make a thumbnail of the image
	// If crop is true, the image will be cropped to fit the box with the given width and height
	// If crop is false, the image will be resized to fit the box (the actual image size will be smaller than the box)
	Thumbnail(params ThumbnailParams) (Image, error)
	// Export the image to a byte array with the given parameters (format, quality, stripMetadata)
	Export(params ImageExportParams) ([]byte, ImageMetadata, error)
	// Get the image metadata
	GetInfo() (ImageMetadata, error)
	// Copy the internal vips.ImageRef
	CopyImageRef() (Image, error)
}

type ImageLoader interface {
	Image
	LoadFromReader(reader io.Reader) error
}

// wrapper for vips.ImageRef
type vipsImage struct {
	image *vips.ImageRef
}

func newVipsImageFromVipsImage(image *vips.ImageRef) Image {
	return &vipsImage{
		image: image,
	}
}

func newVipsImageFromReader(reader io.Reader) (Image, error) {
	i := &vipsImage{}
	err := i.LoadFromReader(reader)
	if err != nil {
		return nil, err
	}
	return i, err
}

func (i *vipsImage) CopyImageRef() (Image, error) {
	if i.image == nil {
		return nil, errors.New("image not loaded")
	}

	copy, err := i.image.Copy()
	if err != nil {
		return nil, err
	}
	return newVipsImageFromVipsImage(copy), nil
}

func (i *vipsImage) Close() error {
	if i.image != nil {
		i.image.Close()
		i.image = nil
	}
	return nil
}

func (i *vipsImage) LoadFromReader(reader io.Reader) error {
	if i.image != nil {
		return nil
	}
	importParams := vips.NewImportParams()
	importParams.NumPages.Set(-1)
	buf, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	image, err := vips.LoadImageFromBuffer(buf, importParams)
	if err == nil {
		i.image = image
	}
	return err
}

func (i *vipsImage) Thumbnail(params ThumbnailParams) (Image, error) {
	if i.image == nil {
		return nil, errors.New("image not loaded")
	}

	if params.Fit == "" {
		params.Fit = FitTypeContain
	}

	width := params.Width
	height := params.Height
	fit := params.Fit

	// get the image size
	imageWidth := i.image.Width()
	imageHeight := i.image.Height()
	aspectRatio := float64(imageWidth) / float64(imageHeight)

	// when crop is false, the image will be resized to fit the box (the actual image size will be smaller than the box)
	// if width or height is 0, the image will be resized using crop since the image will allways be smaller than the box
	if (fit == FitTypeCover) && (width > 0 && height > 0) {

		maxWidth := width
		maxHeight := height

		realImageHeight := imageHeight / i.image.Pages()
		realAspectRatio := float64(imageWidth) / float64(realImageHeight)

		if imageWidth > maxWidth || realImageHeight > maxHeight {
			if realAspectRatio > 1 {
				maxHeight = int(float64(maxWidth) / realAspectRatio)
				if maxHeight < height {
					maxHeight = height
					maxWidth = int(float64(height) * realAspectRatio)
				}
			} else {
				maxWidth = int(float64(maxHeight) * realAspectRatio)
				if maxWidth < width {
					maxWidth = width
					maxHeight = int(float64(width) / realAspectRatio)
				}
			}
		}

		width = maxWidth
		height = maxHeight * i.image.Pages()

	} else {
		if width == 0 {
			width = i.image.Width()
			if height > 0 {
				width = int(float64(height) * aspectRatio)
			}
		}

		if height == 0 {
			height = i.image.Height()
			if width > 0 {
				height = int(float64(width) / aspectRatio)
			} else {
				width = i.image.Width()
			}
		}
	}

	err := i.image.ThumbnailWithSize(width, height, vips.InterestingNone, vips.SizeDown)
	if err != nil {
		return nil, err
	}

	if params.Fit == FitTypeCover {
		// calculate coordinates for cropping
		x := 0
		y := 0
		if i.image.Width() > params.Width {
			x = (i.image.Width() - params.Width) / 2
		}
		if i.image.Height()/i.image.Pages() > params.Height {
			y = (i.image.Height()/i.image.Pages() - params.Height) / 2
		}

		err = i.image.Crop(x, y, params.Width, params.Height)
		if err != nil {
			return nil, err
		}
	}

	return i, nil
}

func (i *vipsImage) GetInfo() (ImageMetadata, error) {
	if i.image == nil {
		return ImageMetadata{}, errors.New("image not loaded")
	}
	return ImageMetadata{
		Width:  i.image.Width(),
		Height: i.image.Height(),
		Format: vips.ImageTypes[i.image.Format()],
	}, nil
}

func (i *vipsImage) Export(params ImageExportParams) ([]byte, ImageMetadata, error) {
	if i.image == nil {
		return nil, ImageMetadata{}, errors.New("image not loaded")
	}
	nativeFormat := i.image.Format()
	if params.Format != "" {
		switch params.Format {
		case "jpeg", "jpg":
			nativeFormat = vips.ImageTypeJPEG
		case "png":
			nativeFormat = vips.ImageTypePNG
		case "webp":
			nativeFormat = vips.ImageTypeWEBP
		case "tiff":
			nativeFormat = vips.ImageTypeTIFF
		case "gif":
			nativeFormat = vips.ImageTypeGIF
		case "svg":
			nativeFormat = vips.ImageTypeSVG
		case "heif":
			nativeFormat = vips.ImageTypeHEIF
		case "pdf":
			nativeFormat = vips.ImageTypePDF
		case "avif":
			nativeFormat = vips.ImageTypeAVIF
		default:
			return nil, ImageMetadata{}, errors.New("invalid format")

		}
	}

	if params.Quality == 0 {
		params.Quality = 80
	}

	reader, metadata, err := i.image.Export(&vips.ExportParams{
		Quality:       params.Quality,
		Format:        nativeFormat,
		StripMetadata: params.StripMetadata,
	})

	// reader,err:

	if err != nil {
		return nil, ImageMetadata{}, err
	}
	return reader, ImageMetadata{
		Width:  metadata.Width,
		Height: metadata.Height,
		Format: vips.ImageTypes[metadata.Format],
	}, nil

}

type ImageFile struct {
	image Image
	file  io.ReadCloser
}

func NewImageFile(file io.ReadCloser) *ImageFile {
	return &ImageFile{
		file: file,
	}
}

func (i *ImageFile) Close() error {
	if i.image != nil {
		i.image.Close()
		i.image = nil
	}
	return i.file.Close()
}

func (i *ImageFile) Load() error {
	if i.image != nil {
		return nil
	}
	image, err := newVipsImageFromReader(i.file)
	if err != nil {
		return err
	}
	i.image = image
	return nil
}

func (i *ImageFile) Thumbnail(params ThumbnailParams) (Image, error) {
	err := i.Load()
	if err != nil {
		return nil, err
	}
	return i.image.Thumbnail(params)
}

func (i *ImageFile) GetInfo() (ImageMetadata, error) {
	err := i.Load()
	if err != nil {
		return ImageMetadata{}, err
	}
	return i.image.GetInfo()
}

func (i *ImageFile) Export(params ImageExportParams) ([]byte, ImageMetadata, error) {
	err := i.Load()
	if err != nil {
		return nil, ImageMetadata{}, err
	}
	return i.image.Export(params)
}

func (i *ImageFile) CopyImageRef() (Image, error) {
	err := i.Load()
	if err != nil {
		return nil, err
	}
	copy, err := i.image.CopyImageRef()
	if err != nil {
		return nil, err
	}

	return &ImageFile{
		image: copy,
		file:  i.file,
	}, nil
}
