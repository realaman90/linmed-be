package server

import (
	"io"
	"net/http"
)

func (s *Server) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Extract the product or profile picture type from the request
	imageType := r.URL.Query().Get("type") // "product" or "profile"
	if imageType == "" {
		s.Logger.Error("image type not provided")
		errorResposne(w, http.StatusBadRequest, "image type must be provided")
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Limit your file size to 10 MB
	if err != nil {
		s.Logger.Error("file size exceeding 10 MB")
		errorResposne(w, http.StatusBadRequest, "file size exceeding 10 MB")
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("image")
	if err != nil {
		s.Logger.Error("image not found")
		errorResposne(w, http.StatusBadRequest, "image not found")
		return
	}
	defer file.Close()

	// Construct the S3 key based on the image type
	var folder string
	if imageType == "product" {
		folder = "products"
	} else if imageType == "profile" {
		folder = "profiles"
	} else {
		s.Logger.Error("invalid image type")
		errorResposne(w, http.StatusBadRequest, "invalid image type")
		return
	}

	// Use the original file name for the S3 key
	imageName := r.FormValue("name") // Get the image name from the form
	if imageName == "" {
		s.Logger.Error("image name not provided")
		errorResposne(w, http.StatusBadRequest, "image name must be provided")
		return
	}

	// Construct the full S3 key
	imageKey := folder + "/" + imageName

	// Call the existing AWS upload method
	err = s.S3Client.UploadImage(s.Config.BucketName, file, imageKey) // Use the constructed key
	if err != nil {
		s.Logger.Error("unable to upload", err.Error())
		errorResposne(w, http.StatusInternalServerError, "unable to upload to s3")
		return
	}

	// Respond with success
	writeJSONResponse(w, http.StatusAccepted, "image uploaded successfully")
}

func (s *Server) GetImage(w http.ResponseWriter, r *http.Request) {
	// Extract the product or profile picture type and image name from the request
	imageType := r.URL.Query().Get("type") // "product" or "profile"
	imageName := r.URL.Query().Get("name") // Name of the image

	if imageType == "" || imageName == "" {
		s.Logger.Error("image type or name not provided")
		errorResposne(w, http.StatusBadRequest, "image type and name must be provided")
		return
	}

	// Construct the S3 key based on the image type
	var folder string
	if imageType == "product" {
		folder = "products"
	} else if imageType == "profile" {
		folder = "profiles"
	} else {
		s.Logger.Error("invalid image type")
		errorResposne(w, http.StatusBadRequest, "invalid image type")
		return
	}

	// Construct the full S3 key
	imageKey := folder + "/" + imageName

	// Call the S3 GetImage method
	imageBody, err := s.S3Client.GetImage(s.Config.BucketName, imageKey)
	if err != nil {
		s.Logger.Error("unable to retrieve image", err.Error())
		errorResposne(w, http.StatusInternalServerError, "unable to retrieve image")
		return
	}
	defer imageBody.Close()

	// Set the appropriate content type (assuming it's a JPEG image)
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)

	// Copy the image data to the response writer
	_, err = io.Copy(w, imageBody)
	if err != nil {
		s.Logger.Error("unable to write image to response", err.Error())
		errorResposne(w, http.StatusInternalServerError, "unable to write image to response")
		return
	}
}
