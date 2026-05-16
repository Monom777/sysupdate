# A convenient and fast utility for updating Arch Linux (and derivatives), written in Go using the Bubble Tea TUI framework.

* 🌐 Multilingual Support: Language selection on the first launch (settings are saved).
* 📦 Pacman & AUR Integration: Automatically counts available updates using checkupdates and yay/paru.
* 🧹 Auto Cache Cleanup: Performs a safe package cache cleanup (pacman -Sc) after a successful update.
* 🌐 Network Check: Quick ping to 1.1.1.1 to prevent freezing when there is no internet connection.
* 🔄 Self-Update: Ability to update the utility source code from GitHub and rebuild.
* 📜 History & Timer: Logs updates with precise timing of the process.

----------------------------------------------------------------------

### From AUR (Recommended) ###

You can easily install sysupdate-git using your favorite AUR helper:
```bash
yay -S sysupdate-git
```

### Manual Build

If you prefer to build the utility manually from the source code, follow these steps:
```bash
git clone https://github.com/monom777/sysupdate-git.git
cd sysupdate-git
go build -o sysupdate .
mkdir -p ~/.local/bin
mv sysupdate ~/.local/bin/
sysupdate
```
----------------------------------------------------------------------

👤 Developer: Nazar (monomom777) -- dovgalnazar94@gmail.com
