package main
// hlswatch - keep track of hls viewer stats
// Copyright (C) 2017 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.


// --------------------------------------------------------------------------------------
//  imports
// --------------------------------------------------------------------------------------

import(
    "log"
    "flag"
    "runtime"
    "os"
    "os/signal"
    "net/http"
    "syscall"
    "context"
    "time"

    "github.com/faryon93/hlswatch/handler"
    "github.com/faryon93/hlswatch/state"
)


// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

var (
    // configuration options
    listen          string
    hlsPath         string
    shutdownTimeout time.Duration
    viewerTimeout   time.Duration
    cycleTime       time.Duration

    // runtime variables
    Context *state.State = state.New()
)


// --------------------------------------------------------------------------------------
//  application entry
// --------------------------------------------------------------------------------------

func main() {
    log.Println("hlswatch version 0.1 #54dasf78")

    // setup go environment to use all available cpu cores
    runtime.GOMAXPROCS(runtime.NumCPU())

    // parse command line arguments
    flag.StringVar(&listen, "listen", ":3000", "")
    flag.StringVar(&hlsPath, "hlspath", "/tmp/hls", "")
    flag.DurationVar(&viewerTimeout, "viewertimeout", 10, "")
    flag.DurationVar(&shutdownTimeout, "shutdowntimeout", 5, "")
    flag.DurationVar(&cycleTime, "cycletime", 1, "")
    flag.Parse()

    // setup the http static file server serving the playlists
    // TODO: gzip compression for playlist, caching in ram, inotify, ...
    mux := http.NewServeMux()
    mux.Handle("/", handler.Hls(Context, http.FileServer(http.Dir(hlsPath))))
    srv := &http.Server{Addr: listen, Handler: mux}

    // serve the content via http
    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Println("failed start http server:", err.Error())
            os.Exit(-1) // TODO: clean shutdown
        }
    }()
    log.Println("http is listening on", listen)

    // fire the statistics computation task
    go StatisticsTask(Context)

    // wait for a signal to shutdown the application
    wait(os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    log.Println("gracefully shutting down application...")

    // gracefully shutdown the server
    ctx, _ := context.WithTimeout(context.Background(), shutdownTimeout * time.Second)
    srv.Shutdown(ctx)

    log.Println("application successfully exited")
}


// --------------------------------------------------------------------------------------
//  helper functions
// --------------------------------------------------------------------------------------

func wait(sig ...os.Signal) {
    signals := make(chan os.Signal)
    signal.Notify(signals, sig...)
    <- signals
}
