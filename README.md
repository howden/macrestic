# macrestic

Small wrapper for [restic](https://restic.readthedocs.io) that performs encrypted, incremental backups of local files on a Mac to Backblaze B2.

### Prerequisites

1. Install restic - `brew install restic`
2. Install macrestic - `chmod +x macrestic; mv macrestic /usr/local/bin/`
3. In the Backblaze web ui:
   1. Create a new B2 repository. Set to private, don't enable the object lock.
   2. Generate an "app key" and note down a copy of the values

### Setup

Store some values in the Keychain:

```shell
security add-generic-password -s "backup-restic-b2-account-id" -a restic_backup -w
# enter your B2 account id (generated in step 3.ii) 

security add-generic-password -s "backup-restic-b2-account-key" -a restic_backup -w
# enter your B2 account key (generated in step 3.ii) 

security add-generic-password -s "backup-restic-repo" -a restic_backup -w
# b2:<bucket id>:/restic/<device name>
# e.g. b2:mybackupsbucket:/restic/mymac

security add-generic-password -s "backup-restic-repo-password" -a restic_backup -w
# use a randomly generated string - store it somewhere safe.
```

Initialise the restic repository:
```shell
macrestic init
```

Define some file paths to backup:
```shell
echo "Documents/" >> ~/restic_includes.txt
echo "Pictures/" >> ~/restic_includes.txt
```

When you run a backup for the file time, you will need to grant access the backup directories. To do this, open System Preferences -> Security and Privacy -> Files and Folders (or Full Disk Access), and give permission to macrestic.


### Run a backup
```shell
macrestic backup --verbose --files-from restic_includes.txt
```

You can use launchd to run the backup command automatically on a schedule, e.g. daily :)