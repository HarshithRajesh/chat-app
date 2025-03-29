// package database
//
// import(
//   "fmt"
//   "log"
//   "os"
//   "gorm.io/driver/postgres"
//   "gorm.io/gorm"
// )
//
// var DB *gorm.DB
//
// func connectDB(){
//   dsn := os.Getenv("DATABASE_URL")
//   DB,err = gorm.Open(postgres.Open(dsn),&gorm.Config{})
//   if err != nil{
//     log.Fatal("Failed to connect to dastabase :",err)
//   }
//   fmt.Println("Suc")
// }
