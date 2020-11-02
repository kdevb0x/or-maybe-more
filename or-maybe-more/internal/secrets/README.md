Package secrets simply opens an encrypted 7zip archive of ".pw" files, using the value of "OOM_SECRETS_AR_PASS" environment variable as the password, or will prompt interactively if unset, and for each file contained in the archive, creates an environment variable using the filename without the ".pw" file extension; its value set equal
to the contents of the first line of the file. If there are multiple lines, it only uses the first.

At the present time, this packages offers no means of creating such an encrypted
secrets archive, so this must be created manually using 7zip utility.

For example:

On linux:

	`7z a -sdel -p secrets.7z traviscl.pw ...`

will prompt interactively for the password, and delete the file after adding it.


See `7z --help` for more info on its use.

**IMPORTANT:**
The current version uses a single password, so all of the files inside of the
encrypted archive must have used the same password when they were added to the
archive, even if they were added at different times.

Copyright 2020 kdevb0x ltd
