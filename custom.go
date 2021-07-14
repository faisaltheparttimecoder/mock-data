package main

import (
	"fmt"
	"github.com/icrowley/fake"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
	"syreclabs.com/go/faker"
)

type Skeleton struct {
	Custom []TableModel `yaml:"Custom"`
}

type TableModel struct {
	Schema string        `yaml:"Schema"`
	Table  string        `yaml:"Table"`
	Column []ColumnModel `yaml:"Column"`
}

type ColumnModel struct {
	Name      string   `yaml:"Name"`
	Type      string   `yaml:"Type"`
	UserData  []string `yaml:"UserData"`
	Realistic string   `yaml:"Realistic"`
}

var fakeRealisticMaps = func() map[string]interface{} {
	return map[string]interface{}{
		"AddressCity":                  faker.Address().City(),                    // => "North Dessie"
		"AddressStreetName":            faker.Address().StreetName(),              // => "Buckridge Lakes"
		"AddressStreetAddress":         faker.Address().StreetAddress(),           // => "586 Sylvester Turnpike"
		"AddressSecondaryAddress":      faker.Address().SecondaryAddress(),        // => "Apt. 411"
		"AddressBuildingNumber":        faker.Address().BuildingNumber(),          // => "754"
		"AddressPostcode":              faker.Address().Postcode(),                // => "31340"
		"AddressPostcodeByState":       faker.Address().PostcodeByState("US"),     // => "46511"
		"AddressZipCode":               faker.Address().ZipCode(),                 // => 49448-5835.
		"AddressZipCodeByState":        faker.Address().ZipCodeByState("US"),      // => 76073-4516.
		"AddressTimeZone":              faker.Address().TimeZone(),                // => "Asia/Taipei"
		"AddressCityPrefix":            faker.Address().CityPrefix(),              // => "East"
		"AddressCitySuffix":            faker.Address().CitySuffix(),              // => "town"
		"AddressStreetSuffix":          faker.Address().StreetSuffix(),            // => "Square"
		"AddressState":                 faker.Address().State(),                   // => "Maryland"
		"AddressStateAbbr":             faker.Address().StateAbbr(),               // => "IL"
		"AddressCountry":               faker.Address().Country(),                 // => "Uruguay"
		"AddressCountryCode":           faker.Address().CountryCode(),             // => "JP"
		"AddressContinent":             fake.Continent(),                          // => Africa
		"AddressLatitude":              faker.Address().Latitude(),                // => (float32) -38.811367
		"AddressLongitude":             faker.Address().Longitude(),               // => (float32) 89.2171
		"AddressString":                faker.Address().String(),                  // => "6071 Heaney Island Suite 553, Ebbaville Texas 37307"Address
		"AppName":                      faker.App().Name(),                        // => "Alphazap"
		"AppVersion":                   faker.App().Version(),                     // => "2.6.0"
		"AppAuthor":                    faker.App().Author(),                      // => "Dorian Shields"
		"AppString":                    faker.App().String(),                      // => "Tempsoft 4.51"
		"AvatarUrl":                    faker.Avatar().Url("jpg", 100, 200),       // => "http://robohash.org/NX34rZw7s0VFzgWY.jpg?size=100x200"
		"AvatarString":                 faker.Avatar().String(),                   // => "http://robohash.org/XRWjFigoImqdeDuA.png?size=300x300"
		"BitcoinAddress":               faker.Bitcoin().Address(),                 // => "1GpEKM5UvD4XDLMirpNLoDnRVrGutogMj2"
		"BitcoinString":                faker.Bitcoin().String(),                  // => 14HD2RPKotUN8KV4qujTZDqr84KSAmseQZ
		"BusinessCreditCardNumber":     faker.Business().CreditCardNumber(),       // => "1234-2121-1221-1211"
		"BusinessCreditCardExpiryDate": faker.Business().CreditCardExpiryDate(),   // => "2015-11-11"
		"BusinessCreditCardType":       faker.Business().CreditCardType(),         // => "mastercard"
		"CalendarDay":                  fake.Day(),                                // => 22
		"CalendarDate":                 faker.Date(),                              // => 2017-01-01
		"CalendarMonth":                fake.Month(),                              // => October
		"CalendarMonthNum":             fake.MonthNum(),                           // => 10
		"CalendarTime":                 faker.Time(),                              // => 00:10:11
		"CalendarWeekDay":              fake.WeekDay(),                            // => Monday
		"CalendarWeekDayShort":         fake.WeekDayShort(),                       // => Mon
		"CalendarWeekDayNumber":        fake.WeekdayNum(),                         // => 6
		"CalendarYear":                 fake.Year(2010, 2021),                     // => 2017
		"CodeIsbn10":                   faker.Code().Isbn10(),                     // => "048931033-8"
		"CodeIsbn13":                   faker.Code().Isbn13(),                     // => "391668236072-1"
		"CodeEan13":                    faker.Code().Ean13(),                      // => "7742864258656"
		"CodeEan8":                     faker.Code().Ean8(),                       // => "03079010"
		"CodeRut":                      faker.Code().Rut(),                        // => "14371602-3"
		"CodeAbn":                      faker.Code().Abn(),                        // => "57914951376"
		"CommerceColor":                faker.Commerce().Color(),                  // => "lime"
		"CommerceDepartment":           faker.Commerce().Department(),             // => "Electronics, Health & Baby"
		"CommerceProductName":          faker.Commerce().ProductName(),            // => "Ergonomic Granite Shoes"
		"CommercePrice":                faker.Commerce().Price(),                  // => (float32) 97.79
		"CompanyName":                  faker.Company().Name(),                    // => "Aufderhar LLC"
		"CompanySuffix":                faker.Company().Suffix(),                  // => "Inc"
		"CompanyCatchPhrase":           faker.Company().CatchPhrase(),             // => "Universal logistical artificial intelligence"
		"CompanyBs":                    faker.Company().Bs(),                      // => "engage distributed applications"
		"CompanyEin":                   faker.Company().Ein(),                     // => "58-6520513"
		"CompanyDunsNumber":            faker.Company().DunsNumber(),              // => "16-708-2968"
		"CompanyLogo":                  faker.Company().Logo(),                    // => "http://www.biz-logo.com/examples/015.gif"
		"CompanyString":                faker.Company().String(),                  // => Lind, Leuschke and Braun.
		"FinanceCreditCard":            faker.Finance().CreditCard(faker.CC_VISA), // => "4190418835414""":
		"HackerSaySomethingSmart":      faker.Hacker().SaySomethingSmart(),        // => "If we connect the bus, we can get to the XML microchip through the digital TCP sensor!"
		"HackerAbbreviation":           faker.Hacker().Abbreviation(),             // => "HTTP"
		"HackerAdjective":              faker.Hacker().Adjective(),                // => "cross-platform"
		"HackerNoun":                   faker.Hacker().Noun(),                     // => "interface"
		"HackerVerb":                   faker.Hacker().Verb(),                     // => "bypass"
		"HackerIngVerb":                faker.Hacker().IngVerb(),                  // => "parsing"
		"InternetEmail":                faker.Internet().Email(),                  // => "maritza@farrell.org"
		"InternetFreeEmail":            faker.Internet().FreeEmail(),              // => "sven_rice@hotmail.com"
		"InternetSafeEmail":            faker.Internet().SafeEmail(),              // => "theron.nikolaus@example.net"
		"InternetUserName":             faker.Internet().UserName(),               // => "micah_pfeffer"
		"InternetPassword":             faker.Internet().Password(8, 14),          // => "s5CzvVp6Ye"
		"InternetDomainName":           faker.Internet().DomainName(),             // => "rolfson.info"
		"InternetDomainWord":           faker.Internet().DomainWord(),             // => "heller"
		"InternetDomainSuffix":         faker.Internet().DomainSuffix(),           // => "net"
		"InternetMacAddress":           faker.Internet().MacAddress(),             // => "15:a9:83:29:76:26"
		"InternetIpV4Address":          faker.Internet().IpV4Address(),            // => "121.204.82.227"
		"InternetIpV6Address":          faker.Internet().IpV6Address(),            // => "c697:392f:6a0e:bf6d:77e1:714a:10ab:0dbc"
		"InternetUrl":                  faker.Internet().Url(),                    // => "http://sporerhamill.net/kyla.schmitt"
		"InternetSlug":                 faker.Internet().Slug(),                   // => "officiis-commodi"
		"LoremCharacter":               faker.Lorem().Character(),                 // => "c"
		"LoremCharacters":              faker.Lorem().Characters(17),              // => "wqFyJIrXYfVP7cL9M"
		"LoremWord":                    faker.Lorem().Word(),                      // => "veritatis"
		"LoremSentence":                faker.Lorem().Sentence(10),                // => "Et officia atque dolor deserunt quam harum in quibusdam est."
		"LoremString":                  faker.Lorem().String(),                    // => Molestiae provident similique animi illum iure dolorem.
		"OtherWords":                   fake.Words(),                              // <= deserunt aut dignissimos ut
		"OtherWord":                    fake.Word(),                               // <= repellendus
		"NameFullName":                 faker.Name().Name(),                       // => "Natasha Hartmann"
		"NameFullNameWithPrefix":       fake.FullNameWithPrefix(),                 // => "Mr. Natasha Hartmann"
		"NameFirstName":                faker.Name().FirstName(),                  // => "Carolina"
		"NameLastName":                 faker.Name().LastName(),                   // => "Kohler"
		"NameMaleFullName":             fake.MaleFullName(),                       // <= James Matthews
		"NameMaleFirstName":            fake.MaleFirstName(),                      // <= Raymond
		"NameMaleLastName":             fake.MaleLastName(),                       // <= Harrison
		"NameMaleFullNameWithPrefix":   fake.MaleFullNameWithPrefix(),             // <= Mr. Dr. Peter Woods
		"NameFemaleFullName":           fake.FemaleFullName(),                     // <= Sandra Ferguson I II III IV V MD DDS PhD DVM
		"NameFemaleFirstName":          fake.FemaleFirstName(),                    // <= Debra
		"NameFemaleLastName":           fake.FemaleLastName(),                     // <= Ramirez
		"NameFemaleFullNameWithPrefix": fake.FemaleFullNameWithPrefix(),           // <= Mrs. Ms. Miss Diana Brooks
		"NameGender":                   fake.Gender(),                             // => "Female", "Male"
		"NameGenderAbbrev":             fake.GenderAbbrev(),                       // => "f", "m"
		"NamePrefix":                   faker.Name().Prefix(),                     // => "Dr."
		"NameSuffix":                   faker.Name().Suffix(),                     // => "Jr."
		"NameTitle":                    faker.Name().Title(),                      // => "Chief Functionality Orchestrator"
		"NameString":                   faker.Name().String(),                     // String is an alias for Name.
		"NumberNumber":                 faker.Number().Number(5),                  // => "43202"
		"NumberNumberInt":              faker.Number().NumberInt(3),               // => 213
		"NumberNumberInt32":            faker.Number().NumberInt32(5),             // => 92938
		"NumberNumberInt64":            faker.Number().NumberInt64(19),            // => 1689541633257139096
		"NumberDecimal":                faker.Number().Decimal(8, 2),              // => "879420.60"
		"NumberDigit":                  faker.Number().Digit(),                    // => "7"
		"NumberHexadecimal":            faker.Number().Hexadecimal(4),             // => "e7f3"
		"NumberBetween":                faker.Number().Between(-100, 100),         // => "-47"
		"NumberPositive":               faker.Number().Positive(1000),             // => "3"
		"NumberNegative":               faker.Number().Negative(-1000),            // => "-16"
		"PhoneNumberPhoneNumber":       faker.PhoneNumber().PhoneNumber(),         // => "1-599-267-6597 x537"
		"PhoneNumberCellPhone":         faker.PhoneNumber().CellPhone(),           // => "+49-131-0003060"
		"PhoneNumberAreaCode":          RandomInt(200, 999),                       // => "903"
		"PhoneNumberExchangeCode":      RandomInt(200, 999),                       // => "574"
		"PhoneNumberSubscriberNumber":  faker.PhoneNumber().SubscriberNumber(4),   // => "1512"
		"PhoneNumberString":            faker.PhoneNumber().String(),              // => (929) 929-3019 x1771
		"TeamName":                     faker.Team().Name(),                       // => "Colorado cats"
		"TeamCreature":                 faker.Team().Creature(),                   // => "cats"
		"TeamState":                    faker.Team().State(),                      // => "Oregon"
		"TeamString":                   faker.Team().String(),                     // => North Carolina goats
	}
}

// Generate a YAML of the mock plan related to this table
func GenerateMockPlan() {
	Infof("Generating a skeleton YAML for the list of table provided")

	// Check if the tables provided exists
	whereClause := generateWhereClause()
	tableList := dbExtractTables(whereClause)

	// If there is any then extract the column and data type
	if len(tableList) > 0 {
		columns := columnExtractor(tableList)
		s := BuildSkeletonYaml(columns)
		RegisterSkeletonYamlToFile(s)
	} else {
		Warn("No table available to generate the table skeleton, closing the program")
	}
}

// Realistic data generator
func RealisticDataBuilder(keyMap string) interface{} {
	Debugf("Getting data from realistic map for the key: %s", keyMap)
	// All the fake keys available
	var fakeMaps = fakeRealisticMaps()

	// If the request is to give the data via the key
	// then send the data
	if !IsStringEmpty(keyMap) {
		data, ok := fakeMaps[keyMap]
		if !ok {
			Fatalf("The requested realistic request key \"%s\" doesn't exists or currently not implemented, check your yaml", keyMap)
		}
		return data
	}
	return ""
}

// Build the skeleton Yaml
func BuildSkeletonYaml(column []TableCollection) Skeleton {
	Debugf("Build skeleton yaml")
	var s Skeleton
	for _, v := range column {
		var t = &TableModel{
			Schema: v.Schema,
			Table:  v.Table,
		}
		for _, u := range v.Columns {
			var c = &ColumnModel{
				Name: u.Column,
				Type: u.Datatype,
			}
			t.Column = append(t.Column, *c)
		}
		s.Custom = append(s.Custom, *t)
	}
	return s
}

// Create a yaml file and write the contents onto the file
func RegisterSkeletonYamlToFile(skeleton Skeleton) {
	Debugf("Save the yaml skeleton to file")
	f := fmt.Sprintf("%s_skeleton_%s.yaml", programName, ExecutionTimestamp)

	// Marshal the content
	d, err := yaml.Marshal(&skeleton)
	if err != nil {
		Fatalf("Error when marshalling the yaml file, err: %v", err)
	}

	// Save to file
	err = WriteToFile(f, string(d))
	if err != nil {
		Fatalf("Error when saving the yaml to file %s, err: %v", f, err)
	}

	Infof("The YAML is saved to file: %s/%s", CurrentDir(), f)
}

// Load the custom configuration and mock data based on that configuration
func MockCustoms() {
	Infof("Loading the table using custom configuration")
	c := Skeleton{}

	// Read the configuration
	c.ReadConfiguration()

	// Start Loading the data
	c.LoadDataByConfiguration()

	// If the program skipped the tables lets the users know
	skipTablesWarning()
}

// Read Configuration and load it to skeleton struct
func (c *Skeleton) ReadConfiguration() {
	Debugf("Reading the configuration file: %s", cmdOptions.File)

	// Using viper reading the configuration
	viper.SetConfigFile(cmdOptions.File)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		Fatalf("Error reading the YAML file, err: %v", err)
	}

	// Load the configuration onto struct
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

// Load the data based on custom configuration
func (c *Skeleton) LoadDataByConfiguration() {
	Infof("Loading data to the table based on what is defined by file %s", cmdOptions.File)

	// Open db connection
	db := ConnectDB()
	defer db.Close()

	for _, s := range c.Custom {
		// Initialize the mocking process
		tab := GenerateTableName(s.Table, s.Schema)
		msg := fmt.Sprintf("Mocking Table %s", tab)
		bar := StartProgressBar(msg, cmdOptions.Rows)

		// Name the for loop to break when we encounter error
	DataTypePickerLoop:
		// Loop through the row count and start loading the data
		for i := 0; i < cmdOptions.Rows; i++ {
			var data []string
			var col []string

			// Column info
			for _, v := range s.Column {
				var d interface{}
				var err error
				if len(v.UserData) > 0 { // User asked to use only the one that is provided
					d = RandomPickerFromArray(v.UserData)
				} else if !IsStringEmpty(v.Realistic) { // User wants to pick from one of the supported realistic data set
					d = RealisticDataBuilder(v.Realistic)
				} else { // If the user said for this column choose anything
					d, err = BuildData(v.Type)
					if err != nil {
						if strings.HasPrefix(fmt.Sprint(err), "unsupported datatypes found") {
							Debugf("Table %s skipped: %v", tab, err)
							skippedTab = append(skippedTab, tab)
							bar.Add(cmdOptions.Rows)
							break DataTypePickerLoop
						} else {
							Fatalf("Error when building data for table %s: %v", tab, err)
						}
					}
				}
				col = append(col, v.Name)
				data = append(data, fmt.Sprintf("%v", d))
			}

			// Copy the data to the table
			CopyData(tab, col, data, db)
			bar.Add(1)
		}
	}
}
