package main
func main(){
	app := APP{}
	app.Initialise(DBUser, DBPassword, DBName)
	app.Run("localhost:3333")
}