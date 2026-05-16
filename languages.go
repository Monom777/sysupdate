package main

type Translation struct {
	Title        string
	MenuOptions  []string
	Controls     string
	Ready        string
	Canceled     string
	HistoryTitle string
	EmptyHistory string
	BackOption   string
	NoInternet   string
}

var allWorldLanguages = []string{
	"Українська", "English", "Русский", "Polski", "Deutsch", "Français", "Español",
	"Italiano", "Português", "Nederlands", "Türkçe", "Čeština", "Slovenčina", "Magyar",
	"Română", "Български", "Ελληνικά", "Svenska", "Norsk", "Dansk", "Suomi", "Eesti",
	"Latviešu", "Lietuvių", "Српски", "Hrvatski", "Slovenščina", "Македонски", "Shqip",
	"Հայերեն", "ქართული", "عربي", "עברית", "فарسی", "हिन्दी", "বাংলা", "ਪੰਜਾਬੀ", "ગુજરાતી",
	"தமிழ்", "తెలుగు", "ಕನ್ನಡ", "മലയാളം", "ไทย", "Tiếng Việt", "한국어", "日本語", "简体中文",
	"繁體中文", "Қазақ тілі", "Oʻzbekcha", "Кыргызча", "Toҷикӣ", "Türkmençe", "Azərbaycan",
}

var languages = map[string]Translation{
	"Українська": {
		Title: "⚡ Утиліта оновлення системи",
		MenuOptions: []string{
			"Обновити %s + Системні пакети",
			"Обновити тільки %s",
			"Обновити тільки системні пакети (Pacman)",
			"Обновити тільки Флетхаб (Flatpak)",
			"Оновити ВСЕ (Pacman + %s + Flatpak)",
			"🔄 Оновити утиліту з GitHub",
			"Переглянути історію оновлень",
			"Скинути налаштування (Hard Reset)",
			"Вихід",
		},
		Controls:     "Керування: Стрілочки/WS — рух | Enter/D — вибір | Q — вихід",
		Ready:        "\n✅ Все готово!",
		Canceled:     "Скасовано.",
		HistoryTitle: "📜 Історія оновлень (Останні записи):\n",
		EmptyHistory: "  Історія поки що порожня.",
		BackOption:   " Натисніть Q або Enter, щоб повернутися",
		NoInternet:   "❌ Немає підключення до інтернету! Перевірте мережу.",
	},
	"Русский": {
		Title: "⚡ Утилита обновления системы",
		MenuOptions: []string{
			"Обновить %s + Системные пакеты",
			"Обновить только %s",
			"Обновить только системные пакеты (Pacman)",
			"Обновить только Флэтхаб (Flatpak)",
			"Обновить ВСЕ (Pacman + %s + Flatpak)",
			"🔄 Обновить утилиту с GitHub",
			"Просмотреть историю обновлений",
			"Сбросить настройки (Hard Reset)",
			"Выход",
		},
		Controls:     "Управление: Стрелочки/WS — движение | Enter/D — выбор | Q — выход",
		Ready:        "\n✅ Всё готово!",
		Canceled:     "Отменено.",
		HistoryTitle: "📜 История обновлений (Последние записи):\n",
		EmptyHistory: "  История пока пуста.",
		BackOption:   " Нажмите Q или Enter, чтобы вернуться",
		NoInternet:   "❌ Нет подключения к интернету! Проверьте сеть.",
	},
	"English": {
		Title: "⚡ System Update Utility",
		MenuOptions: []string{
			"Update %s + System Packages",
			"Update %s only",
			"Update System Packages (Pacman) only",
			"Update Flatpak only",
			"Update EVERYTHING (Pacman + %s + Flatpak)",
			"🔄 Update Utility from GitHub",
			"View Update History",
			"Reset Settings (Hard Reset)",
			"Exit",
		},
		Controls:     "Controls: Arrows/WS — move | Enter/D — select | Q — exit",
		Ready:        "\n✅ Done!",
		Canceled:     "Canceled.",
		HistoryTitle: "📜 Update History (Latest entries):\n",
		EmptyHistory: "  No history found yet.",
		BackOption:   " Press Q or Enter to go back",
		NoInternet:   "❌ No internet connection! Please check your network.",
	},
}

func getTranslation(lang string) Translation {
	if t, ok := languages[lang]; ok {
		return t
	}
	return languages["English"]
}
