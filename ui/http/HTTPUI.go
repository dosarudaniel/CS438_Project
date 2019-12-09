package http

import (
	"encoding/json"

	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
	"github.com/2_alt_hw2/Peerster/ui"
	"github.com/2_alt_hw2/Peerster/types"
	"github.com/2_alt_hw2/Peerster/file_sharing"
	"github.com/2_alt_hw2/Peerster/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func unmarshalBody(body io.Reader, structure interface{}) error {
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(reqBody, structure)
	if err != nil {
		return err
	}
	return nil
}

func readChanAndSend(w http.ResponseWriter, object interface{}) {
	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func returnAllPeers(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := make(chan []string)
		uiRequester.RequestPullPeers(func(peers []string) { channel <- peers })
		readChanAndSend(w, <-channel)
	}
}

func pushPeer(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		wg.Add(1)
		uiRequester.RequestPushPeer(string(reqBody), func() { wg.Done() })
		wg.Wait()
	}
}

func returnAllContacts(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := make(chan []string)
		uiRequester.RequestPullContacts(func(peers []string) { channel <- peers })
		readChanAndSend(w, <-channel)
	}
}

func pushContact(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		wg.Add(1)
		uiRequester.RequestPushContact(string(reqBody), func() { wg.Done() })
		wg.Wait()
	}
}

func returnAllRumors(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		offset, err := strconv.Atoi(mux.Vars(r)["offset"])
		if err == nil {
			channel := make(chan []*types.RumorMessage)
			uiRequester.RequestPullRumors(uint32(offset), func(rumors []*types.RumorMessage) { channel <- rumors })
			readChanAndSend(w, <-channel)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func pushRumor(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rumor types.Message
		err := unmarshalBody(r.Body, &rumor)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		wg.Add(1)
		uiRequester.RequestPushRumor(rumor, func() { wg.Done() })
		wg.Wait()
	}
}

func pullSettings(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := make(chan types.Settings)
		uiRequester.RequestPullSettings(func(settings types.Settings) { channel <- settings })
		readChanAndSend(w, <-channel)
	}
}

func pushSettings(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var settings types.Settings
		err := unmarshalBody(r.Body, &settings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		wg.Add(1)
		uiRequester.RequestPushSettings(settings, func() { wg.Done() })
		wg.Wait()
	}
}

func returnAllPrivateMessages(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		offset, err := strconv.Atoi(mux.Vars(r)["offset"])
		if err == nil {
			channel := make(chan []*types.PrivateMessage)
			uiRequester.RequestPullPrivateMessage(uint32(offset), func(rumors []*types.PrivateMessage) { channel <- rumors })
			readChanAndSend(w, <-channel)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func pushPrivateMessage(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rumor types.Message
		err := unmarshalBody(r.Body, &rumor)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		wg.Add(1)
		uiRequester.RequestPushPrivateMessage(rumor, func() { wg.Done() })
		wg.Wait()
	}
}

func returnAllSharedFiles(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := make(chan []file_sharing.File)
		uiRequester.RequestPullSharedFiles(func(files []file_sharing.File) { channel <- files })
		readChanAndSend(w, <-channel)
	}
}

func returnAllSharableFiles(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := make(chan []string)
		uiRequester.RequestPullSharableFiles(func(files []string) { channel <- files })
		readChanAndSend(w, <-channel)
	}
}

func pushSharedFile(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		wg.Add(1)
		uiRequester.RequestPushSharedFile(mux.Vars(r)["filename"], func() { wg.Done() })
		wg.Wait()
	}
}

func pushDownloadFile(uiRequester *ui.Requester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		wg.Add(1)
		file := file_sharing.File{
			Name:   mux.Vars(r)["filename"],
			Origin: mux.Vars(r)["origin"],
			Hash:   mux.Vars(r)["hash"],
		}
		uiRequester.RequestPushDownloadFile(file, func() { wg.Done() })
		wg.Wait()
	}
}

// NewHttpUI creates a new uiRequester based on an HTTP API
func NewHttpUI(address string, log *logger.Logger) (ui.BackendInterface, error) {
	uiRequester := ui.NewRequester()

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	router.HandleFunc("/peers", returnAllPeers(uiRequester)).Methods("GET")
	router.HandleFunc("/peer", pushPeer(uiRequester)).Methods("POST")
	router.HandleFunc("/contacts", returnAllContacts(uiRequester)).Methods("GET")
	router.HandleFunc("/contact", pushContact(uiRequester)).Methods("POST")
	router.HandleFunc("/rumors/{offset}", returnAllRumors(uiRequester)).Methods("GET")
	router.HandleFunc("/rumor", pushRumor(uiRequester)).Methods("POST")
	router.HandleFunc("/settings", pullSettings(uiRequester)).Methods("GET")
	router.HandleFunc("/settings", pushSettings(uiRequester)).Methods("POST")
	router.HandleFunc("/privateMessages/{offset}", returnAllPrivateMessages(uiRequester)).Methods("GET")
	router.HandleFunc("/privateMessage", pushPrivateMessage(uiRequester)).Methods("POST")
	router.HandleFunc("/sharedFiles", returnAllSharedFiles(uiRequester)).Methods("GET")
	router.HandleFunc("/sharableFiles", returnAllSharableFiles(uiRequester)).Methods("GET")
	router.HandleFunc("/sharedFile/{filename}", pushSharedFile(uiRequester)).Methods("POST")
	router.HandleFunc("/downloadFile/{origin}/{filename}/{hash}", pushDownloadFile(uiRequester)).Methods("POST")

	// CORS handling
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	var err error = nil
	go func() {
		err = http.ListenAndServe(address, handlers.CORS(headers, methods, origins)(router))
	}()
	log.Debug("Waiting for HTTP server to come up")
	time.Sleep(time.Millisecond * 200) // Let the server return it's error
	log.Debug("Waiting for HTTP: done")

	return uiRequester.BackendInterface(), err
}
