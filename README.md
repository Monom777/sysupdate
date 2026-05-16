# ⚡ sysupdate

A fast and convenient system update utility for Arch Linux, written in Go using the `Bubble Tea` TUI framework.

## ✨ Features
* 🌐 **Multi-language Support:** Built-in language search and support upon the first launch.
* 📦 **Pacman & AUR Integration:** Automatically counts available updates before running.
* 🧹 **Auto Cache Cleanup:** Automatically removes stale package files (`pacman -Sc`) after a successful update.
* 🌐 **Network Check:** Quick ping to `1.1.1.1` at startup to prevent timeout errors.
* 🔄 **Self-Update:** Ability to update the utility itself from GitHub directly from the TUI menu.
* 📜 **History & Timer:** Logs all updates with precise execution time tracking.

## 🛠️ Installation

### Build from source:
```bash
git clone [https://github.com/YOUR_NICKNAME/sysupdate.git](https://github.com/YOUR_NICKNAME/sysupdate.git)
cd sysupdate
go build -o sysupdate .
mv sysupdate ~/.local/bin/

---

## 🇷🇺 Russian Version (`README.ru.md`)

```markdown
# ⚡ sysupdate

Утилита для быстрого и удобного обновления системы на Arch Linux, написанная на Go с использованием TUI-фреймворка `Bubble Tea`.

## ✨ Особые возможности (Features)
* 🌐 **Мультиязычность:** Поддержка множества языков со встроенным поиском при первом запуске.
* 📦 **Работа с AUR и Pacman:** Автоматический подсчет доступных для обновления пакетов.
* 🧹 **Автоочистка кэша:** Удаляет устаревшие пакеты (`pacman -Sc`) после успешного апдейта.
* 🌐 **Проверка сети:** Быстрый пинг `1.1.1.1` перед стартом во избежание ошибок таймаута.
* 🔄 **Self-Update:** Возможность обновить саму утилиту с GitHub прямо из TUI-меню.
* 📜 **История и таймер:** Логирование обновлений с фиксацией точного времени выполнения.

## 🛠️ Установка (Installation)

### Компиляция из исходников:
```bash
git clone [https://github.com/ТВОЙ_НИК/sysupdate.git](https://github.com/ТВОЙ_НИК/sysupdate.git)
cd sysupdate
go build -o sysupdate .
mv sysupdate ~/.local/bin/
