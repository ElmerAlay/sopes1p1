package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type memStruct struct {
	Total_mem int
	Free_mem  int
	Porcent   int
}

type cpuinfo struct {
	Pid   string
	Name  string
	User  string
	State string
	Ram   string
}

type rend struct {
	CPU string
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/memoria", ramInfo).Methods("GET")
	router.HandleFunc("/", cpu).Methods("GET")
	router.HandleFunc("/cpu", rendcpu).Methods("GET")
	router.HandleFunc("/kill", deleteProcess).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", router))
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
		porcentaje := ((ramTotalMB - ramFreeMB) * 100) / ramTotalMB

		memResponse := memStruct{ramTotalMB, ramFreeMB, porcentaje}
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
	cmd := exec.Command("bash", "-c", "top -bn1 | awk '/Cpu/ { cpu = 100 - $8 }; END { print cpu }'")
	b, e := cmd.Output()
	if e != nil {
		log.Printf("failed due to :%vn", e)
		panic(e)
	}

	rendimiento := rend{CPU: string(b)[0 : len(string(b))-1]}
	jsonResponse, errorjson := json.Marshal(rendimiento)
	if errorjson != nil {
		http.Error(w, errorjson.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

type respuesta struct {
	Pid string
}

func deleteProcess(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var u respuesta
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return
	}

	cmd := exec.Command("bash", "-c", "kill -9 "+u.Pid)
	b, e := cmd.Output()
	if e != nil {
		log.Printf("failed due to :%vn", e)
		panic(e)
	}

	fmt.Println(b)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
