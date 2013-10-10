package main

import "fmt"
import "bitbucket.org/shatteredscreens/federalservices/character"
import "bitbucket.org/shatteredscreens/federalservices/federation"
import "bitbucket.org/shatteredscreens/federalservices/loot"
import "bitbucket.org/shatteredscreens/federalservices/profile"
import "bitbucket.org/shatteredscreens/federalservices/realm"
import "bitbucket.org/shatteredscreens/federalservices/zone"

func main() {
	fmt.Println("main")
	
	character.Hello ()
	federation.Hello()
	loot.Hello()
	profile.Hello()
	realm.Hello()
	zone.Hello()
}

