/*
To create supplier via HTTP POST to endpoint:
/api/supplier

Example payload:
{
  "name": "Name",
  "description": "Description",
  "contact": {
    "first_name": "Foo",
    "last_name": "Bar",
    "company_name": "It's a bar",
    "phone": "08 333",
    "mobile": "021 333",
    "fax": "wat",
    "email": "email",
    "twitter": "twitter.com",
    "website": "example.com",
    "physical_address1": "1 place",
    "physical_address2": "or the other",
    "physical_suburb": "Remuera",
    "physical_city": "AKL",
    "physical_postcode": "0000",
    "physical_state": "Georgia",
    "physical_country_id": "USA",
    "postal_address1": "Norm",
    "postal_address2": "Street",
    "postal_suburb": "Bla",
    "postal_city": "PCity",
    "postal_postcode": "2222",
    "postal_state": "Georgina",
    "postal_country_id": "NZ"
  }
}
*/

// Package types contains various structs.
package main

// Supplier contains basic supplier info, with contact details nested.
type Supplier struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Contact     Contact `json:"contact,omitempty"`
}

// Contact contains the supplier's specific contact info from the vend_contact table.
type Contact struct {
	FirstName         string `json:"first_name,omitempty"`
	LastName          string `json:"last_name,omitempty"`
	CompanyName       string `json:"company_name,omitempty"`
	Phone             string `json:"phone,omitempty"`
	Mobile            string `json:"mobile,omitempty"`
	Fax               string `json:"fax,omitempty"`
	Email             string `json:"email,omitempty"`
	Twitter           string `json:"twitter,omitempty"`
	Website           string `json:"website,omitempty"`
	PhysicalAddress1  string `json:"physical_address1,omitempty"`
	PhysicalAddress2  string `json:"physical_address2,omitempty"`
	PhysicalSuburb    string `json:"physical_suburb,omitempty"`
	PhysicalCity      string `json:"physical_city,omitempty"`
	PhysicalPostcode  string `json:"physicaL_postcode,omitempty"`
	PhysicalState     string `json:"physical_state,omitempty"`
	PhysicalCountryID string `json:"physical_country_id,omitempty"`
	PostalAddress1    string `json:"postal_address1,omitempty"`
	PostalAddress2    string `json:"postal_address2,omitempty"`
	PostalSuburb      string `json:"postal_suburb,omitempty"`
	PostalCity        string `json:"postal_city,omitempty"`
	PostalPostcode    string `json:"postal_postcode,omitempty"`
	PostalState       string `json:"postal_state,omitempty"`
	PostalCountryID   string `json:"postal_country_id,omitempty"`
}
