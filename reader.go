package cartogopher

import (
	"encoding/csv"
	"os"
)

// MapReader contains all our necessary data for the various methods to function
type MapReader struct {
	Headers        []string
	HeaderIndexMap map[string]int
	Reader         *csv.Reader
}

// CreateHeaderIndexMap creates a map of header strings to their indices in the array generated by
// encoding/csv's reader. For instance, if your CSV file looks something like
// this:
// 	---------------------
// 	| one | two | three |
// 	---------------------
// 	|  A  |  B  |   C   |
// 	---------------------
// Go's generated array for the header row will be [ "one", "two", "three" ].
// Cartogopher's generated map for the header row will be { "one": 1, "two": 2, "three": 3 }
func (m *MapReader) CreateHeaderIndexMap(headers []string) map[string]int {
	headerIndexMap := make(map[string]int, len(headers))

	for index, header := range headers {
		headerIndexMap[header] = index
	}

	return headerIndexMap
}

// CreateRowMap takes a given CSV array and returns a map of column names to the values contained therein.
// For instance, if your CSV file looks something like this:
// 	---------------------
// 	| one | two | three |
// 	---------------------
// 	|  A  |  B  |   C   |
// 	---------------------
// The return result will be:
// 	{
// 		"one": "A",
// 		"two": "B",
// 		"three": "C",
// 	}
func (m *MapReader) CreateRowMap(csvRow []string) map[string]string {
	result := map[string]string{}
	for header, index := range m.HeaderIndexMap {
		result[header] = csvRow[index]
	}

	return result
}

// Read mimics the built-in CSV reader Read method
func (m *MapReader) Read() (map[string]string, error) {
	csvRow, err := m.Reader.Read()
	if err != nil {
		return nil, err
	}
	return m.CreateRowMap(csvRow), nil
}

// ReadAll mimics the built-in CSV reader ReadAll method
func (m *MapReader) ReadAll() ([]map[string]string, error) {
	records, err := m.Reader.ReadAll()
	if err != nil {
		return nil, err
	}
	results := []map[string]string{}
	for _, record := range records {
		results = append(results, m.CreateRowMap(record))
	}
	return results, nil
}

// NewReader returns a new MapReader struct
func NewReader(file *os.File) (*MapReader, error) {
	// Create our reader
	reader := csv.NewReader(file)

	// Create our resulting struct
	output := &MapReader{}

	inputHeaders, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Use our methods (defined above) to populate our struct fields
	output.Headers = inputHeaders
	output.Reader = reader
	output.HeaderIndexMap = output.CreateHeaderIndexMap(inputHeaders)

	return output, nil
}