# dh

Personal utility to hop between directories

## installations
```bash
go install github.com/dimfu/dh
```

and add the following function to your .zshrc or .bashrc
```bash
goto() {
  cd $(dh goto $1);
}
```
