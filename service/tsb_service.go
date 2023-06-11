package TSBService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pbnjay/grate"
	_ "github.com/pbnjay/grate/xls"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"tsb/config"
	"tsb/helper"
	"tsb/models"
	"tsb/schema"
	response_trait "tsb/trait"
)

var DB = config.DBConnect()
var sheetName = "Smarka"

var filePath = helper.GetEnv("TSB_FILE_PATH")

// GetData function is the main function to handle incoming requests, validate input, fetch and process data
func GetData(ctx *fiber.Ctx) error {
	var input schema.RequestModel

	// Parse the body of the request into the input variable
	err := ctx.BodyParser(&input)

	if err != nil {
		errors := map[string]string{"message": "Invalid request"}
		return response_trait.BadRequest(ctx, errors)
	}

	validate := validator.New()

	// Validate the parsed input
	validationError := validate.Struct(&input)

	if validationError != nil {
		errors := map[string]string{"message": validationError.(validator.ValidationErrors).Error()}
		return response_trait.UnprocessableEntity(ctx, errors)
	}

	// Send HTTP request with validated input
	res, err := sendRequest(input)

	if err != nil {
		return err
	}

	var response schema.TSBModel

	err = json.Unmarshal(res, &response)

	if err != nil {
		errors := map[string]string{"message": err.Error()}
		return response_trait.NotFound(ctx, errors)
	}

	// Download the file mentioned in the response
	err = getActualFile(&response)

	if err != nil {
		errors := map[string]string{"message": err.Error()}
		return response_trait.NotFound(ctx, errors)
	}

	// Run ReadFile function concurrently to read and process the downloaded file
	go ReadFile(filePath)

	return response_trait.Success(ctx, map[string]string{"message": "TSB Excel reading started whole data will be saved to the Database in about 5 minutes"})
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

	// Read the entire response body
	responseBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fmt.Println("Successfully request to TSB API")

	return responseBody, nil
}

func getActualFile(res *schema.TSBModel) error {
	filePath = filePath + res.ActualFile.Name // Update the file path by appending the file name
	url := res.ActualFile.FilePath            // Retrieve the URL from where the file needs to be downloaded

	// Create a new file with the updated path
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	// Close the file after the function finishes
	defer file.Close()

	// Send an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	// Close the response body after the function finishes
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

	fmt.Printf("Successfully download file from %s\n Kasko Year: %s\n Kasko Month: %s\n", url, res.ActualFile.KaskoYear, res.ActualFile.KaskoMonth)

	return nil
}

func ReadFile(fileName string) {
	// Open the spreadsheet file
	file, err := grate.Open(fileName)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully open file")

	// Defer a function to close the spreadsheet file and remove it from the file system
	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}

		err := os.Remove(fileName)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Successfully remove file")
	}()

	// Get data from the specified sheet in the spreadsheet
	data, _ := file.Get(sheetName)
	fmt.Println("Successfully get data")

	// Initialize the row counter
	i := 0

	// Loop through the rows of the data
	for data.Next() {
		if i == 0 || i == 1 { // Skip the first two rows
			i++
			continue
		}

		// Convert the row data into a slice of strings
		row := data.Strings()

		// Prepare the data to be inserted into the database and get the brand and model IDs
		brandId, modelId := prepareData(row)
		fmt.Printf("Brand ID: %d - Model ID: %d\n", brandId, modelId)

		// Initialize the year counter
		l := 0

		// loop through the number of years
		for j := 6; j <= 16; j++ {
			year := time.Now().Year()
			if l != 0 {
				year = year - l
			}

			cascoValue, _ := strconv.ParseFloat(row[j], 64)

			cascoStruct := models.CascoValue{
				ModelId:       modelId,
				Casco:         cascoValue,
				ModelYear:     uint16(year),
				LastUpdatedAt: time.Now(),
			}

			upsertCasco(cascoStruct)

			l++ // Increment the year counter
		}

		i++ // Increment the row counter
	}
}

func prepareData(data []string) (uint32, uint32) {
	brandCode, _ := strconv.ParseUint(data[0], 10, 16)

	brand := models.Brand{
		BrandCode: uint16(brandCode),
		Name:      data[2],
	}

	brandId := upsertBrand(brand)

	model := models.Model{
		BrandId: brandId,
		Name:    data[3],
	}

	modelId := upsertModel(model)

	return brandId, modelId
}

func upsertBrand(brand models.Brand) uint32 {
	DB.FirstOrCreate(&brand, models.Brand{BrandCode: brand.BrandCode})

	return brand.ID
}

func upsertModel(model models.Model) uint32 {
	DB.FirstOrCreate(&model, models.Model{BrandId: model.BrandId, Name: model.Name})

	return model.ID
}

func upsertCasco(casco models.CascoValue) {
	DB.FirstOrCreate(
		&casco,
		models.CascoValue{
			ModelId:   casco.ModelId,
			ModelYear: casco.ModelYear,
			Casco:     casco.Casco,
		},
	)
}
