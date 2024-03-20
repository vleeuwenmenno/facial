package insightfaceapi

type SearchResult struct {
	SimilarFaces []FaceSearch `json:"similar_faces"`
}

type Face struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Embedding string `json:"embedding"`
}

type FaceSearch struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Age        int     `json:"age"`
	Gender     string  `json:"gender"`
	Embedding  string  `json:"embedding"`
	Similarity float64 `json:"similarity"`
}

type FaceUpload struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Embedding string `json:"embedding"`
}

type FaceVerify struct {
	Similarity float64 `json:"similarity"`
	Status     bool    `json:"status"`
}

type ApiResponse[T any] struct {
	StatusCode int `json:"status_code"`
	Result     T   `json:"result"`
}

type GetFacesResult struct {
	Faces []Face
}
