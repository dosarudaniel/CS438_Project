package web

import (
	"context"
	"github.com/dosarudaniel/CS438_Project/chord"
	"github.com/dosarudaniel/CS438_Project/services/client_service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var server = struct {
	chordNode *chord.ChordNode
	guiIPAddr string
}{}

// RunServer runs the web server for gossiper
func RunServer(guiIPAddr string, chordNode *chord.ChordNode) {
	var err error

	server.chordNode = chordNode
	server.guiIPAddr = guiIPAddr

	r := gin.Default()
	r.LoadHTMLGlob("./web/static/*")
	r.Static("assets/", string(http.Dir("./web/static"))) // TODO: find a nicer solution

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r.GET("/", indexHandler)
	r.POST("/upload_file", postUploadFile)
	r.POST("/search_file", postSearchFileHandler)
	r.GET("/download_file", getDownloadFile)

	err = r.Run(guiIPAddr)
	if err != nil {
		log.Printf("Web server could not start at %s : %v", guiIPAddr, err)
	}
}

// GET "/"
// @return index.tmpl with filled in peer name
func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", struct {
		Name string
		IP   string
	}{server.chordNode.ID(), server.chordNode.IP()})
}

func postSearchFileHandler(c *gin.Context) {
	var err error

	reqBody := struct {
		Query string `json:"query"`
	}{}

	err = c.BindJSON(&reqBody)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fileRecords, err := server.chordNode.SearchFile(context.Background(), &client_service.Query{
		Query: reqBody.Query,
	})

	c.JSON(http.StatusOK, fileRecords.FileRecords)
}

func postUploadFile(c *gin.Context) {
	// single file
	file, err := c.FormFile("file_to_share")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	absolutePathToFile, err := JoinToAbsolutePath("_upload", file.Filename)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = c.SaveUploadedFile(file, absolutePathToFile)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = server.chordNode.UploadFile(context.Background(), &client_service.Filename{
		Filename: file.Filename,
	})
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func getDownloadFile(c *gin.Context) {
	var err error

	filenameToDownload := c.Query("filename")
	ownerIP := c.Query("owner_ip")

	if filenameToDownload == "" || ownerIP == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = server.chordNode.RequestFileFromIP(filenameToDownload, filenameToDownload, ownerIP)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	time.Sleep(time.Second)

	absolutePathToFile, err := JoinToAbsolutePath("_download", filenameToDownload)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filenameToDownload)
	c.File(absolutePathToFile)
}

// JoinToAbsolutePath returns an absolute path to the project appended by the given strings
// in sequential order
func JoinToAbsolutePath(pathRelativeToExecFolder ...string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	absolutePathSlice := append([]string{workingDir}, pathRelativeToExecFolder...)

	return filepath.Join(absolutePathSlice...), nil
}
