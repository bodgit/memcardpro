Output split.gif

Require mkdir
Require memcardpro
Require find

Set FontSize 16
Set Width 1000
Set Height 600
Set Padding 32

Sleep 1.5s
Type "mkdir temp"
Enter
Sleep 1s
Type "memcardpro"
Tab
Type "split temp"
Sleep 1s
Type " testd"
Tab
Type "psx"
Tab
Type "Mem"
Tab 3
Sleep 1s
Type "?.mcd"
Sleep 2.5s
Enter
Sleep 1.5s
Type "find temp -type f"
Sleep 1s
Enter
Sleep 5s
Ctrl+D
