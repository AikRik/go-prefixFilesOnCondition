package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", upload)
	http.HandleFunc("/upload", uploadFiles)
	http.ListenAndServe("localhost:8080", nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `
		This will go through the file and will per line look for the value of each string separated by a space(" ").<br>
		Specify the value you wish to search for and the value you want to replace it with.<br>
		Optionally, you can specify from which space " " on you wish to search.<br><br>

		<form action="/upload" method="POST" enctype="multipart/form-data">
			<input type="file" name="files" multiple required="true"><br>
			Search Value: <input type="text" name="searchValue" required="true"><br>
			Replace Value: <input type="text" name="replaceValue" required="true"><br>
			Space Counter: <input type="number" name="column"><br>
			<input type="submit">
		</form>
	`)
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	files := r.MultipartForm.File["files"]
	searchValue := r.FormValue("searchValue")
	replaceValue := r.FormValue("replaceValue")
	index, _ := strconv.Atoi(r.FormValue("number"))

	for _, file := range files {
		condB := false
		f, _ := file.Open()
		scanner := bufio.NewScanner(f)
		fileContent := ""

		i := 0
		for scanner.Scan() {
			stringArr := strings.Split(scanner.Text(), " ")
			text := stringArr[0]

			j := 0
			for j < len(stringArr) {
				// if index
				if j >= index {
					if stringArr[j] == searchValue {
						stringArr[j] = replaceValue
					} else {
						condB = true
					}
				}

				text += " " + stringArr[j]
				j++
			}

			fileContent += text + "\n"
			i++
		}

		filename := file.Filename
		if condB == true {
			filename = "B_" + filename
		}

		ioutil.WriteFile(filename, []byte(fileContent), 0644)
		f.Close()

		fmt.Fprintf(w, `Upload completed`)
	}
}
