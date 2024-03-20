package main

import (
	"flag"
	"fmt"
	insightfaceapi "insightface-api-wrapper/insightface-api"
	"os"
)

func main() {
	addFace := flag.Bool("add", false, "Add a face")
	verifyFace := flag.Bool("verify", false, "Verify a face vs the given name")
	search := flag.Bool("search", false, "Search for a face")
	list := flag.Bool("list", false, "List all faces")

	name := flag.String("name", "", "Provide a name of the person in the image")
	delete := flag.String("delete", "", "Delete a face")
	update := flag.String("update", "", "Update a face with the given ID")
	file := flag.String("file", "", "File path to upload")
	host := flag.String("host", "http://localhost:1800", "Host of the API")

	page := flag.Int("page", 1, "Page number")
	perPage := flag.Int("per-page", 10, "Number of items per page")

	flag.Parse()

	api := insightfaceapi.NewAPIWrapper(*host)

	if *file != "" && *name != "" && *addFace {
		fmt.Printf("Uploading selfie for %s ...\n", *name)

		result, err := api.AddSelfie(*file, *name)
		if err != nil {
			fmt.Println("Error uploading selfie:", err)
			os.Exit(1)
		}
		face := result.Result

		fmt.Printf("%d : %s (Age: %d)\n", face.ID, face.Name, face.Age)
		return
	}

	if *file != "" && *name != "" && *verifyFace {
		fmt.Printf("Uploading verification for %s ...\n", *name)

		result, err := api.VerifyImage(*file, *name)
		if err != nil {
			fmt.Println("Error uploading verification:", err)
			os.Exit(1)
		}

		fmt.Printf("Similarity: %f\n", result.Result.Similarity)
	}

	if *file != "" && *search {
		fmt.Println("Searching for face ...")

		result, err := api.SearchFace(*file)
		if err != nil {
			fmt.Println("Error searching face:", err)
			os.Exit(1)
		}

		fmt.Printf("Found %d faces ...\n", len(result.Result.SimilarFaces))
		for _, face := range result.Result.SimilarFaces {
			fmt.Printf("%d : %s (Age: %d)\n", face.ID, face.Name, face.Age)
		}
	}

	if *list {
		fmt.Println("Listing all faces ...")

		result, err := api.ListFaces(*page, *perPage)
		if err != nil {
			fmt.Println("Error listing faces:", err)
			os.Exit(1)
		}

		for _, face := range result.Result.Faces {
			fmt.Printf("%d : %s (Age: %d)\n", face.ID, face.Name, face.Age)
		}
	}

	if *delete != "" {
		fmt.Printf("Deleting face %s ...\n", *delete)

		err := api.DeleteFace(*delete)
		if err != nil {
			fmt.Println("Error deleting face:", err)
			os.Exit(1)
		}

		fmt.Println("Face deleted")
	}

	if *update != "" && *name != "" && *file != "" {
		fmt.Printf("Updating face %s ...\n", *update)

		result, err := api.UpdateFace(*update, *name, *file)
		if err != nil {
			fmt.Println("Error updating face:", err)
			os.Exit(1)
		}
		face := result.Result

		fmt.Println("Face updated")
		fmt.Printf("%s (Age: %d)\n", face.Name, face.Age)
	}
}
