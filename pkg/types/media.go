package types

type MediaType string

const (
	MediaTypeVideo MediaType = "video"
	MediaTypeImage MediaType = "image"
)

type MediaDTO struct {
	ID   uint      `json:"id"`
	Type MediaType `json:"type"`
	URL  string    `json:"url"`
	// TODO: Thêm các trường metadata (duration, quality, ...)
}
