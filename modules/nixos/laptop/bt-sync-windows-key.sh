# usage: sudo bt-sync-windows-key  (after pairing on Linux, then on Windows; copies the Windows link key into BlueZ)
ADAPTER=${ADAPTER:-28:7F:CF:AF:D6:E1}
DEVICE=${DEVICE:-3C:B0:ED:BA:2F:26}
PART=${PART:-/dev/sda2}

[[ $EUID -eq 0 ]] || { echo "run as root: sudo bt-sync-windows-key" >&2; exit 1; }

win_adapter=$(tr -d : <<<"$ADAPTER" | tr 'A-F' 'a-f')
win_device=$(tr -d : <<<"$DEVICE" | tr 'A-F' 'a-f')
info=/var/lib/bluetooth/$ADAPTER/$DEVICE/info
[[ -f $info ]] || { echo "no BlueZ pairing at $info, pair on Linux first" >&2; exit 1; }
[[ $(realpath "$info") == "$info" ]] || { echo "$info resolves elsewhere: $(realpath "$info")" >&2; exit 1; }

work=$(mktemp -d)
mnt=$work/win
mkdir "$mnt"
cleanup() { mountpoint -q "$mnt" && umount "$mnt"; rm -rf "$work"; }
trap cleanup EXIT

mount -t ntfs3 -o ro "$PART" "$mnt" 2>/dev/null \
  || ntfs-3g -o ro "$PART" "$mnt"

hive=$mnt/Windows/System32/config/SYSTEM
[[ -f $hive ]] || { echo "no Windows SYSTEM hive on $PART" >&2; exit 1; }

reged -x "$hive" 'HKEY_LOCAL_MACHINE\SYSTEM' \
  "\\ControlSet001\\Services\\BTHPORT\\Parameters\\Keys\\$win_adapter" "$work/keys.reg" >/dev/null

key=$(sed -e ':a' -e '/\\\r\?$/{N;s/\\\r\?\n *//;ba}' "$work/keys.reg" \
  | grep -i "^\"$win_device\"=hex:" | sed 's/.*=hex://; s/[, \r]//g' | tr 'a-f' 'A-F')
[[ $key =~ ^[0-9A-F]{32}$ ]] || { echo "no classic link key for $DEVICE in the Windows registry" >&2; exit 1; }

cp -p "$info" "$info.bak"
cmp -s "$info" "$info.bak"
awk -v k="$key" '/^\[/{s=($0=="[LinkKey]")} s&&/^Key=/{$0="Key="k} {print}' "$info.bak" > "$work/info"
grep -q "^Key=$key$" "$work/info" || { echo "no [LinkKey] section in $info" >&2; exit 1; }
install -m 600 "$work/info" "$info"

systemctl restart bluetooth
echo "link key copied from Windows (backup: $info.bak), bluetooth restarted"
