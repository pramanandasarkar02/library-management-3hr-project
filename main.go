package main

import (
	"pramanandasarkar02/library-management/internal"
)

func main() {
	dataReadSeassion := internal.InitDb()
	internal.StartServer(dataReadSeassion)
}
