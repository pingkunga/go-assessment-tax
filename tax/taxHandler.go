package tax

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

func CalculationsHandler(c echo.Context) error {
	var taxRequest TaxRequest
	if err := c.Bind(&taxRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := ValidateTaxRequest(taxRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("ValidateError: %v", err.Error()))
	}

	TaxResponse := CalculateTax(taxRequest)
	return c.JSON(http.StatusOK, TaxResponse)
}

func BatchCalculationsHandler(c echo.Context) error {
	// Multipart form
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("taxFile")
	if err != nil {
		return err
	}

	//https://stackoverflow.com/questions/48349927/how-to-write-a-directory-with-a-file-in-golang
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	//taxes.csv
	CheckBatchTxPath("Uploads")
	file.Filename = generateFileName(file.Filename)
	dst, err := os.Create(filepath.Join("Uploads", file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	//Process
	taxRequests, err := ImportTaxCSV(filepath.Join("Uploads", file.Filename))

	//Clear file
	_ = os.Remove(filepath.Join("Uploads", file.Filename))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("ImportError: %v", err.Error()))
	}

	//Calculate
	taxBatchsResponse, err := CalculateTaxBatch(taxRequests)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("CalculateError: %v", err.Error()))
	}

	return c.JSON(http.StatusOK, taxBatchsResponse)
}

func CheckBatchTxPath(p string) {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.Mkdir(p, 0770)
	}
}

func fileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func generateFileName(fileName string) string {
	return fmt.Sprintf("%s-%d%s", fileNameWithoutExt(fileName), time.Now().Unix(), filepath.Ext(fileName))
}
