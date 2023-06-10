package TSBService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"os"
	"tsb/helper"
	"tsb/schema"
	response_trait "tsb/trait"
)

var filePath = helper.GetEnv("TSB_FILE_PATH")

func GetData(ctx *fiber.Ctx) error {
	var input schema.RequestModel

	err := ctx.BodyParser(&input)

	if err != nil {
		errors := map[string]string{"message": "Invalid request"}
		return response_trait.BadRequest(ctx, errors)
	}

	validate := validator.New()

	validationError := validate.Struct(&input)

	if validationError != nil {
		errors := map[string]string{"message": validationError.(validator.ValidationErrors).Error()}
		return response_trait.UnprocessableEntity(ctx, errors)
	}

	res, err := sendRequest(input)

	if err != nil {
		return err
	}

	var response schema.TSBModel

	err = json.Unmarshal([]byte(res), &response)

	if err != nil {
		errors := map[string]string{"message": err.Error()}
		return response_trait.NotFound(ctx, errors)
	}

	err = getActualFile(&response)

	return err
}

func sendRequest(Input schema.RequestModel) ([]byte, error) {
	payload, err := json.Marshal(Input)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := http.Post(helper.GetEnv("TSB_URL"), "application/json", bytes.NewBufferString(string(payload)))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println("Successfully request to TSB API")

	return responseBody, nil
}

func getActualFile(res *schema.TSBModel) error {
	filePath = filePath + res.ActualFile.Name
	url := res.ActualFile.FilePath

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer file.Close()

	// Send an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer response.Body.Close()

	// Check if the response was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error response: %s", response.Status)
	}

	// Copy the response body to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	fmt.Printf("Successfully download file from %s\n Kakso Year: %s\n Kasko Month: %s\n", url, res.ActualFile.KaskoYear, res.ActualFile.KaskoMonth)

	return nil
}
