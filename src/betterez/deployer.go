package betterez

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"

	"github.com/bsphere/le_go"
	"gopkg.in/mgo.v2"
)

// DeployToServer removnig existing db and restoring one from backup
func DeployToServer(archiveName string, deploymentData *DeploymentData) (bool, error) {
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		return false, err
	}
	session.SetMode(mgo.Monotonic, true)
	session.DB(deploymentData.DatabaseName).DropDatabase()
	session.Close()
	err = restoreFromArchive(archiveName, deploymentData)
	if err != nil {
		return false, err
	}
	return true, nil
}

func restoreFromArchive(archiveName string, deploymentData *DeploymentData) error {
	leToken := os.Getenv("LE_TOKEN")
	le, _ := le_go.Connect(leToken)
	commands := []*exec.Cmd{}
	commands = append(commands, exec.Command("rm", "-rf", "dump"))
	commands = append(commands, exec.Command("tar", "-xzf", archiveName))
	commands = append(commands, exec.Command("rm", "-f", fmt.Sprintf("dump/%s/system.users.bson", deploymentData.DatabaseName)))
	commands = append(commands, exec.Command("rm", "-f", fmt.Sprintf("dump/%s/system.users.metadata.json", deploymentData.DatabaseName)))
	var out bytes.Buffer
	if deploymentData.IsAuthenticated() {
		commands = append(commands, exec.Command("mongorestore", "--authenticationDatabase", "admin",
			"-u", deploymentData.Username, "-p", deploymentData.Password))
	} else {
		commands = append(commands, exec.Command("mongorestore"))
	}
	for cmdIndex, cmd := range commands {
		le, _ = le_go.Connect(leToken)
		if cmdIndex < (len(commands) - 1) {
			le.Printf("running %s", cmd.Args)
		} else {
			le.Print("Running mongo restore, this will take a while...")
		}
		cmd.Stdout = &out
		err := cmd.Run()
		le, _ = le_go.Connect(leToken)
		if err != nil {
			le.Printf("Error %v while running %v\n%s", err, cmd.Args, out.String())
			return err
		}
		le.Printf("cmd done with %s", out)
		le.Close()
	}
	le, _ = le_go.Connect(leToken)
	le.Printf("Done restoring mongo...")
	le.Printf("Connecting mongo...")
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		le.Printf("%v error connecting mongo.", err)
		return err
	}
	betterezRole := make([]mgo.Role, 1)
	betterezRole[0] = "dbOwner"
	le.Printf("updating %s, user %s,password:%s", deploymentData.DatabaseName, deploymentData.Username, fmt.Sprintf("%x", sha256.Sum256([]byte(deploymentData.Password))))
	betterezUser := &mgo.User{Password: deploymentData.Password, Username: deploymentData.Username, Roles: betterezRole}
	session.DB(deploymentData.DatabaseName).RemoveUser(deploymentData.Username)
	err = session.DB(deploymentData.DatabaseName).UpsertUser(betterezUser)
	if err != nil {
		le.Printf("%v error while adding user", err)
	} else {
		le.Println("user updated!")
	}
	session.Close()
	return err
}
