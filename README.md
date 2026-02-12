# idasenctl

idasenctl is a command line tool for ikea idasen desk that allows you to move the desk up and down using presets. really nerdy.

## How to install

### Binary Releases (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/samueltorres/idasenctl/releases).

### Homebrew (macOS)

```bash
brew tap samueltorres/homebrew-tap
brew install idasenctl
```

### Linux (amd64)

```bash
curl -LO https://github.com/samueltorres/idasenctl/releases/latest/download/idasenctl_<version>_linux_amd64.tar.gz
tar -xzf idasenctl_<version>_linux_amd64.tar.gz
sudo install -m 0755 idasenctl /usr/local/bin/idasenctl
```

### From Source

Install it using `go install`:

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

### List your desks and presets

List all configured desks:

```bash
idasenctl desk list
```

List presets for your default desk:

```bash
idasenctl preset list
```

List presets for a specific desk:

```bash
idasenctl preset list --desk "my-desk"
```

### 3. Move the desk

You can move the desk to a preset by using the `set` command:

```bash
idasenctl set stand
```

## Daemon Mode & Scheduled Movements

idasenctl now supports running as a daemon with scheduled desk movements. This allows you to automatically move your desk based on predefined schedules.

### Setting up schedules

First, create a schedule using the `schedule add` command:

```bash
# Schedule to move to "sit" preset at 9:00 AM on weekdays
idasenctl schedule add "morning-sit" --time 09:00 --preset sit --days monday,tuesday,wednesday,thursday,friday

# Schedule to move to "stand" preset at 2:00 PM on weekdays  
idasenctl schedule add "afternoon-stand" --time 14:00 --preset stand --days mon,tue,wed,thu,fri

# Schedule for a specific desk (if you have multiple desks)
idasenctl schedule add "evening-sit" --time 17:30 --preset sit --desk "my-desk" --days 1,2,3,4,5
```

### Managing schedules

List all configured schedules:

```bash
idasenctl schedule list
```

Remove a schedule:

```bash
idasenctl schedule remove "morning-sit"
```

### Running the daemon

Start the daemon to begin automated desk movements:

```bash
idasenctl daemon
```

The daemon will:
- Check for scheduled movements every minute
- Send an OS notification 10 seconds before moving the desk
- Execute the scheduled movement at the specified time
- Only run schedules on the configured days of the week

### Running in the background (system service)

If you want the scheduler to run automatically in the background, use your OS service manager.

#### macOS (launchd)

1. Copy the sample LaunchAgent file:

```bash
cp contrib/launchd/com.samueltorres.idasenctl.plist ~/Library/LaunchAgents/
```

2. Edit the plist to point `ProgramArguments[0]` to your `idasenctl` binary path.

3. Load and start it:

```bash
launchctl load ~/Library/LaunchAgents/com.samueltorres.idasenctl.plist
launchctl start com.samueltorres.idasenctl
```

To stop it:

```bash
launchctl stop com.samueltorres.idasenctl
launchctl unload ~/Library/LaunchAgents/com.samueltorres.idasenctl.plist
```

#### Linux (systemd user service)

1. Copy the sample service file:

```bash
mkdir -p ~/.config/systemd/user
cp contrib/systemd/idasenctl.service ~/.config/systemd/user/
```

2. Edit `ExecStart` to point to your `idasenctl` binary path (for example `~/go/bin/idasenctl`).

3. Enable and start it:

```bash
systemctl --user daemon-reload
systemctl --user enable --now idasenctl.service
```

To stop it:

```bash
systemctl --user disable --now idasenctl.service
```

### macOS Shortcuts

You can integrate idasenctl with macOS Shortcuts in a few different ways. All of these use the built-in **Run Shell Script** action.

#### Option 1: Start/Stop the daemon from Shortcuts

Create two shortcuts (or a single shortcut with a menu):

- **Start daemon**
```bash
/usr/local/bin/idasenctl daemon
```

- **Stop daemon**
```bash
pkill -f "idasenctl daemon"
```

Notes:
- Update `/usr/local/bin/idasenctl` to your actual binary path.
- If you're using the launchd service from above, you can replace these with:
  - Start: `launchctl start com.samueltorres.idasenctl`
  - Stop: `launchctl stop com.samueltorres.idasenctl`

#### Option 2: One-off move (sit/stand) from Shortcuts

Create a shortcut that runs:

```bash
/usr/local/bin/idasenctl set sit
```

Or:

```bash
/usr/local/bin/idasenctl set stand
```

You can also prompt for the preset name by adding **Ask for Input** in Shortcuts and using it in the command:

```bash
/usr/local/bin/idasenctl set "$SHORTCUT_INPUT"
```

#### Option 3: Use Shortcuts Personal Automations (no daemon)

If you prefer not to run a background service, create a **Personal Automation** in the Shortcuts app:

1. Open Shortcuts → Automation → New Personal Automation.
2. Choose a time of day (e.g., 9:00 AM).
3. Add **Run Shell Script** with:

```bash
/usr/local/bin/idasenctl set sit
```

Repeat for other times/presets.

Notes:
- Personal Automations are scheduled by Shortcuts itself, so you don't need `idasenctl daemon`.
- Make sure the binary path is correct and accessible.

### Schedule Configuration

Schedules support the following options:

- **Time**: 24-hour format (e.g., `09:00`, `14:30`)
- **Days**: Days of the week can be specified as:
  - Full names: `monday`, `tuesday`, etc.
  - Abbreviated: `mon`, `tue`, etc. 
  - Numbers: `0` (Sunday), `1` (Monday), ..., `6` (Saturday)
- **Preset**: Any preset you've configured for the desk
- **Desk**: Specific desk name (defaults to your default desk)
- **Enabled**: Whether the schedule is active (default: true)

### Example Configuration

After setting up schedules, your configuration file (`~/.idasenctl.yaml`) will look like:

```yaml
desks:
  my-desk:
    name: my-desk
    address: "AA:BB:CC:DD:EE:FF"
    presets:
      sit:
        name: sit
        height: 0.75
      stand:
        name: stand
        height: 1.10
defaultDesk: my-desk
schedules:
  - name: morning-sit
    time: "09:00"
    deskName: my-desk
    presetName: sit
    enabled: true
    days: [1, 2, 3, 4, 5]
  - name: afternoon-stand
    time: "14:00"
    deskName: my-desk
    presetName: stand
    enabled: true
    days: [1, 2, 3, 4, 5]
```

### Notifications

The daemon sends OS notifications 10 seconds before moving your desk using the [beeep](https://github.com/gen2brain/beeep) library, which provides cross-platform desktop notifications:

- **macOS**: Uses native notification system
- **Linux**: Uses libnotify/notify-send 
- **Windows**: Uses Windows toast notifications

No additional setup required - notifications should work out of the box on all platforms.
