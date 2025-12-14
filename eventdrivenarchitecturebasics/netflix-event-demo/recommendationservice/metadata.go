package recommendationservice

// Simulated Video Metadata Store (like Netflix Video Catalog)
var videoCategoryMap = map[string]string{
	"67890": "sports",
	"11111": "movies",
	"22222": "news",
	"33333": "kids",
}

// ResolveCategory returns category for a video
func ResolveCategory(videoID string) string {
	if category, ok := videoCategoryMap[videoID]; ok {
		return category
	}
	return "unknown"
}
