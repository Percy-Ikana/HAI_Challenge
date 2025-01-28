# Running
- Ensure you have Go installed
- clone the repository
```bash
git clone https://github.com/Percy-Ikana/HAI_Challenge
cd HAI_Challenge
```
- Ensure that `input.xml` contains the desired XML data.
- Modify `config.json` if you want to change the mappings.
- Run the program
```bash
go run .
```

# Configuring the Mappings:

the structure of the config map is a list of JSON objects with the following fields: , see `config.json` or `expandedConfig.json` for examples.

`json_field`: the `string` name of the JSON field. 

`xml_fields`: a list `[]` of `strings` containing the XML fields that are used for this JSON field.

`transform`:  a `string` that denotes what transform function to use to determine the final JSON value, this must exist within the `switch` contained within `transform`, at present, only `concat_strings`, `to_int`, `age_from_dob`, and `""` (no transform) are supported without modification. Invalid transformations will default to `""`. 

`children`: an optional list of any children this element should have, children use the same 

format. For an example, see `expandedConfig.json`


# Modifying the input XML:

For new XML fields to be properly used for mapping/transformations, they must be added to the XML struct, as well as the `patientData` map within `transformPatient`. XML elements not referred to in the mappings config will simply remain unused. 
