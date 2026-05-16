package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	logPath    = "/.cache/sysupdate.log"
	configPath = "/.config/sysupdate/config.json"
)

type Config struct {
	Language  string `json:"language"`
	AurHelper string `json:"aur_helper"`
	RepoPath  string `json:"repo_path"`
}

type sessionState int

const (
	stateLangSelect sessionState = iota
	stateAurSelect
	stateMenu
	stateHistory
)

type model struct {
	state       sessionState
	langCursor  int
	aurCursor   int
	aurOptions  []string
	currentLang string
	aurHelper   string
	repoPath    string
	t           Translation

	cursor      int
	textInput   textinput.Model
	historyLog  string
	pacmanCount int
	aurCount    int
	hasInternet bool
}

func checkInternet() bool {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", "1.1.1.1:53", timeout)
	return err == nil
}

func getPacmanUpdates() int {
	cmd := exec.Command("checkupdates")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0
	}
	return len(strings.Split(strings.TrimSpace(out.String()), "\n"))
}

func getAurUpdates(helper string) int {
	if _, err := exec.LookPath(helper); err != nil {
		return 0
	}
	cmd := exec.Command(helper, "-Qu")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	count := 0
	for _, l := range lines {
		if l != "" && !strings.Contains(l, " Unable to read") {
			count++
		}
	}
	return count
}

func initialModel(cfg Config) model {
	ti := textinput.New()
	ti.Placeholder = "Type to search language... / Пошук мови..."
	ti.Focus()

	m := model{
		aurOptions:  []string{"yay", "paru"},
		textInput:   ti,
		hasInternet: checkInternet(),
		repoPath:    cfg.RepoPath,
	}

	// Якщо утиліта запущена з папки з .git, оновлюємо або зберігаємо шлях до репозиторію
	if _, err := os.Stat(".git"); err == nil {
		if dir, err := os.Getwd(); err == nil {
			m.repoPath = dir
		}
	}

	if cfg.Language != "" && cfg.AurHelper != "" {
		m.state = stateMenu
		m.currentLang = cfg.Language
		m.aurHelper = cfg.AurHelper
		m.t = getTranslation(cfg.Language)
		if m.hasInternet {
			m.pacmanCount = getPacmanUpdates()
			m.aurCount = getAurUpdates(cfg.AurHelper)
		}
		// Перезаписуємо конфіг, якщо шлях до репо оновився
		saveConfig(m.currentLang, m.aurHelper, m.repoPath)
	} else {
		m.state = stateLangSelect
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" || key == "q" || key == "Q" || key == "й" || key == "Й" {
			if m.state == stateHistory {
				m.state = stateMenu
				return m, nil
			}
			return m, tea.Quit
		}

		if m.state == stateLangSelect {
			switch key {
			case "up", "w", "W", "ц", "Ц":
				if m.langCursor > 0 {
					m.langCursor--
				}
			case "down", "s", "S", "і", "І":
				filteredCount := len(m.getFilteredLangs())
				if m.langCursor < filteredCount-1 {
					m.langCursor++
				}
			case "enter", "d", "D", "в", "В":
				filtered := m.getFilteredLangs()
				if len(filtered) > 0 && m.langCursor < len(filtered) {
					m.currentLang = filtered[m.langCursor]
					m.t = getTranslation(m.currentLang)
					m.state = stateAurSelect
				}
			default:
				m.textInput, cmd = m.textInput.Update(msg)
				filteredCount := len(m.getFilteredLangs())
				if m.langCursor >= filteredCount && filteredCount > 0 {
					m.langCursor = filteredCount - 1
				}
			}
			return m, cmd
		}

		if m.state == stateAurSelect {
			switch key {
			case "up", "w", "W", "ц", "Ц":
				if m.aurCursor > 0 {
					m.aurCursor--
				}
			case "down", "s", "S", "і", "І":
				if m.aurCursor < len(m.aurOptions)-1 {
					m.aurCursor++
				}
			case "enter", "d", "D", "в", "В":
				m.aurHelper = m.aurOptions[m.aurCursor]
				saveConfig(m.currentLang, m.aurHelper, m.repoPath)
				if m.hasInternet {
					m.pacmanCount = getPacmanUpdates()
					m.aurCount = getAurUpdates(m.aurHelper)
				}
				m.state = stateMenu
			}
			return m, nil
		}

		if m.state == stateHistory {
			if key == "enter" || key == "d" || key == "D" || key == "в" || key == "В" {
				m.state = stateMenu
			}
			return m, nil
		}

		switch key {
		case "up", "w", "W", "ц", "Ц":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "s", "S", "і", "І":
			if m.cursor < len(m.t.MenuOptions)-1 {
				m.cursor++
			}
		case "enter", "d", "D", "в", "В":
			if m.cursor == 6 { // Історія переїхала на індекс 6
				m.state = stateHistory
				m.historyLog = readLog()
				return m, nil
			}
			if m.cursor == 7 { // Скидання на 7
				hardReset()
				return m, tea.Quit
			}
			if m.cursor == 8 { // Вихід на 8
				return m, tea.Quit
			}
			if !m.hasInternet && m.cursor < 6 {
				return m, nil
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) getFilteredLangs() []string {
	var filtered []string
	query := strings.ToLower(m.textInput.Value())
	for _, lang := range allWorldLanguages {
		if strings.Contains(strings.ToLower(lang), query) {
			filtered = append(filtered, lang)
		}
	}
	return filtered
}

func (m model) View() string {
	if m.state == stateLangSelect {
		s := "🌐 Select Language (Search):\n"
		s += m.textInput.View() + "\n\n"
		filtered := m.getFilteredLangs()
		if len(filtered) == 0 {
			s += "  ❌ Not found\n"
		} else {
			for i, lang := range filtered {
				cursor := " "
				if m.langCursor == i {
					cursor = "👉"
				}
				s += fmt.Sprintf("%s %s\n", cursor, lang)
			}
		}
		return s
	}

	if m.state == stateAurSelect {
		s := "📦 Select your AUR Helper:\n\n"
		for i, option := range m.aurOptions {
			cursor := " "
			if m.aurCursor == i {
				cursor = "👉"
			}
			s += fmt.Sprintf("%s %s\n", cursor, option)
		}
		return s
	}

	if m.state == stateHistory {
		s := m.t.HistoryTitle
		if m.historyLog == "" {
			s += m.t.EmptyHistory + "\n"
		} else {
			s += m.historyLog + "\n"
		}
		s += "\n" + m.t.BackOption + "\n"
		return s
	}

	s := m.t.Title + "\n"
	if !m.hasInternet {
		s += m.t.NoInternet + "\n"
	}
	s += "\n"

	for i, option := range m.t.MenuOptions {
		cursor := " "
		if m.cursor == i {
			cursor = "👉"
		}

		prefix := ""
		if m.hasInternet {
			switch i {
			case 0:
				prefix = fmt.Sprintf("[%s: %d, Pacman: %d] ", m.aurHelper, m.aurCount, m.pacmanCount)
			case 1:
				prefix = fmt.Sprintf("[%s: %d] ", m.aurHelper, m.aurCount)
			case 2:
				prefix = fmt.Sprintf("[Pacman: %d] ", m.pacmanCount)
			case 4:
				prefix = fmt.Sprintf("[Total: %d] ", m.aurCount+m.pacmanCount)
			}
		}

		var formattedOption string
		if strings.Contains(option, "%s") {
			formattedOption = fmt.Sprintf(option, m.aurHelper)
		} else {
			formattedOption = option
		}

		s += fmt.Sprintf("%s %s%s\n", cursor, prefix, formattedOption)
	}
	s += fmt.Sprintf("\n(%s)\n", m.t.Controls)
	return s
}

func loadConfig() Config {
	var cfg Config
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}
	data, err := os.ReadFile(home + configPath)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	return cfg
}

func saveConfig(lang, helper, repoPath string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	fullPath := home + configPath
	_ = os.MkdirAll(filepath.Dir(fullPath), 0755)
	cfg := Config{Language: lang, AurHelper: helper, RepoPath: repoPath}
	data, _ := json.MarshalIndent(cfg, "", "  ")
	_ = os.WriteFile(fullPath, data, 0644)
}

func hardReset() {
	home, err := os.UserHomeDir()
	if err == nil {
		_ = os.Remove(home + configPath)
		fmt.Println("\n Налаштування повністю скинуто!")
	}
}

func writeLog(message string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	f, err := os.OpenFile(home+logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, _ = f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
}

func readLog() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(home + logPath)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	start := 0
	if len(lines) > 11 {
		start = len(lines) - 11
	}
	var lastLines []string
	for _, line := range lines[start:] {
		if strings.TrimSpace(line) != "" {
			lastLines = append(lastLines, "  "+line)
		}
	}
	return strings.Join(lastLines, "\n")
}

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}

func runCmdInDir(dir, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}

func main() {
	cfg := loadConfig()
	p := tea.NewProgram(initialModel(cfg))
	resModel, err := p.Run()
	if err != nil {
		os.Exit(1)
	}

	m := resModel.(model)
	if m.state == stateHistory || m.cursor >= 7 || !m.hasInternet {
		return
	}

	startTime := time.Now()

	switch m.cursor {
	case 0:
		writeLog("Started: Update " + m.aurHelper + " + System Packages")
		runCmd("sudo", "pacman", "-Syu")
		runCmd(m.aurHelper, "-Sua")
		runCmd("sudo", "pacman", "-Sc", "--noconfirm")
		runCmd(m.aurHelper, "-Sc", "--noconfirm")
		duration := time.Since(startTime).Round(time.Second)
		writeLog(fmt.Sprintf("Success: Update %s + System Packages (Cleaned Cache, took %s)", m.aurHelper, duration))
	case 1:
		writeLog("Started: Update " + m.aurHelper + " only")
		runCmd(m.aurHelper, "-Sua")
		runCmd(m.aurHelper, "-Sc", "--noconfirm")
		duration := time.Since(startTime).Round(time.Second)
		writeLog(fmt.Sprintf("Success: Update %s only (Cleaned Cache, took %s)", m.aurHelper, duration))
	case 2:
		writeLog("Started: Update Pacman only")
		runCmd("sudo", "pacman", "-Syu")
		runCmd("sudo", "pacman", "-Sc", "--noconfirm")
		duration := time.Since(startTime).Round(time.Second)
		writeLog(fmt.Sprintf("Success: Update Pacman only (Cleaned Cache, took %s)", duration))
	case 3:
		writeLog("Started: Update Flatpak only")
		if _, err := exec.LookPath("flatpak"); err == nil {
			runCmd("flatpak", "update", "-y")
			runCmd("flatpak", "uninstall", "--unused", "-y")
		}
		duration := time.Since(startTime).Round(time.Second)
		writeLog(fmt.Sprintf("Success: Update Flatpak only (took %s)", duration))
	case 4:
		writeLog("Started: Full System Update (All)")
		runCmd("sudo", "pacman", "-Syu")
		runCmd(m.aurHelper, "-Sua")
		if _, err := exec.LookPath("flatpak"); err == nil {
			runCmd("flatpak", "update", "-y")
			runCmd("flatpak", "uninstall", "--unused", "-y")
		}
		runCmd("sudo", "pacman", "-Sc", "--noconfirm")
		runCmd(m.aurHelper, "-Sc", "--noconfirm")
		duration := time.Since(startTime).Round(time.Second)
		writeLog(fmt.Sprintf("Success: Full System Update (Cleaned Cache, took %s)", duration))
	case 5:
		if m.repoPath == "" {
			fmt.Println("\n❌ Помилка: Шлях до репозиторію GitHub не знайдено.")
			fmt.Println("Будь ласка, запустіть утиліту один раз безпосередньо з папки репозиторію, щоб вона запам'ятала шлях.")
			return
		}
		writeLog("Started: Self-Update from GitHub")
		fmt.Println("\n🔄 Оновлення коду з GitHub...")
		runCmdInDir(m.repoPath, "git", "pull")
		fmt.Println("🔨 Перезбірка бінарника...")
		runCmdInDir(m.repoPath, "go", "build", "-o", "sysupdate", ".")
		home, _ := os.UserHomeDir()
		targetPath := filepath.Join(home, ".local/bin/sysupdate")
		_ = os.Rename(filepath.Join(m.repoPath, "sysupdate"), targetPath)
		duration := time.Since(startTime).Round(time.Second)
		writeLog(fmt.Sprintf("Success: Self-Update completed successfully (took %s)", duration))
		fmt.Println("✅ Утиліту успішно оновлено!")
	}
}
