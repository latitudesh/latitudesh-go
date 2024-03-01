package latitude

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"

	projects "github.com/latitudesh/latitudesh-go/infrastructure/projects"
)

const (
	apiURLEnvVar           = "LATITUDE_API_URL"
	latitudeAccTestVar     = "LATITUDE_TEST_ACTUAL_API"
	TestProjectPrefix      = "LATITUDE_TEST_PROJECT_"
	testPlanVar            = "LATITUDE_TEST_PLAN"
	testSiteVar            = "LATITUDE_TEST_SITE"
	testOperatingSystemVar = "LATITUDE_TEST_OS"
	testSSHKeyVar          = "LATITUDE_TEST_SSH_KEY"
	testUserDataContentVar = "LATITUDE_TEST_USER_DATA_CONTENT"
	testRecorderEnv        = "LATITUDE_TEST_RECORDER"

	testRecorderRecord   = "record"
	testRecorderPlay     = "play"
	testRecorderDisabled = "disabled"
	recorderDefaultMode  = recorder.ModePassthrough

	// defaults should be available to most users
	testSiteDefault            = "SAO"
	TestPlanDefault            = "c2-small-x86"
	TestRegionDefault          = "SAO"
	testSSHKeyDefault          = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDQZtz6DPH4Y04vYLdOch5xOzDY7cdGWpYjBFx5H7ZzieVoRwartZAVTGX4qFT9aoyCuuE6qXYcTj6G1CdO5fb8iOtU6K3FdzVyw/WQ/c4sCehEL+wbYrOnXJSYMhLsUAFhZ69tTdmQSgctbv44yP32Z4xiE4zc/Bk465F3u4Zi1Jj883fyAgzahTWXOxpmvYAEuS6Qv6w4yJc6giiGFVYmu+N6h9j348UgbpToYiCSnSM4iNa9fs7sBGufOa9FuXtggPfXtpyk9f05AhkKEjPlCXcDNAq0GsvN2QEx3tYw6i5ze0qehv6EBAtwx3PLrj636O6IgSh0DgrZBih9NBov"
	testUserDataContentDefault = "bGF0aXR1ZGVzaCB1c2VyIGRhdGEgZXhhbXBsZQ=="
	testOperatingSystemDefault = "ubuntu_22_04_x64_lts"
)

const (
	TestProjectType        = "projects"
	TestProjectEnvironment = "Development"
)

func TestPlan() string {
	envPlan := os.Getenv(testPlanVar)
	if envPlan != "" {
		return envPlan
	}
	return TestPlanDefault
}

func TestUserDataContent() string {
	envUserDataContent := os.Getenv(testUserDataContentVar)
	if envUserDataContent != "" {
		return envUserDataContent
	}
	return testUserDataContentDefault
}

func TestSite() string {
	envSite := os.Getenv(testSiteVar)
	if envSite != "" {
		return envSite
	}
	return testSiteDefault
}

func TestOperatingSystem() string {
	envOS := os.Getenv(testOperatingSystemVar)
	if envOS != "" {
		return envOS
	}
	return testOperatingSystemDefault
}

func TestSSHKey() string {
	envPlan := os.Getenv(testSSHKeyVar)
	if envPlan != "" {
		return envPlan
	}
	return testSSHKeyDefault
}

func RandString8() string {
	// test recorder needs replayable names, not randoms
	mode, _ := testRecordMode()
	if mode != recorder.ModePassthrough {
		return "testrand"
	}

	n := 8
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("acdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// setupWithProject returns a client, project id, and teardown function
// configured for a new project with a test recorder for the named test
func SetupWithProject(t *testing.T) (*Client, string, func()) {
	c, stopRecord := Setup(t)
	rs := TestProjectPrefix + RandString8()
	pcr := projects.ProjectCreateRequest{
		Data: projects.ProjectCreateData{
			Type: TestProjectType,
			Attributes: projects.ProjectCreateAttributes{
				Name:        rs,
				Environment: TestProjectEnvironment,
			},
		},
	}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}

	return c, p.ID, func() {
		_, err := c.Projects.Delete(p.ID)
		if err != nil {
			panic(fmt.Errorf("while deleting %s: %s", p.Name, err))
		}
		stopRecord()
	}
}

func Setup(t *testing.T) (*Client, func()) {
	name := t.Name()
	apiToken := os.Getenv(authTokenEnvVar)
	if apiToken == "" {
		t.Fatalf("If you want to run latitude test, you must export %s.", authTokenEnvVar)
	}

	mode, err := testRecordMode()
	if err != nil {
		t.Fatal(err)
	}
	apiURL := os.Getenv(apiURLEnvVar)
	if apiURL == "" {
		apiURL = baseURL
	}
	r, stopRecord := testRecorder(t, name, mode)
	httpClient := http.DefaultClient
	httpClient.Transport = r
	c, err := NewClientWithBaseURL(apiToken, httpClient, apiURL)
	if err != nil {
		t.Fatal(err)
	}

	return c, stopRecord
}

func ProjectTeardown(c *Client) {
	ps, _, err := c.Projects.List(nil)
	if err != nil {
		panic(fmt.Errorf("while teardown: %s", err))
	}
	for _, p := range ps {
		fmt.Println(p.ID)
		if strings.HasPrefix(p.Name, TestProjectPrefix) {
			_, err := c.Projects.Delete(p.ID)
			if err != nil {
				panic(fmt.Errorf("while deleting %s: %s", p.Name, err))
			}
		}
	}
}

func SkipUnlessAcceptanceTestsAllowed(t *testing.T) {
	if os.Getenv(latitudeAccTestVar) == "" {
		t.Skipf("%s is not set", latitudeAccTestVar)
	}
}

func testRecordMode() (recorder.Mode, error) {
	modeRaw := os.Getenv(testRecorderEnv)
	mode := recorderDefaultMode

	switch strings.ToLower(modeRaw) {
	case testRecorderRecord:
		mode = recorder.ModeRecordOnly
	case testRecorderPlay:
		mode = recorder.ModeReplayOnly
	case "":
		// no-op
	case testRecorderDisabled:
		// no-op
	default:
		return mode, fmt.Errorf("invalid %s mode: %s", testRecorderEnv, modeRaw)
	}
	return mode, nil
}

func testRecorder(t *testing.T, name string, mode recorder.Mode) (*recorder.Recorder, func()) {
	rOptions := recorder.Options{
		CassetteName:  path.Join("fixtures", name),
		Mode:          mode,
		RealTransport: nil,
	}

	r, err := recorder.NewWithOptions(&rOptions)
	if err != nil {
		t.Fatal(err)
	}

	r.AddHook(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Auth-Token")
		return nil
	}, recorder.HookKind(1))

	return r, func() {
		if err := r.Stop(); err != nil {
			t.Fatal(err)
		}
	}
}
