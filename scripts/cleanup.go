package main

import (
"fmt"
"github.com/GDSC-Phenikaa/ctf-backend/db"
"github.com/GDSC-Phenikaa/ctf-backend/models"
"github.com/joho/godotenv"
)

func main() {
godotenv.Load(".env")
database, _ := db.Connect()
res := database.Where("content LIKE ?", "Bulk Auto-Generated Question %").Unscoped().Delete(&models.Question{})
fmt.Printf("Deleted %d fake questions\n", res.RowsAffected)
}
