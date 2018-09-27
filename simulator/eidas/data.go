package eidas

/*
import (
	"encoding/json"
)
*/

type Atributes struct {
	GetAttributes struct {
		Attributes []struct {
			Name  string `json:"Name"`
			Value string `json:"Value"`
		} `json:"Attributes"`
	} `json:"Get_Attributes"`
}

func GetDataSetByName(name string, server string) []byte {
	switch name {
	case "Representative Attributes":
		return Representative(server)
	case "Person Informacion":
		return PersonInformacion()
	}
	return nil
}

func Representative(server string) []byte {

	in := []byte(`
	{
  "Get_Attributes":{
    "Attributes":[
      {"Name": "eidasconnector", "Value": "http://` + server + `/SpecificConnector/ServiceProvider"},
      {"Name": "nodeMetadataUrl", "Value": "http://` + server + `/SpecificConnector/ServiceProvider"},
      {"Name":"citizenEidas", "Value": "CA"},
      {"Name":"returnUrl", "Value": "http://` + server + `/SP/ReturnPage"},
      {"Name":"eidasNameIdentifier", "Value": "unspecified"},
      {"Name":"eidasloa", "Value": "A"},
      {"Name":"eidasloaCompareType", "Value": "minimum"},
      {"Name":"eidasSPType", "Value": "public"},
      {"Name":"BirthName", "Value": "BirthName"},
      {"Name":"BirthNameType", "Value": "false"},
      {"Name":"CurrentAddress", "Value": "CurrentAddress"},
      {"Name":"CurrentAddressType", "Value": "false"},
      {"Name":"FamilyName", "Value": "FamilyName"},
      {"Name":"FamilyNameType", "Value": "true"},
      {"Name":"FirstName", "Value": "FirstName"},
      {"Name":"FirstNameType", "Value": "true"},
      {"Name":"DateOfBirth", "Value": "DateOfBirth"},
      {"Name":"DateOfBirthType", "Value": "true"},
      {"Name":"Gender", "Value": "Gender"},
      {"Name":"GenderType", "Value": "false"},
      {"Name":"PersonIdentifier", "Value": "PersonIdentifier"},
      {"Name":"PersonIdentifierType", "Value": "true"},
      {"Name":"PlaceOfBirth", "Value": "PlaceOfBirth"},
      {"Name":"PlaceOfBirthType", "Value": "false"},
      {"Name":"AdditionalAttribute", "Value": "AdditionalAttribute"},
      {"Name":"AdditionalAttributeType", "Value": "false"},
      {"Name":"D-2012-17-EUIdentifier", "Value": "D-2012-17-EUIdentifier"},
      {"Name":"D-2012-17-EUIdentifierType", "Value": "false"},
      {"Name":"EORI", "Value": "EORI"},
      {"Name":"EORIType", "Value": "false"},
      {"Name":"LEI", "Value": "LEI"},
      {"Name":"LEIType", "Value": "false"},
      {"Name":"LegalName", "Value": "LegalName"},
      {"Name":"LegalNameType", "Value": "true"},
      {"Name":"LegalAddress", "Value": "LegalAddress"},
      {"Name":"LegalAddressType", "Value": "false"},
      {"Name":"LegalPersonIdentifier", "Value": "LegalPersonIdentifier"},
      {"Name":"LegalPersonIdentifierType", "Value": "true"},
      {"Name":"SEED", "Value": "SEED"},
      {"Name":"SEEDType", "Value": "false"},
      {"Name":"SIC", "Value": "SIC"},
      {"Name":"SICType", "Value": "false"},
      {"Name":"TaxReference", "Value": "TaxReference"},
      {"Name":"TaxReferenceType", "Value": "false"},
      {"Name":"VATRegistration", "Value": "VATRegistration"},
      {"Name":"VATRegistrationType", "Value": "false"},
      {"Name":"LegalAdditionalAttribute", "Value": "LegalAdditionalAttribute"},
      {"Name":"LegalAdditionalAttributeType", "Value": "false"},
      {"Name":"allTypeEidas", "Value": "none"},
      {"Name":"RepresentativeBirthName", "Value": "RepresentativeBirthName"},
      {"Name":"RepresentativeBirthNameType", "Value": "none"},
      {"Name":"RepresentativeCurrentAddress", "Value": "RepresentativeCurrentAddress"},
      {"Name":"RepresentativeCurrentAddressType", "Value": "none"},
      {"Name":"RepresentativeFamilyName", "Value": "RepresentativeFamilyName"},
      {"Name":"RepresentativeFamilyNameType", "Value": "none"},
      {"Name":"RepresentativeFirstName", "Value": "RepresentativeFirstName"},
      {"Name":"RepresentativeFirstNameType", "Value": "none"},
      {"Name":"RepresentativeDateOfBirth", "Value": "RepresentativeDateOfBirth"},
      {"Name":"RepresentativeDateOfBirthType", "Value": "none"},
      {"Name":"RepresentativeGender", "Value": "RepresentativeGender"},
      {"Name":"RepresentativeGenderType", "Value": "none"},
      {"Name":"RepresentativePersonIdentifier", "Value": "RepresentativePersonIdentifier"},
      {"Name":"RepresentativePersonIdentifierType", "Value": "none"},
      {"Name":"RepresentativePlaceOfBirth", "Value": "RepresentativePlaceOfBirth"},
      {"Name":"RepresentativePlaceOfBirthType", "Value": "none"},
      {"Name":"RepresentativeD-2012-17-EUIdentifier", "Value": "RepresentativeD-2012-17-EUIdentifier"},
      {"Name":"RepresentativeD-2012-17-EUIdentifierType", "Value": "none"},
      {"Name":"RepresentativeEORI", "Value": "RepresentativeEORI"},
      {"Name":"RepresentativeEORIType", "Value": "none"},
      {"Name":"RepresentativeLEI", "Value": "RepresentativeLEI"},
      {"Name":"RepresentativeLEIType", "Value": "none"},
      {"Name":"RepresentativeLegalName", "Value": "RepresentativeLegalName"},
      {"Name":"RepresentativeLegalNameType", "Value": "none"},
      {"Name":"RepresentativeLegalAddress", "Value": "RepresentativeLegalAddress"},
      {"Name":"RepresentativeLegalAddressType", "Value": "none"},
      {"Name":"RepresentativeLegalPersonIdentifier", "Value": "RepresentativeLegalPersonIdentifier"},
      {"Name":"RepresentativeLegalPersonIdentifierType", "Value": "none"},
      {"Name":"RepresentativeSEED", "Value": "RepresentativeSEED"},
      {"Name":"RepresentativeSEEDType", "Value": "none"},
      {"Name":"RepresentativeSIC", "Value": "RepresentativeSIC"},
      {"Name":"RepresentativeSICType", "Value": "none"},
      {"Name":"RepresentativeTaxReference", "Value": "RepresentativeTaxReference"},
      {"Name":"RepresentativeTaxReferenceType", "Value": "none"},
      {"Name":"RepresentativeVATRegistration", "Value": "RepresentativeVATRegistration"},
      {"Name":"RepresentativeVATRegistrationType", "Value": "none"},
      {"Name":"spType", "Value":"public"}
  ]
  }
}`)

	return in
}

func PersonInformacion() []byte {
	in := []byte(
		`{
		  "Get_Attributes": {
		    "Attributes": [
		      {
		        "Name": "requestId",
		        "Value": ""
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/naturalperson/CurrentFamilyName",
		        "Value": "http://eidas.europa.eu/attributes/naturalperson/CurrentFamilyName"
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/naturalperson/CurrentGivenName",
		        "Value": "http://eidas.europa.eu/attributes/naturalperson/CurrentGivenName"
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/naturalperson/DateOfBirth",
		        "Value": "http://eidas.europa.eu/attributes/naturalperson/DateOfBirth"
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/naturalperson/PersonIdentifier",
		        "Value": "http://eidas.europa.eu/attributes/naturalperson/PersonIdentifier"
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/legalperson/LegalName",
		        "Value": "http://eidas.europa.eu/attributes/legalperson/LegalName"
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/legalperson/LegalPersonIdentifier",
		        "Value": "http://eidas.europa.eu/attributes/legalperson/LegalPersonIdentifier"
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/naturalperson/BirthName",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/naturalperson/BirthName",
		        "Value": "true"
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/naturalperson/CurrentAddress",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/naturalperson/CurrentAddress",
		        "Value": "true"
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/naturalperson/Gender",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/naturalperson/Gender",
		        "Value": "true"
		      },
		      {
		        "Name": "http://eidas.europa.eu/attributes/naturalperson/PlaceOfBirth",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/naturalperson/PlaceOfBirth",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/naturalperson/AdditionalAttribute",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/legalperson/D-2012-17-EUIdentifier",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/legalperson/EORI",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/legalperson/LEI",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/legalperson/LegalPersonAddress",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/legalperson/SEED",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/legalperson/SIC",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/legalperson/TaxReference",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/legalperson/VATRegistrationNumber",
		        "Value": "true"
		      },
		      {
		        "Name": "__checkbox_http://eidas.europa.eu/attributes/legalperson/LegalAdditionalAttribute",
		        "Value": "true"
		      }
		    ]
			  }
		}`)
	return in
}
