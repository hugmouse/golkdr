package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/hugmouse/golkdr"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Your phone number: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	text = strings.ReplaceAll(text, "\n", "")

	number, err := strconv.Atoi(text)
	if err != nil {
		log.Fatal(err)
	}

	u := golkdr.NewUser(uint(number))
	err = u.RequestSMS()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("SMS code from FNS_Russia: ")
	text, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	text = strings.ReplaceAll(text, "\n", "")

	textToInt, err := strconv.Atoi(text)
	if err != nil {
		log.Fatal(err)
	}

	err = u.SetCodeFromSMS(textToInt)
	if err != nil {
		log.Fatal(err)
	}

	info, err := json.MarshalIndent(u.AuthInfo, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("./key.json", info, 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Login successful! Authorization info saved in ./key.json")
}
