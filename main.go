package main

import(
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "net/http"
	"time"
	"excel-to-db/excel"
	"excel-to-db/models"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "fmt"
    "log" 
)

 var userCollection *mongo.Collection

func main() {
    // Initialize Gin router
    r := gin.Default()
    
    // Connect to MongoDB
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Assign the MongoDB collection
    userCollection = client.Database("your_database").Collection("your_collection")

    // Define route handlers
    r.POST("/upload", uploadFile)
    r.GET("/users", getUser)
    r.GET("/users/:id", getUserID)
    r.PUT("/users/:id", updateUser)
    r.DELETE("/users/:id", deleteUser)

    // Start the server
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }

}


func uploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }

    filePath := "./" + file.Filename
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    users, err := excel.ReadCSV(filePath) // Assuming ReadCSV reads CSV format
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read CSV"})
        return
    }

    var docs []interface{}
    for _, user := range users {
        docs = append(docs, user)
    }

    // Assuming userCollection is your MongoDB collection
    _, err = userCollection.InsertMany(context.Background(), docs)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into database"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "File processed successfully"})
}


func getUserID(c *gin.Context) {
    id := c.Param("id")
    fmt.Println("Fetching user with ID:", id)
    
    // Convert id to ObjectID
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    
    var user models.Users // Adjust the model type accordingly
    filter := bson.M{"_id": objID}
    
    err = userCollection.FindOne(context.Background(), filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

func getUser(c *gin.Context) {
    filter := bson.M{}
    var users []models.Users
    cur, err := userCollection.Find(context.Background(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
        return
    }
    defer cur.Close(context.Background())

    for cur.Next(context.Background()) {
        var user models.Users
        if err := cur.Decode(&user); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user"})
            return
        }
        users = append(users, user)
    }
    if err := cur.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
        return
    }

    c.JSON(http.StatusOK, users)
}

func updateUser(c *gin.Context) {
    id := c.Param("id")

    // Convert id to ObjectID
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    fmt.Println("Converted ObjectID:", objID)

    var user models.Users
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    fmt.Println("Received user updates:",user)

    // Set the filter with the ObjectID
    filter := bson.M{"_id": objID}

    // Prepare the update
    update := bson.M{"$set": user}

    // Perform the update operation
    _, err = userCollection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}


func deleteUser(c *gin.Context) {
    id := c.Param("id")

    // Convert id to ObjectID
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    // Perform the deletion operation
    _, err = userCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

