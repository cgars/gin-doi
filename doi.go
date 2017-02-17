package main

import (
	"flag"
	"net/http"
	"log"
	"./src"
)

func main() {
	var (
		maxWorkers   = flag.Int("max_workers", 5, "The number of workers to start")
		maxQueueSize = flag.Int("max_queue_size", 100, "The size of job queue")
		port         = flag.String("port", "8081", "The server port")
		source       = flag.String("source", "https://repo.gin.g-node.org", "The default URI")
		baseTarget   = flag.String("target", "./", "The default base path for storgae")
	)
	flag.Parse()
	ds := ginDoi.GinDataSource{GinURL: *source}
	storage := ginDoi.LocalStorage{Path:*baseTarget, Source:ds}

	// Create the job queue.
	jobQueue := make(chan ginDoi.Job, *maxQueueSize)
	// Start the dispatcher.
	dispatcher := ginDoi.NewDispatcher(jobQueue, *maxWorkers)
	dispatcher.Run(ginDoi.NewWorker)

	// Start the HTTP handler.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ginDoi.InitDoiJob(w, r, &ds)
	})
	http.HandleFunc("/do/", func(w http.ResponseWriter, r *http.Request) {
		ginDoi.DoDoiJob(w,r,jobQueue, storage)
	})
	http.Handle("/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("/assets"))))

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

