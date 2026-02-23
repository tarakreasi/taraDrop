package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Configuration constants
const (
	UploadFolder       = "uploads"
	MaxUploadSize      = 16 * 1024 * 1024 * 1024 // 16GB
	ServerReadTimeout  = 1 * time.Hour
	ServerWriteTimeout = 1 * time.Hour
)

//go:embed templates/index.html
var content embed.FS

var (
	serverApp   fyne.App
	mainWindow  fyne.Window
	statusLabel *widget.Label
)

func main() {
	// Initialize environment
	setupUploadDirectory()

	// 1. Setup GUI App
	serverApp = app.New()
	mainWindow = serverApp.NewWindow("taraDrop Control Panel")
	mainWindow.Resize(fyne.NewSize(400, 300))
	mainWindow.SetFixedSize(true)

	// 2. Setup Server Handler
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/upload_chunk", handleUploadChunk)

	// 3. Start HTTP Server in Background
	go startServer()

	// 4. Build GUI Layout
	buildUI()

	// 5. Run GUI
	mainWindow.ShowAndRun()
}

func startServer() {
	var listener net.Listener
	var err error
	var port string

	// Try port 80 first on Windows, fallback to 5000
	if runtime.GOOS == "windows" {
		port = ":80"
		listener, err = net.Listen("tcp", port)
		if err != nil {
			log.Printf("Failed to bind to port 80: %v. Falling back to port 5000.", err)
			port = ":5000"
			listener, err = net.Listen("tcp", port)
		}
	} else {
		// Non-Windows always use 5000
		port = ":5000"
		listener, err = net.Listen("tcp", port)
	}

	if err != nil {
		log.Fatalf("Server failed to bind to any port: %v", err)
	}

	server := &http.Server{
		ReadTimeout:    ServerReadTimeout,
		WriteTimeout:   ServerWriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Update GUI with URL once server is technically "ready"
	ip := getLocalIP()
	urlStr := fmt.Sprintf("http://%s%s", ip, port)
	// If port 80, cleaner URL
	if port == ":80" {
		urlStr = fmt.Sprintf("http://%s", ip)
	}

	log.Printf("Server starting at %s", urlStr)

	// Refresh UI label on main thread
	if statusLabel != nil {
		// Ensure UI updates happen on main thread to avoid races
		go func(u string) {
			// Small delay to ensure UI is drawn before we update it
			time.Sleep(100 * time.Millisecond)
			statusLabel.SetText(fmt.Sprintf("Running at:\n%s", u))
		}(urlStr)
	}

	if err := server.Serve(listener); err != nil {
		log.Printf("Server failed: %v", err)
	}
}

func buildUI() {
	// Header
	title := widget.NewLabelWithStyle("taraDrop Server", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// IP/Status (Initial state, will be updated by startServer)
	statusLabel = widget.NewLabel("Starting server...")
	statusLabel.Alignment = fyne.TextAlignCenter

	// Upload Folder Button
	uploadPath, _ := filepath.Abs(UploadFolder)

	openFolderBtn := widget.NewButton("📂 Open Upload Folder", func() {
		openDir(uploadPath)
	})

	// Link Instruction
	instruction := widget.NewLabel("Ensure your phone is on the same WiFi.")
	instruction.Alignment = fyne.TextAlignCenter
	instruction.TextStyle = fyne.TextStyle{Italic: true}

	// Stop Button
	stopBtn := widget.NewButtonWithIcon("STOP SERVER", nil, func() {
		serverApp.Quit()
	})
	// Simple red-ish feedback logic isn't built-in easily without custom theme,
	// but standard button is fine.

	// Footer Link
	linkUrl, _ := url.Parse("https://tarakreasi.com")
	footer := widget.NewHyperlink("create by Tri Wantoro | tarakreasi.com", linkUrl)
	footer.Alignment = fyne.TextAlignCenter

	// Layout
	content := container.New(layout.NewVBoxLayout(),
		layout.NewSpacer(),
		title,
		layout.NewSpacer(),
		statusLabel,
		layout.NewSpacer(),
		openFolderBtn,
		instruction,
		layout.NewSpacer(),
		widget.NewSeparator(),
		stopBtn,
		layout.NewSpacer(),
		footer,
	)

	mainWindow.SetContent(content)
}

// openDir opens the directory using the OS default file manager
func openDir(path string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		log.Printf("Unsupported OS for opening folder: %s", runtime.GOOS)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to open folder: %v", err)
	}
}

// setupUploadDirectory ensures the upload directory exists
func setupUploadDirectory() {
	if _, err := os.Stat(UploadFolder); os.IsNotExist(err) {
		if err := os.Mkdir(UploadFolder, 0755); err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
	}
}

// handleIndex serves the embedded HTML template
func handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := content.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

// handleUploadChunk processes file chunks uploaded by the client
func handleUploadChunk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	// Parse multipart form (max memory 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get filename correctly (decode/sanitize)
	filename := r.FormValue("filename")

	// Basic sanitation
	filename = filepath.Base(filename)

	chunkIndex := r.FormValue("chunkIndex")

	// Validate filename
	if filename == "." || filename == "/" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	savePath := filepath.Join(UploadFolder, filename)

	// Determine file opening flags
	flags := os.O_WRONLY | os.O_CREATE
	if chunkIndex == "0" {
		flags |= os.O_TRUNC // Overwrite if it's the first chunk
	} else {
		flags |= os.O_APPEND // Append for subsequent chunks
	}

	// Open file for writing
	out, err := os.OpenFile(savePath, flags, 0644)
	if err != nil {
		log.Printf("Error opening file %s: %v", savePath, err)
		http.Error(w, "Unable to create/open file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Copy chunk to file
	if _, err = io.Copy(out, file); err != nil {
		log.Printf("Error writing to file %s: %v", savePath, err)
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	// Just minimal logging to avoid cluttering the GUI terminal if run from there
	// Ideally we could log to a GUI console widget, but not requested.
	if idx, err := strconv.Atoi(chunkIndex); err == nil && idx%10 == 0 {
		// Log every 10th chunk to stdout for debug
		fmt.Printf("Received chunk %d for %s\n", idx, filename)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Chunk received"))
}

// getLocalIP attempts to find the local machine's IP address
func getLocalIP() string {
	// Method 1: Connect to a public DNS server
	// This doesn't actually establish a connection, but asks the OS
	// which interface it would use to reach that IP.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}

	// Method 2: Fallback to iterating interfaces
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}
