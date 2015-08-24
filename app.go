// Used to create a large number of suppliers at once, given no import feature in Vend.
// Takes a CSV file input, then creates each supplier with a POST to the /api/supplier endpoint.
package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Create a log file.
	lf, err := os.Create("./vend_supplier_bulk_add.log")
	if err != nil {
		// Can't log the error, so just panic.
		fmt.Println("Couldn't create logfile. Where are we running?")
		panic(0)
	}
	// Ensure log file is closed at the end.
	defer lf.Close()
	log.SetOutput(lf)

	// Define command line flags.
	var authToken, domainPrefix, filePath string
	flag.StringVar(&authToken, "t", "", "Vend API token. (Vend -> Setup -> API Access)")
	flag.StringVar(&domainPrefix, "d", "", "Vend Store name (xxxx.vendhq.com)")
	flag.StringVar(&filePath, "f", "", "Path to CSV file of suppliers.")
	flag.Parse()

	// Log values of all flags.
	log.Printf("Command line arguments provided: -t %s -d %s -f %s",
		authToken, domainPrefix, filePath)

	// Save people who write domain_prefix.vendhq.com.
	if strings.HasSuffix(domainPrefix, ".vendhq.com") {
		domainPrefix = domainPrefix[:len(domainPrefix)-11]
	}

	// Read supplier list from CSV file.
	supplierList, err := readSupplierCSV(filePath)
	if err != nil {
		log.Fatalf("Something went wrong trying to read supplier CSV: %s", err)
	}
	// Post all suppliers to Vend.
	err = postAllSuppliers(supplierList, domainPrefix, authToken)
	if err != nil {
		log.Fatalf("Something went wrong trying to post suppliers: %s", err)
	}
}

// Reads supplier CSV, returns a slice of supplier types.
func readSupplierCSV(filePath string) ([]Supplier, error) {

	exampleHeader := []string{"name", "description", "first_name", "last_name", "company_name",
		"phone", "mobile", "fax", "email", "twitter", "website", "physical_address1",
		"physical_address2", "physical_suburb", "physical_city",
		"physical_postcode", "physical_state", "physical_country_id",
		"postal_address1", "postal_address2", "postal_suburb", "postal_city",
		"postal_postcode", "postal_state", "postal_country_id"}

	// Open our provided CSV file.
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Printf("Could not read from CSV file: %s", err)
		return nil, err
	}
	// Make sure to close at end.
	defer csvFile.Close()

	// Create CSV reader on our file.
	reader := csv.NewReader(csvFile)

	var supplier Supplier
	var supplierList []Supplier

	// Read and store our header line.
	headerRow, err := reader.Read()

	// Check each header in the row is same as our template.
	for i := range headerRow {
		if headerRow[i] != exampleHeader[i] {
			fmt.Println("Found error in header row. Check log!")
			log.Fatalf("No header match for: %s, instead got: %s.",
				string(exampleHeader[i]), string(headerRow[i]))
		}
	}

	// Read the rest of the data from the CSV.
	rawData, err := reader.ReadAll()

	// Loop through rows and assign them to supplier type.
	for _, row := range rawData {

		supplier.Name = row[0]
		supplier.Description = row[1]
		supplier.Contact.FirstName = row[2]
		supplier.Contact.LastName = row[3]
		supplier.Contact.CompanyName = row[4]
		supplier.Contact.Phone = row[5]
		supplier.Contact.Mobile = row[6]
		supplier.Contact.Fax = row[7]
		supplier.Contact.Email = row[8]
		supplier.Contact.Twitter = row[9]
		supplier.Contact.Website = row[10]
		supplier.Contact.PhysicalAddress1 = row[11]
		supplier.Contact.PhysicalAddress2 = row[12]
		supplier.Contact.PhysicalSuburb = row[13]
		supplier.Contact.PhysicalCity = row[14]
		supplier.Contact.PhysicalPostcode = row[15]
		supplier.Contact.PhysicalState = row[16]
		supplier.Contact.PhysicalCountryID = row[17]
		supplier.Contact.PostalAddress1 = row[18]
		supplier.Contact.PostalAddress2 = row[19]
		supplier.Contact.PostalSuburb = row[20]
		supplier.Contact.PostalCity = row[21]
		supplier.Contact.PostalPostcode = row[22]
		supplier.Contact.PostalState = row[23]
		supplier.Contact.PostalCountryID = row[24]

		// Append each supplier type to our slice of all suppliers.
		supplierList = append(supplierList, supplier)
	}

	return supplierList, err
}

// Iterate through and post each supplier.
func postAllSuppliers(supplierList []Supplier, domainPrefix, authToken string) error {
	var err error

	for _, supplier := range supplierList {

		supplierJSON, err := createSupplierJSON(supplier)
		if err != nil {
			log.Printf("Something went wrong trying to create supplier JSON: %s", err)
		}

		fmt.Printf("\n\nPosting body:\n %s", string(supplierJSON))

		err = postSupplier(supplierJSON, domainPrefix, authToken)
		if err != nil {
			log.Printf("Something went wrong trying to post supplier: %s", err)
		}
	}

	return err
}

// Create JSON payload ready for POST.
func createSupplierJSON(supplier Supplier) ([]byte, error) {

	// Creates pretty JSON format of our supplier.
	supplierJSON, err := json.MarshalIndent(supplier, "", "\t")
	if err != nil {
		log.Printf("Error marhsalling supplier JSON: %s", err)
		return nil, err
	}

	return supplierJSON, err
}

// Posts supplier JSON payload to supplier endpoint which creates supplier.
func postSupplier(supplier []byte, domainPrefix, authToken string) error {

	client := &http.Client{}

	url := fmt.Sprintf("https://%s.vendhq.com/api/supplier", domainPrefix)
	log.Printf("Posting supplier to endpoint: %s", url)

	// Supplier endpoint 0.X takes a POST request for supplier creation.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(supplier))
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return err
	}

	// Set request headers.
	// Auth token is created in Vend -> Setup -> API Access.
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Set("User-Agent", "Support-tool: supplier-bulk-add - one of JOEYM8's tools.")

	// Perform request.
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to post supplier: %s", err)
		return err
	}

	// Check HTTP response status codes, return error if code is unexpected.
	switch resp.StatusCode {
	case 200:
		log.Printf("Status %d - supplier birthed :D", resp.StatusCode)
	case 401:
		log.Printf("Access denied - check token. Status: %d", resp.StatusCode)
		return err
	case 404:
		log.Printf("URL not found - check domain prefix. Status: %d", resp.StatusCode)
		return err
	case 429:
		log.Printf("Our rates have been limited - sit tight. Status: %d", resp.StatusCode)
	default:
		log.Printf("Got an unknown status code - Google it. Status: %d", resp.StatusCode)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading response body: %s\n", err)
		return err
	}
	fmt.Printf("\n\nGot response:\n")
	os.Stdout.Write(body)

	return err
}
