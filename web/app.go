package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type memStruct struct {
	Total_mem int
	Free_mem  int
}

type cpuinfo struct {
	Pid   string
	Name  string
	User  string
	State string
	Ram   string
}

func main() {
	http.HandleFunc("/memoria", ramInfo)
	http.HandleFunc("/", cpu)
	http.HandleFunc("/", rendcpu)
	http.ListenAndServe(":3000", nil)
}

func ramInfo(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}

	str := string(b)
	listaInfo := strings.Split(string(str), "\n")

	memoriaTotal := strings.Replace((listaInfo[0])[10:24], " ", "", -1)
	memoriaLibre := strings.Replace((listaInfo[1])[10:24], " ", "", -1)

	ramTotalKB, err1 := strconv.Atoi(memoriaTotal)
	ramFreeKB, err2 := strconv.Atoi(memoriaLibre)

	if err1 == nil && err2 == nil {
		ramTotalMB := ramTotalKB / 1024
		ramFreeMB := ramFreeKB / 1024

		memResponse := memStruct{ramTotalMB, ramFreeMB}
		jsonResponse, errorjson := json.Marshal(memResponse)
		if errorjson != nil {
			http.Error(w, errorjson.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		return
	}
}

func cpu(w http.ResponseWriter, r *http.Request) {
	var cuerpo []cpuinfo
	archivos, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Fatal(err)
	}

	contZ := 0
	contS := 0
	contR := 0
	contI := 0
	contT := 0
	total := 0

	for _, archivo := range archivos {
		if archivo.IsDir() {
			b, err := ioutil.ReadFile("/proc/" + archivo.Name() + "/status")
			if err != nil {
				continue
			}

			str := string(b)
			listaInfo := strings.Split(string(str), "\n")

			nombre := strings.Replace((listaInfo[0])[6:], " ", "", -1)
			estado := strings.Replace((listaInfo[2])[7:], " ", "", -1)
			pid := strings.Replace((listaInfo[5])[5:], " ", "", -1)
			memoria := strings.Replace((listaInfo[17])[8:], " ", "", -1)
			id := ""
			if len(strings.Replace((listaInfo[8])[5:], " ", "", -1)) == 7 {
				id = (strings.Replace((listaInfo[8])[5:6], " ", "", -1))
			} else if len(strings.Replace((listaInfo[8])[5:], " ", "", -1)) == 19 {
				id = (strings.Replace((listaInfo[8])[5:9], " ", "", -1))
			} else if len(strings.Replace((listaInfo[8])[5:], " ", "", -1)) == 15 {
				id = (strings.Replace((listaInfo[8])[5:8], " ", "", -1))
			} else if len(strings.Replace((listaInfo[8])[5:], " ", "", -1)) == 23 {
				id = (strings.Replace((listaInfo[8])[5:10], " ", "", -1))
			}

			user, err := ioutil.ReadFile("/etc/passwd")
			if err != nil {
				return
			}
			str1 := string(user)
			users := strings.Split(string(str1), "\n")
			for i := 0; i < len(users)-1; i++ {
				userids := strings.Split(string(users[i]), ":")
				nid, _ := strconv.Atoi(userids[2])
				aux, _ := strconv.Atoi(id)

				if nid == aux {
					id = userids[0]
				}
			}

			est := estado[0:1]
			if est == "S" {
				contS++
			} else if est == "R" {
				contR++
			} else if est == "Z" {
				contZ++
			} else if est == "T" {
				contT++
			} else if est == "I" {
				contI++
			}

			total++

			cpuStruct := cpuinfo{Pid: pid, Name: nombre, User: id, State: estado, Ram: memoria}

			cuerpo = append(cuerpo, cpuStruct)
		}
	}

	totalres := cpuinfo{strconv.Itoa(contR), strconv.Itoa(contS), strconv.Itoa(contT), strconv.Itoa(contZ), strconv.Itoa(total)}
	cuerpo = append(cuerpo, totalres)

	jsonResponse, errorjson := json.Marshal(cuerpo)
	if errorjson != nil {
		http.Error(w, errorjson.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func rendcpu(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}

	str := string(b)
	listaInfo := strings.Split(string(str), "\n")

	memoriaTotal := strings.Replace((listaInfo[0])[10:24], " ", "", -1)
	memoriaLibre := strings.Replace((listaInfo[1])[10:24], " ", "", -1)

	ramTotalKB, err1 := strconv.Atoi(memoriaTotal)
	ramFreeKB, err2 := strconv.Atoi(memoriaLibre)

	if err1 == nil && err2 == nil {
		ramTotalMB := ramTotalKB / 1024
		ramFreeMB := ramFreeKB / 1024

		memResponse := memStruct{ramTotalMB, ramFreeMB}
		jsonResponse, errorjson := json.Marshal(memResponse)
		if errorjson != nil {
			http.Error(w, errorjson.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		return
	}
}
