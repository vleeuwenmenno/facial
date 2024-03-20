package insightfaceapi

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// APIWrapper class
type APIWrapper struct {
	host string
}

// NewAPIWrapper creates a new APIWrapper
func NewAPIWrapper(host string) *APIWrapper {
	return &APIWrapper{host: host}
}

// SearchResult search for a face in the database
func (wrapper *APIWrapper) SearchFace(imagePath string) (ApiResponse[SearchResult], error) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	fileContents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	response, err := sendPostRequest(wrapper.host+"/face-search", postContent{
		fname: imagePath,
		ftype: "file",
		fdata: fileContents,
	})
	if err != nil {
		return ApiResponse[SearchResult]{}, err
	}

	if response.StatusCode != 200 {
		return ApiResponse[SearchResult]{}, fmt.Errorf("error searching face: %s", response.Status)
	}

	defer response.Body.Close()
	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return ApiResponse[SearchResult]{}, err
	}

	var searchResult ApiResponse[SearchResult]
	err = json.Unmarshal(responseBodyBytes, &searchResult)
	if err != nil {
		return ApiResponse[SearchResult]{}, err
	}

	return searchResult, err
}

// VerifyImage verifies the image with the given name and returns similarity
func (wrapper *APIWrapper) VerifyImage(imagePath string, name string) (ApiResponse[FaceVerify], error) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	fileContents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	file.Close()

	response, err := sendPostRequest(wrapper.host+"/face-verification?name="+name, postContent{
		fname: imagePath,
		ftype: "file",
		fdata: fileContents,
	})
	if err != nil {
		return ApiResponse[FaceVerify]{}, err
	}

	if response.StatusCode != 200 {
		return ApiResponse[FaceVerify]{}, fmt.Errorf("error uploading selfie: %s", response.Status)
	}

	defer response.Body.Close()
	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return ApiResponse[FaceVerify]{}, err
	}

	var uploadSelfieResult ApiResponse[FaceVerify]
	err = json.Unmarshal(responseBodyBytes, &uploadSelfieResult)
	if err != nil {
		return ApiResponse[FaceVerify]{}, err
	}
	return uploadSelfieResult, err
}

// AddSelfie uploads a selfie to the database
func (wrapper *APIWrapper) AddSelfie(imagePath string, name string) (ApiResponse[Face], error) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	fileContents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	file.Close()

	response, err := sendPostRequest(wrapper.host+"/upload-selfie?name="+name, postContent{
		fname: imagePath,
		ftype: "file",
		fdata: fileContents,
	})
	if err != nil {
		return ApiResponse[Face]{}, err
	}

	if response.StatusCode != 200 {
		return ApiResponse[Face]{}, fmt.Errorf("error uploading selfie: %s", response.Status)
	}

	defer response.Body.Close()
	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return ApiResponse[Face]{}, err
	}

	var uploadSelfieResult ApiResponse[Face]
	err = json.Unmarshal(responseBodyBytes, &uploadSelfieResult)
	if err != nil {
		return ApiResponse[Face]{}, err
	}
	return uploadSelfieResult, err
}

// ListFaces returns a list of faces from the database
func (wrapper *APIWrapper) ListFaces(page int, limit int) (ApiResponse[GetFacesResult], error) {
	response, err := sendGetRequest(wrapper.host + "/faces?page=" + fmt.Sprint(page) + "&page_size=" + fmt.Sprint(limit))
	if err != nil {
		return ApiResponse[GetFacesResult]{}, err
	}

	if response.StatusCode != 200 {
		return ApiResponse[GetFacesResult]{}, fmt.Errorf("error getting faces: %s", response.Status)
	}

	defer response.Body.Close()
	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return ApiResponse[GetFacesResult]{}, err
	}

	var getFacesResult ApiResponse[GetFacesResult]
	err = json.Unmarshal(responseBodyBytes, &getFacesResult)
	if err != nil {
		return ApiResponse[GetFacesResult]{}, err
	}
	return getFacesResult, err
}

// UpdateFace updates the face with the given ID. If name is not empty, it will also update the name of the face.
func (wrapper *APIWrapper) UpdateFace(id string, name string, imagePath string) (ApiResponse[FaceUpload], error) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	fileContents, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	file.Close()

	url := fmt.Sprintf(wrapper.host+"/faces?id=%s", id)
	if name != "" {
		url += "&name=" + name
	}

	response, err := sendPostRequest(url, postContent{
		fname: imagePath,
		ftype: "file",
		fdata: fileContents,
	})
	if err != nil {
		return ApiResponse[FaceUpload]{}, err
	}

	if response.StatusCode != 200 {
		return ApiResponse[FaceUpload]{}, fmt.Errorf("error updating face: %s", response.Status)
	}

	defer response.Body.Close()
	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return ApiResponse[FaceUpload]{}, err
	}

	var uploadSelfieResult ApiResponse[FaceUpload]
	err = json.Unmarshal(responseBodyBytes, &uploadSelfieResult)
	if err != nil {
		return ApiResponse[FaceUpload]{}, err
	}
	return uploadSelfieResult, err
}

// UpdateFaceImage updates the face with the given name.
func (wrapper *APIWrapper) UpdateFaceImage(id string, imagePath string) (ApiResponse[FaceUpload], error) {
	return wrapper.UpdateFace(id, "", imagePath)
}

// DeleteFace deletes the face with the given ID from the database
func (wrapper *APIWrapper) DeleteFace(name string) error {
	response, err := sendDeleteRequest(wrapper.host + "/faces?name=" + fmt.Sprint(name))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("error deleting face: %s", response.Status)
	}

	return nil
}
