/*
 * Ghostdreamer LTS — Multi-Layout / Multi-Theme Go Host Engine
 * Version: v0.1.0-LTS
 * Features:
 * - Single binary deployment via go:embed (index.html).
 * - Selectable layouts (default, vlc) and themes (dark, vlc-orange, youtube).
 * - Live state synchronization, thread-safe sync.RWMutex, history logging.
 * - GitHub release auto-update check endpoint (/api/version).
 * - OS-tolerant network binding (127.0.0.1 with dynamic fallback).
 * Usage: go run main.go (or compile with go build -ldflags="-s -w")
 */

package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

//go:embed index.html
var contentFS embed.FS

const AppVersion = "v0.1.0-LTS"
const GitHubRepo = "NeoMarkiss/ghostdreamer"

type AppState struct {
	sync.RWMutex
	AudioSource       string            `json:"audio_source"`
	PrivateFeedToggle bool              `json:"private_feed_toggle"`
	VisualizerPreset  string            `json:"visualizer_preset"`
	FrequencyTemplate string            `json:"frequency_template"`
	StreamStatus      string            `json:"stream_status"`
	TargetFPS         int               `json:"target_fps"`
	Layout            string            `json:"layout"`
	Theme             string            `json:"theme"`
	ShowVisualizer    bool              `json:"show_visualizer"`
	ShowPlaylist      bool              `json:"show_playlist"`
	ShowControls      bool              `json:"show_controls"`
	CustomThemes      map[string]Theme  `json:"custom_themes"`
	History           []History         `json:"history"`
}

type Theme struct {
	Name    string `json:"name"`
	Bg      string `json:"bg"`
	Panel   string `json:"panel"`
	Border  string `json:"border"`
	Text    string `json:"text"`
	Muted   string `json:"muted"`
	Accent  string `json:"accent"`
	Accent2 string `json:"accent2"`
}

type History struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Profile   string    `json:"profile"`
	Metric    string    `json:"metric"`
	Passed    bool      `json:"passed"`
}

var defaultThemes = map[string]Theme{
	"dark":       {Name: "Dark", Bg: "#0f172a", Panel: "#1e293b", Border: "#334155", Text: "#f8fafc", Muted: "#94a3b8", Accent: "#3b82f6", Accent2: "#60a5fa"},
	"vlc-orange": {Name: "VLC Orange", Bg: "#1e1e1e", Panel: "#2d2d2d", Border: "#3c3c3c", Text: "#e0e0e0", Muted: "#a0a0a0", Accent: "#ff8800", Accent2: "#ffaa33"},
	"youtube":    {Name: "YouTube", Bg: "#0f0f0f", Panel: "#212121", Border: "#3f3f3f", Text: "#ffffff", Muted: "#aaaaaa", Accent: "#ff0000", Accent2: "#cc0000"},
}

var globalState = &AppState{
	AudioSource:       "",
	PrivateFeedToggle: true,
	VisualizerPreset:  "Frequency Bars",
	FrequencyTemplate: "Standard 256",
	StreamStatus:      "IDLE / READY",
	TargetFPS:         30,
	Layout:            "default",
	Theme:             "dark",
	ShowVisualizer:    true,
	ShowPlaylist:      true,
	ShowControls:      true,
	CustomThemes:      make(map[string]Theme),
	History: []History{{
		ID:        "hist-001",
		Timestamp: time.Now(),
		Profile:   "Ghostdreamer LTS Handshake",
		Metric:    "Web Audio & Multi-Theme Ready",
		Passed:    true,
	}},
}

type Manifest struct {
	Name            string `json:"name"`
	ShortName       string `json:"short_name"`
	StartURL        string `json:"start_url"`
	Display         string `json:"display"`
	BackgroundColor string `json:"background_color"`
	ThemeColor      string `json:"theme_color"`
}

type VersionInfo struct {
	Current         string `json:"current"`
	Latest          string `json:"latest"`
	UpdateAvailable bool   `json:"update_available"`
	ReleaseURL      string `json:"release_url"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/manifest.json", handleManifest)
	mux.HandleFunc("/api/state", handleGetState)
	mux.HandleFunc("/api/update", handleUpdateState)
	mux.HandleFunc("/api/history/add", handleAddHistory)
	mux.HandleFunc("/api/theme/add", handleAddTheme)
	mux.HandleFunc("/api/version", handleVersionCheck)

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatalf("Failed to bind port: %v", err)
		}
	}

	addr := listener.Addr().String()
	log.Printf("=========================================================")
	log.Printf("Ghostdreamer LTS Portal Initialized: http://%s", addr)
	log.Printf("Multi-Theme & Auto-Update Engine Active (%s)", AppVersion)
	log.Printf("=========================================================")

	log.Fatal(http.Serve(listener, mux))
}

func resolveTheme() Theme {
	globalState.RLock()
	defer globalState.RUnlock()
	if t, ok := globalState.CustomThemes[globalState.Theme]; ok {
		return t
	}
	if t, ok := defaultThemes[globalState.Theme]; ok {
		return t
	}
	return defaultThemes["dark"]
}

func themeToCSS(t Theme) string {
	return fmt.Sprintf(`--bg:%s;--panel:%s;--border:%s;--text:%s;--muted:%s;--accent:%s;--accent2:%s;`,
		t.Bg, t.Panel, t.Border, t.Text, t.Muted, t.Accent, t.Accent2)
}

func handleManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	th := resolveTheme()
	json.NewEncoder(w).Encode(Manifest{
		Name:            "Ghostdreamer LTS",
		ShortName:       "Ghostdreamer",
		StartURL:        "/?mode=standalone",
		Display:         "standalone",
		BackgroundColor: th.Bg,
		ThemeColor:      th.Accent,
	})
}

func handleGetState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	globalState.RLock()
	defer globalState.RUnlock()
	json.NewEncoder(w).Encode(globalState)
}

func handleUpdateState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload struct {
		AudioSource       string `json:"audio_source"`
		PrivateFeedToggle bool   `json:"private_feed_toggle"`
		VisualizerPreset  string `json:"visualizer_preset"`
		FrequencyTemplate string `json:"frequency_template"`
		TargetFPS         int    `json:"target_fps"`
		Layout            string `json:"layout"`
		Theme             string `json:"theme"`
		ShowVisualizer    bool   `json:"show_visualizer"`
		ShowPlaylist      bool   `json:"show_playlist"`
		ShowControls      bool   `json:"show_controls"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	globalState.Lock()
	if payload.AudioSource != "" {
		globalState.AudioSource = payload.AudioSource
	}
	globalState.PrivateFeedToggle = payload.PrivateFeedToggle
	if payload.VisualizerPreset != "" {
		globalState.VisualizerPreset = payload.VisualizerPreset
	}
	if payload.FrequencyTemplate != "" {
		globalState.FrequencyTemplate = payload.FrequencyTemplate
	}
	if payload.TargetFPS > 0 {
		globalState.TargetFPS = payload.TargetFPS
	}
	if payload.Layout != "" {
		globalState.Layout = payload.Layout
	}
	if payload.Theme != "" {
		globalState.Theme = payload.Theme
	}
	globalState.ShowVisualizer = payload.ShowVisualizer
	globalState.ShowPlaylist = payload.ShowPlaylist
	globalState.ShowControls = payload.ShowControls
	globalState.StreamStatus = "CONFIGURED / READY"
	globalState.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success"}`))
}

func handleAddTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var t Theme
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil || t.Name == "" {
		http.Error(w, "Invalid theme", http.StatusBadRequest)
		return
	}
	globalState.Lock()
	globalState.CustomThemes[t.Name] = t
	globalState.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"added"}`))
}

func handleAddHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Profile string `json:"profile"`
		Metric  string `json:"metric"`
		Passed  bool   `json:"passed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	globalState.Lock()
	newHist := History{
		ID:        fmt.Sprintf("hist-%03d", len(globalState.History)+1),
		Timestamp: time.Now(),
		Profile:   req.Profile,
		Metric:    req.Metric,
		Passed:    req.Passed,
	}
	globalState.History = append([]History{newHist}, globalState.History...)
	globalState.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newHist)
}

func handleVersionCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	info := VersionInfo{
		Current:         AppVersion,
		Latest:          AppVersion,
		UpdateAvailable: false,
		ReleaseURL:      fmt.Sprintf("https://github.com/%s/releases", GitHubRepo),
	}
	client := &http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", GitHubRepo), nil)
	if err == nil {
		req.Header.Set("User-Agent", "Ghostdreamer-LTS-Host")
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var gitHubRelease struct {
				TagName string `json:"tag_name"`
				HTMLURL string `json:"html_url"`
			}
			if json.NewDecoder(resp.Body).Decode(&gitHubRelease) == nil && gitHubRelease.TagName != "" {
				info.Latest = gitHubRelease.TagName
				info.ReleaseURL = gitHubRelease.HTMLURL
				info.UpdateAvailable = (info.Current != info.Latest)
			}
		}
	}
	json.NewEncoder(w).Encode(info)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	globalState.RLock()
	layout := globalState.Layout
	globalState.RUnlock()

	tmplName := "default"
	if layout == "vlc" {
		tmplName = "vlc"
	}

	tmpl, err := template.New(tmplName).ParseFS(contentFS, "index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Template parse error: %v", err), http.StatusInternalServerError)
		return
	}

	globalState.RLock()
	data := struct {
		*AppState
		ThemeCSS string
	}{AppState: globalState, ThemeCSS: themeToCSS(resolveTheme())}
	globalState.RUnlock()

	tmpl.Execute(w, data)
}