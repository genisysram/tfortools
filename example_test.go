// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfortools

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func ExampleGenerateUsageDecorated() {
	cfg := NewConfig(OptCols)
	help := GenerateUsageDecorated("-f", []struct{ X, Y int }{}, cfg)
	fmt.Println(help)
	// output:
	// The template passed to the --f option operates on a
	//
	// []struct {
	//	X int
	//	Y int
	// }
	//
	// Some new functions have been added to Go's template language
	//
	// - 'cols' can be used to extract certain columns from a table consisting of a
	//   slice or array of structs.  It returns a new slice of structs which contain
	//   only the fields requested by the caller.   For example, given a slice of structs
	//
	//   {{cols . "Name" "Address"}}
	//
	//   returns a new slice of structs, each element of which is a structure with only
	//   two fields, 'Name' and 'Address'.
}

func ExampleGenerateUsageUndecorated() {
	i := struct {
		X       int    `tfortools:"This is an int"`
		Y       string `json:"omitempty" tfortools:"This is a string"`
		hidden  float64
		Invalid chan int
	}{}
	help := GenerateUsageUndecorated(i)
	fmt.Println(help)
	// output:
	// struct {
	// 	X int    // This is an int
	// 	Y string `json:"omitempty"` // This is a string
	// }
}

func ExampleTemplateFunctionNames() {
	cfg := NewConfig(OptCols, OptRows)
	err := cfg.AddCustomFn(strings.TrimSpace, "trim",
		"- trim trims leading and trailing whitespace from string")
	if err != nil {
		panic(err)
	}
	for _, fn := range TemplateFunctionNames(cfg) {
		fmt.Println(fn)
	}
	// output:
	// cols
	// rows
	// trim
}

func ExampleTemplateFunctionHelpSingle() {
	cfg := NewConfig(OptCols, OptRows)
	err := cfg.AddCustomFn(strings.TrimSpace, "trim",
		"- trim trims leading and trailing whitespace from string")
	if err != nil {
		panic(err)
	}
	help, err := TemplateFunctionHelpSingle("cols", cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(help)

	help, err = TemplateFunctionHelpSingle("trim", cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(help)
	// output:
	// - 'cols' can be used to extract certain columns from a table consisting of a
	//   slice or array of structs.  It returns a new slice of structs which contain
	//   only the fields requested by the caller.   For example, given a slice of structs
	//
	//   {{cols . "Name" "Address"}}
	//
	//   returns a new slice of structs, each element of which is a structure with only
	//   two fields, 'Name' and 'Address'.
	//
	// - trim trims leading and trailing whitespace from string
}

func ExampleOutputToTemplate() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// print the surname of the person whose middlename is lexographically smallest.
	script := `{{select (head (sort . "MiddleName")) "Surname"}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Caesar
}

func ExampleOptFilter() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Print the surname of all people whose first name is Marcus
	script := `{{range (filter . "FirstName" "Marcus")}}{{println .Surname}}{{end}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Cicero
	// Crassus
}

func ExampleOptFilterContains() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Count the number of people whose middle name contains a 'ul'
	script := `{{len (filterContains . "MiddleName" "ul")}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// 2
}

func ExampleOptFilterHasPrefix() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Print all the surnames that start with 'Ci'
	script := `{{select (filterHasPrefix . "Surname" "Ci") "Surname"}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Cicero
}

func ExampleOptFilterHasSuffix() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Print all the surnames that end with 'us'
	script := `{{select (filterHasSuffix . "Surname" "us") "Surname"}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Crassus
}

func ExampleOptFilterFolded() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Output the first and surnames of all people whose first name is marcus
	script := `{{range (filterFolded . "FirstName" "marcus")}}{{println .FirstName .Surname}}{{end}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Marcus Cicero
	// Marcus Crassus
}

func ExampleOptFilterRegexp() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Output the first and last names of all people whose middle name ends in 'ius' and whose
	// second letter is 'u'
	script := `{{range (filterRegexp . "MiddleName" "^.u.*ius$")}}{{println .FirstName .Surname}}{{end}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Marcus Cicero
	// Gaius Caesar
}

func ExampleOptToJSON() {
	data := []struct {
		Name       string
		AgeAtDeath int
		Battles    []string
	}{
		{"Caesar", 55, []string{"Battle of Alesia", "Battle of Dyrrhachium", "Battle of the Nile"}},
		{"Alexander", 32, []string{"Battle of Issus", "Battle of Gaugamela", "Battle of the Hydaspes"}},
	}

	script := `{{tojson .}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// [
	// 	{
	//		"Name": "Caesar",
	//		"AgeAtDeath": 55,
	//		"Battles": [
	//			"Battle of Alesia",
	//			"Battle of Dyrrhachium",
	//			"Battle of the Nile"
	//		]
	//	},
	//	{
	//		"Name": "Alexander",
	//		"AgeAtDeath": 32,
	//		"Battles": [
	//			"Battle of Issus",
	//			"Battle of Gaugamela",
	//			"Battle of the Hydaspes"
	//		]
	//	}
	// ]
}

func ExampleOptTableX() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Output the names of people in a nicely formatted table
	script := `{{tablex . 12 8 0}}`
	var b bytes.Buffer
	if err := OutputToTemplate(&b, "names", script, data, nil); err != nil {
		panic(err)
	}

	// Normally you would pass os.Stdout directly into OutputToTemplate.  Here
	// we're outputting the result of the running the script to a buffer.  We need
	// to do this so we can remove the whitespace at the end of each line of the
	// table.  The test fails with the newline present as go tests implementation
	// of output: for examples, trims spaces.

	scanner := bufio.NewScanner(&b)
	for scanner.Scan() {
		fmt.Println(strings.TrimSpace(scanner.Text()))
	}
	// output:
	// FirstName   MiddleName  Surname
	// Marcus      Tullius     Cicero
	// Gaius       Julius      Caesar
	// Marcus      Licinius    Crassus
}

func ExampleOptTableXAlt() {
	data := []struct {
		FirstName string
		Mask      uint32
	}{
		{"Marcus", 255},
		{"Gaius", 10},
		{"Marcus", 6},
	}

	script := `{{tablexalt . 12 8 0}}`
	var b bytes.Buffer
	if err := OutputToTemplate(&b, "names", script, data, nil); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(&b)
	for scanner.Scan() {
		fmt.Println(strings.TrimSpace(scanner.Text()))
	}
	// output:
	// FirstName   Mask
	// "Marcus"    0xff
	// "Gaius"     0xa
	// "Marcus"    0x6
}

func ExampleOptHTableX() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Output the names of people in a series of nicely formatted tables
	script := `{{htablex . 12 8 0}}`
	var b bytes.Buffer
	if err := OutputToTemplate(&b, "names", script, data, nil); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(&b)
	for scanner.Scan() {
		fmt.Println(strings.TrimSpace(scanner.Text()))
	}
	// output:
	// FirstName:  Marcus
	// MiddleName: Tullius
	// Surname:    Cicero
	//
	// FirstName:  Gaius
	// MiddleName: Julius
	// Surname:    Caesar
	//
	// FirstName:  Marcus
	// MiddleName: Licinius
	// Surname:    Crassus
}

func ExampleOptHTableXAlt() {
	data := []struct {
		FirstName string
		Mask      uint32
	}{
		{"Marcus", 255},
		{"Gaius", 10},
		{"Marcus", 6},
	}

	script := `{{htablexalt . 12 8 0}}`
	var b bytes.Buffer
	if err := OutputToTemplate(&b, "names", script, data, nil); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(&b)
	for scanner.Scan() {
		fmt.Println(strings.TrimSpace(scanner.Text()))
	}

	// output:
	// FirstName:  "Marcus"
	// Mask:       0xff
	//
	// FirstName:  "Gaius"
	// Mask:       0xa
	//
	// FirstName:  "Marcus"
	// Mask:       0x6
}

func ExampleOptCols() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Output the first and last names of people in a nicely formatted table
	script := `{{tablex (cols . "FirstName" "Surname") 12 8 0}}`
	var b bytes.Buffer
	if err := OutputToTemplate(&b, "names", script, data, nil); err != nil {
		panic(err)
	}

	// Normally you would pass os.Stdout directly into OutputToTemplate.  Here
	// we're outputting the result of the running the script to a buffer.  We need
	// to do this so we can remove the whitespace at the end of each line of the
	// table.  The test fails with the newline present as go tests implementation
	// of output: for examples, trims spaces.

	scanner := bufio.NewScanner(&b)
	for scanner.Scan() {
		fmt.Println(strings.TrimSpace(scanner.Text()))
	}
	// output:
	// FirstName   Surname
	// Marcus      Cicero
	// Gaius       Caesar
	// Marcus      Crassus
}

func ExampleOptSort() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Output the names of people sorted by their Surnames
	script := `{{tablex (sort . "Surname") 12 8 0}}`
	var b bytes.Buffer
	if err := OutputToTemplate(&b, "names", script, data, nil); err != nil {
		panic(err)
	}

	// Normally you would pass os.Stdout directly into OutputToTemplate.  Here
	// we're outputting the result of the running the script to a buffer.  We need
	// to do this so we can remove the whitespace at the end of each line of the
	// table.  The test fails with the newline present as go tests implementation
	// of output: for examples, trims spaces.

	scanner := bufio.NewScanner(&b)
	for scanner.Scan() {
		fmt.Println(strings.TrimSpace(scanner.Text()))
	}
	// output:
	// FirstName   MiddleName  Surname
	// Gaius       Julius      Caesar
	// Marcus      Tullius     Cicero
	// Marcus      Licinius    Crassus
}

func ExampleOptRows() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Print the surname of the first and third people in the database
	script := `{{range (rows . 0 2)}}{{println .Surname}}{{end}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Cicero
	// Crassus
}

func ExampleOptHead() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Print the surname of the first person in the database
	script := `{{range (head .)}}{{println .Surname}}{{end}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Cicero
}

func ExampleOptTail() {
	data := []struct{ FirstName, MiddleName, Surname string }{
		{"Marcus", "Tullius", "Cicero"},
		{"Gaius", "Julius", "Caesar"},
		{"Marcus", "Licinius", "Crassus"},
	}

	// Print the surname of the first person in the database
	script := `{{range (tail .)}}{{println .Surname}}{{end}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Crassus
}

func ExampleOptDescribe() {
	data := []struct{ FirstName, MiddleName, Surname string }{}

	// Describe the type of data
	script := `{{describe .}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// []struct {
	// 	FirstName  string
	// 	MiddleName string
	//	Surname    string
	// }
}

func ExampleOptPromote() {
	type cred struct {
		Name     string
		Password string
	}

	type u struct {
		Credentials cred
	}

	data := []struct {
		Uninteresting int
		User          u
	}{
		{0, u{cred{"Marcus", "1234"}}},
		{0, u{cred{"Gaius", "0000"}}},
	}

	// Create a new []cred containing the credentials embedded within data,
	// iterate through this new slice printing out the names and passwords.
	// The cred instances rooted at "User.Credentials" in the data object
	// are promoted to the top level in the new slice.
	script := `{{range (promote . "User.Credentials")}}{{printf "%s %s\n" .Name .Password}}{{end}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Marcus 1234
	// Gaius 0000
}

func ExampleOptSliceof() {
	script := `{{index (sliceof .) 0}}`
	if err := OutputToTemplate(os.Stdout, "names", script, 1, nil); err != nil {
		panic(err)
	}
	// output:
	// 1
}

func ExampleOptToTable() {
	data := [][]string{
		{"Message", "Code", "Occurrence"},
		{"Too many GOSUBs", "37", "0.1"},
		{"Too many REPEATs", "44", "0.15"},
	}
	script := `{{with (totable .)}}{{select . "Message"}}{{select . "Code"}}{{select . "Occurrence"}}{{end}}`
	if err := OutputToTemplate(os.Stdout, "errors", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// Too many GOSUBs
	// Too many REPEATs
	// 37
	// 44
	// 0.1
	// 0.15
}

func ExampleOptSelectAlt() {
	data := []struct{ Integer uint32 }{{255}}
	script := `{{selectalt . "Integer"}}`
	if err := OutputToTemplate(os.Stdout, "names", script, data, nil); err != nil {
		panic(err)
	}
	// output:
	// 0xff
}

func ExampleConfig_AddCustomFn() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	cfg := NewConfig(OptAllFns)
	err := cfg.AddCustomFn(func(n []int) int {
		sum := 0
		for _, num := range n {
			sum += num
		}
		return sum
	}, "sum", "- sum \"Returns\" the sum of a slice of integers")
	if err != nil {
		panic(err)
	}

	// Print the sum of a slice of numbers
	script := `{{println (sum .)}}`
	if err = OutputToTemplate(os.Stdout, "sums", script, nums, cfg); err != nil {
		panic(err)
	}
	// output:
	// 55
}
