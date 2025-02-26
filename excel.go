package excel

import(
	"excel-to-db/models"
	//"github.com/xuri/excelize/v2"
    "io"
    "os"
    "encoding/csv"

)

func ReadCSV(filePath string) ([]models.Users, error) {
    // Open the CSV file
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Create a new CSV reader
    csvReader := csv.NewReader(file)

    var users []models.Users

    // Read each record from CSV
    for {
        record, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

        // Assuming User struct and parsing logic
        user := models.Users{
            Name:  record[0],
            Email: record[1],
            // Add more fields as needed
        }

        users = append(users, user)
    }

    return users, nil
}


