package asana

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/mikehouston/asana-go"
	"github.com/codeallthethingz/alexa-asana/api/common"
)

type DAO interface {
	CreateTask(task string, user string) error
}
type DAOImpl struct {
	Client *common.Client
}

func NewDAO(client *common.Client) DAO {
	return &DAOImpl{
		Client: client,
	}
}

type teamsResponse struct {
	Data []*asana.Team
}
type projectsResponse struct {
	Data []*asana.Project
}

type createProjectResponse struct {
	Data asana.Project
}
type createTaskResponse struct {
	Data asana.Task
}
type tasksResponse struct {
	Data []*asana.Task
}
type usersResponse struct {
	Data []*asana.User
}
type userResponse struct {
	Data asana.User
}

var users = map[string]string{
	"amzn1.ask.account.AEJJDVS5D3ZBZUFK6IHCNVGEHASUKTWE6ANKPHOJFLDJIYLRF6LOHBCDAUKBSCC4CXTLYDYASHPQJQXJEUI6GGLHU6EFB7IDBTM6V3YJHD7K6OEHNK4D6VI5TAP3UCOI5U7YERSZIXQB6WYGDWENTRJVHLRZ3N5S5WRM663AWHVQSYR5UN4OSDOM5XUJRYAK6JFJG2VRKRZEWYQ": "howdoicontactwill@gmail.com",
	"amzn1.ask.account.AGDHKOZRZ3IAJ6L4HWPFNCBIU4IBTTU2GKG53COT3NEIIAUZVVNWU4DKIDTGEW3HTRBHU32C3ZMDD37M4DTMXB72V45BALL3GTKK2I35VSXPK4Y7OAEU2AAKNZG5ZVA55IGGYOLNJEZQJ7VTPBCDBNYP2L6TWHKP5HA3WPNMXHRINETGSYOCGJJXLIDL4WF3CJE6NT2PFFTQ7WA": "tracie.valdez6@gmail.com",
}

func (d *DAOImpl) CreateTask(task string, user string) error {
	token := os.Getenv("ASANA_ACCESS_TOKEN")
	if strings.TrimSpace(token) == "" {
		return fmt.Errorf("must supply ASANA_ACCESS_TOKEN as an environment variable")
	}
	workspaceID := os.Getenv("ASANA_WORKSPACE_ID")
	if strings.TrimSpace(workspaceID) == "" {
		return fmt.Errorf("must supply ASANA_WORKSPACE_ID as an environment variable")
	}
	projectID := os.Getenv("ASANA_PROJECT_ID")
	if strings.TrimSpace(projectID) == "" {
		return fmt.Errorf("must supply ASANA_PROJECT_ID as an environment variable")
	}
	assignee := ``
	data := `{
		"data": {
			"name": "` + task + `", ` + assignee + `
			"completed": false,
			  "projects": [
				"` + projectID + `"
			  ]
			}
		  }`

	body, err := d.Client.AuthorizedPost(token, "https://app.asana.com/api/1.0/tasks", data)
	if err != nil {
		return err
	}
	taskResponse := &createTaskResponse{}
	if err := json.Unmarshal(body, taskResponse); err != nil {
		return err
	}
	return nil
}
