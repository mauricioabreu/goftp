package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

func (c *Connection) list(args []string) {
	var filename string
	switch lenargs := len(args); lenargs {
	case 0:
		filename = filepath.Join(c.rootdir, c.workdir)
	case 1:
		filename = filepath.Join(c.rootdir, c.workdir, args[0])
	default:
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}

	c.writeout("150 File status okay; about to open data connection.")
	dc, err := c.dataconnection()
	if err == ErrBadSequence {
		c.writeout("503 Bad sequence of commands.")
		return
	}
	if err != nil {
		log.Println(err)
		c.writeout("425 Can't open data connection.")
		return
	}
	defer dc.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		log.Println(fmt.Sprintf("Could not read file: %s", filename))
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}

	if fileinfo.IsDir() {
		files, err := file.Readdirnames(0) // 0 to read all names
		if err != nil {
			log.Println(err)
			c.writeout("450 Requested file action not taken.")
			return
		}
		for _, name := range formatList(c.curDir(), files) {
			if _, err := fmt.Fprint(dc, name, c.lineterminator()); err != nil {
				log.Println(err)
				c.writeout("426 Connection closed; transfer aborted.")
				return
			}
		}
	} else {
		name := formatList(c.curDir(), []string{filename})[0]
		if _, err := fmt.Fprint(dc, name, c.lineterminator()); err != nil {
			log.Println(err)
			c.writeout("426 Connection closed; transfer aborted.")
			return
		}
	}
	c.writeout("226 Closing data connection. Requested file action successful.")
}

// Format a list of files/directories
// Output:
// drwxr-xr-x  13 mauricio.abreua  staff   416 Dec 28 13:31 videos
// -rw-r--r--   1 mauricio.abreua  staff   431 Dec 27 16:17 foo.txt
func formatList(curDir string, names []string) []string {
	listing := make([]string, 0)
	for _, name := range names {
		file, err := os.Open(filepath.Join(curDir, name))
		if err != nil {
			log.Printf("Skipping %s. Reason %s", name, err)
			continue
		}
		fileinfo, err := file.Stat()
		mode := fileinfo.Mode().String()
		var nlinks int
		var user, group string
		stat, ok := fileinfo.Sys().(*syscall.Stat_t)
		if !ok {
			nlinks = 1 // for systems that don't support inode links
			user = "owner"
			group = "group"
		} else {
			nlinks = links(stat)
			user = getUserFromUID(int(stat.Uid))
			group = getGroupFromGID(int(stat.Gid))
		}
		modtime := fileinfo.ModTime().Format("Jan 2 15:04")
		listing = append(listing, fmt.Sprintf(
			"%s %3d %-8s %-8s %8d %s %s", mode, nlinks, user, group, fileinfo.Size(), modtime, name))
	}
	return listing
}

func links(stat *syscall.Stat_t) int {
	return int(stat.Nlink)
}

func getUserFromUID(uid int) string {
	suid := strconv.Itoa(int(uid))
	user, err := user.LookupId(suid)
	if err != nil {
		return suid // raw uid if the system does not have this info
	}
	return user.Name
}

func getGroupFromGID(gid int) string {
	sgid := strconv.Itoa(int(gid))
	group, err := user.LookupGroupId(sgid)
	if err != nil {
		return sgid // raw uid if the system does not have this info
	}
	return group.Name
}
