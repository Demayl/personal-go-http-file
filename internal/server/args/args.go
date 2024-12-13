package args

import (
	"fmt"
	"log"
	//"internal/server"
	"flag"
	"os"
	"os/user"
	"syscall"
	"path"
	"strings"
//	"reflect"
	"strconv"
)

const (
	OsRead	   = 0o4
	OsUserShift  = 6
	OsGroupShift = 3
	OsOthShift   = 0
	OsUserR 	 = OsRead << OsUserShift
	OsGroupR 	 = OsRead << OsGroupShift
	OsOthR 		 = OsRead << OsOthShift
)

func ParseArgs() (string, int) {
	port := flag.Int("port", 80, "on which port should the server listen on")
	help := flag.Bool("help", false, "shows help")

	flag.Parse();


	if *help {
		showHelp("")
	}
	args_len := len(flag.Args());

	if args_len < 1 {
		showHelp("Missing first argument directory. Usage:\n")
	}

	if args_len > 1 {
		showHelp("Unknown arguments passed after the directory path: " + strings.Join(flag.Args()[1:args_len], ",") );
	}

	path_dir := flag.Args()[0]

	path_info, path_err := os.Stat(path_dir)

	if os.IsNotExist(path_err) {
		showHelp( fmt.Sprintf("Missing directory '%v' ", path_dir) )
	}
	if !path_info.IsDir() {
		showHelp( fmt.Sprintf("'%v' is not a directory", path_dir) )
	}

	if !checkReadable(path_dir) {
		log.Printf("You dont have permission to read from this directory '%v'\n", path_dir)
		os.Exit(1)
	}


	path_full := resolveDir(path_dir)

	return path_full, *port
}


// void shows the help message with optional message and exits
// param message string
// exit 1 on message else exit 0
func showHelp(message string) {

	if message != "" {
		log.Println(message)
	}

	fmt.Println(
		`Start an HTTP server that serves the files in current directory

  ./server [ARGUMENTS] directory
	`)
	flag.PrintDefaults()

	if message != "" {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// TODO move to utils package
// Check if the current user can read the file/directory
// Uses simple file permissions check and it's ignoring any additional ACL's
// It's assumed that root has the file permission ignore flag turned on and it can always read
// param path string
func checkReadable(path_dir string) bool {
	var stat syscall.Stat_t

	if path_dir == "" {
		log.Println("Missing path argument to check if readable")
		return false
	}
	path_info, path_err := os.Lstat(path_dir)

	if path_err != nil {

		log.Printf("Failed to stat '%v': %v\n", path_dir, path_err)
		return false
	}

	if path_err = syscall.Stat(path_dir, &stat); path_err != nil {
		log.Printf("Unable to get stat of '%v'", path_dir)
		return false
	}

	euid := uint32(os.Geteuid());

	if euid == 0 { // root can read, see CAP_DAC_OVERRIDE
		return true
	}

	uid, ugid, err := getCurrentUser()

	if err != nil {
		log.Println( err )
		return false
	}


	if path_info.Mode().Perm()&OsOthR != 0 { // Others can read
		return true
	}

	if inGroup(ugid, stat.Gid) && path_info.Mode().Perm()&OsGroupR != 0 { // Group can read
		return true
	}

	if uid == uint32(stat.Uid) && path_info.Mode().Perm()&OsUserR != 0 { // User can read
		return true
	}


	return false
}

func inGroup(userGroup []uint32, gid uint32) bool {
	for _, gID := range userGroup {
		if uint32(gid) == gID {
			return true
		}
	}
	return false
}

func getCurrentUser() (uint32, []uint32, error) {
	var uid uint32
	var gid []uint32
	curUser, err := user.Current()

	if err != nil {
		return uid, gid, err
	}
	uidInt, err := strconv.ParseUint(curUser.Uid, 10, 32)

	if err != nil {
		return uid, gid, err
	}

	uid = uint32(uidInt)

	userGroups, err := curUser.GroupIds()

	if err != nil {
		return uid, gid, err
	}

	for _, ugid := range userGroups {
		gidInt, err := strconv.ParseUint(ugid, 10, 32)
		if err != nil {
			return uid, gid, err
		}
		gid = append(gid, uint32(gidInt))
	}

	return uid, gid, nil
}

func resolveDir(path_dir string) string {

	var path_full string

	if path_dir == "" {
		return ""
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get the cwd: %v", err)
		os.Exit(1)
	}

	if path.IsAbs(path_dir) {
		path_full = path.Clean(path_dir)
	} else {
		path_full = path.Join( cwd, path.Clean(path_dir) ) 
	}

	return path_full
}
