package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Mappings []Mapping `json:"mappings"`
}

type Mapping struct {
	JSONField string    `json:"json_field"`
	XMLFields []string  `json:"xml_fields"`
	Transform string    `json:"transform"`
	Children  []Mapping `json:"children,omitempty"`
}

type Patient struct {
	ID          string `xml:"ID,attr"`
	FirstName   string `xml:"FirstName"`
	LastName    string `xml:"LastName"`
	DateOfBirth string `xml:"DateOfBirth"`
	//Address     Address `xml:"Address"`
}

//Above and below are things related to address,
//a field I used to test for if the xml data had some
//child fields added so they could be handled

//Address struct for XML unmarshalling
//type Address struct {
//	Street string `xml:"Street"`
//	City   string `xml:"City"`
//	State  string `xml:"State"`
//}

type Patients struct {
	Patients []Patient `xml:"Patient"`
}

func loadXML(filename string) (Patients, error) {
	var patientXML Patients
	data, err := os.ReadFile(filename)
	if err != nil {
		return patientXML, err
	}
	err = xml.Unmarshal(data, &patientXML)
	return patientXML, err
}

func loadConfig(filename string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

//Transformation Functions go here, currently only the two needed are implemented.

func calculateAge(dob string) int {
	layout := "2006-01-02"
	birthdate, err := time.Parse(layout, dob)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return -1
	}
	age := time.Now().Year() - birthdate.Year()
	if time.Now().YearDay() < birthdate.YearDay() {
		age--
	}
	return age
}

func concatStrings(fields []string) string {
	return strings.Join(fields, " ")
}

// End of Transform functions

func transform(mapping Mapping, patientData map[string]string) interface{} {
	xmlValues := make([]string, len(mapping.XMLFields))
	for i, xmlField := range mapping.XMLFields {
		xmlValues[i] = patientData[xmlField]
	}
	var result interface{}

	//new transformations can be added to this switch.

	switch mapping.Transform {
	case "to_int":
		convertedValue, err := strconv.Atoi(xmlValues[0])
		if err != nil {
			//error in conversion, just return the string.
			result = xmlValues[0]
		} else {
			result = convertedValue
		}
	case "concat_strings":
		result = concatStrings(xmlValues)
	case "age_from_dob":
		result = calculateAge(xmlValues[0])
	default:
		//by default, just write the XML to JSON, while handling any children of the JSON field recusively.
		if len(mapping.Children) > 0 {
			childResult := make(map[string]interface{})
			for _, childMapping := range mapping.Children {
				childResult[childMapping.JSONField] = transform(childMapping, patientData)
			}
			result = childResult
		} else {
			result = xmlValues[0]
		}
	}
	return result
}

func transformPatient(patient Patient, config Config) map[string]interface{} {
	patientData := map[string]string{
		"ID":          patient.ID,
		"FirstName":   patient.FirstName,
		"LastName":    patient.LastName,
		"DateOfBirth": patient.DateOfBirth,

		//The below is an example of how this can be extended to handle new elements with children.

		//"Street":      patient.Address.Street,
		//"City":        patient.Address.City,
		//"State":       patient.Address.State,
	}
	mappedPatient := make(map[string]interface{})
	for _, mapping := range config.Mappings {
		mappedPatient[mapping.JSONField] = transform(mapping, patientData)
	}
	return mappedPatient
}

func main() {
	//Loads patients from source
	patients, err := loadXML("input.xml")
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return
	}

	//Loads the mapping config
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	//Transforms the source XML into JSON based on the mappings.
	var result []map[string]interface{}
	for _, patient := range patients.Patients {
		result = append(result, transformPatient(patient, config))
	}

	//Takes the result above and turns it into json.
	jsonData, err := json.MarshalIndent(map[string]interface{}{"patients": result}, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	//Prints
	fmt.Println(string(jsonData))

	//Saves to disk
	err = os.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

}
