package betterez

import (
	"btrzaws"
	"bytes"
	"errors"
	"log"
	"sshconnector"
)

// AnnounceCompletion - notify that the task completed
func AnnounceCompletion() error {
	if serverKeyFile == "" {
		return errors.New("No key file settings")
	}
	sess, err := btrzaws.GetAWSSession()
	if err != nil {
		return err
	}
	instances, err := btrzaws.GetInstancesWithTags(sess, []*btrzaws.AwsTag{
		btrzaws.NewWithValues("tag:Environment", "support"),
		btrzaws.NewWithValues("tag:Service-Type", "http"),
		btrzaws.NewWithValues("tag:Online", "yes"),
		btrzaws.NewWithValues("instance-state-name", "running"),
		btrzaws.NewWithValues("tag:Repository", "connex2"),
	})
	if err != nil {
		return err
	}
	for _, reservation := range instances {
		for _, currentInstance := range reservation.Instances {
			sshSession, err := sshconnector.CreateSSHSession(*currentInstance.PrivateIpAddress,
				"ubuntu", serverKeyFile,
				22, sshconnector.UseKey)
			if err != nil {
				continue
			}
			var stdoutBuf bytes.Buffer
			sshSession.Stdout = &stdoutBuf
			err = sshSession.Run("sudo service connex2 restart")
			if err != nil {
				log.Println("session error:", err)
			}
			defer sshSession.Close()
		}
	}
	return nil
}
