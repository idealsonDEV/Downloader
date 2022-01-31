package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func TaskPrint() {
	down, _ := ioutil.ReadDir("./myget/")
	folders := make([]string, 0)
	for _, d := range down {
		if d.IsDir() {
			folders = append(folders, d.Name())
		}
	}

	folderString := strings.Join(folders, "\n")
	fmt.Println("Currently on going download: ")
	fmt.Println(folderString)

}

/*func Read(task string) (*State, error) {
 	file := filepath.Join(os.Getenv("HOME"), dataFolder, task, stateFileName)
	Printf("Getting data from %s\n", file)
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	s := new(State)
	err = json.Unmarshal(bytes, s)
	return s, err
}*/
