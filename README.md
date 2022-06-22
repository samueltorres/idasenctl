# idasenctl

idasenctl is a command line tool for ikea idasen desk that allows you to move the desk up and down using presets. really nerdy.

## How to install

Currently the  easiest way to install it is by using `go install`

```bash
go install github.com/samueltorres/idasenctl
```

## How to use

### 1. Add your desk

You can add your desk by running the `desk add` command:

```bash
idasenctl desk add
```

This will scan bluetooth devices around you that should be IKEA Idasen desks and show a prompt for you to choose which one.

### 2. Add a preset

Add a preset with the current height using the `preset add` command:

```bash
idasenctl preset add sit --current
```
Or add a preset with the height set explicitly

```bash
idasenctl preset add stand --height 1.00
```

### 3. Move the desk

You can move the desk to a preset by using the `set` command:

```bash
idasenctl set stand
```